package cmd

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"golang.org/x/sync/errgroup"
	"url-shorter/config"
	"url-shorter/database"
	"url-shorter/internal/api"
	"url-shorter/internal/cache"
	"url-shorter/internal/service"
	"url-shorter/pkg/shortcode"
	"url-shorter/pkg/valid"
)

func Web(ctx context.Context) error {
	config.InitConfig("config/config.yaml")
	cfg := config.Cfg

	shortCode := shortcode.NewShortCode(cfg.ShortCodeConfig.Length)
	db, err := database.NewDB(cfg.DBConfig)
	if err != nil {
		return err
	}

	urlService := service.NewUrlService(db, shortCode, cfg.App.DefaultDuration, cfg.App.BaseUrl)

	urlHandler := api.NewUrlHandler(urlService)

	r := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("url", valid.Url)
	}

	r.GET("/:code", urlHandler.RedirectUrl)
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.String(http.StatusOK, "pong")
	// })
	urlShorterApi := r.Group("/api/url")
	{
		urlShorterApi.POST("", urlHandler.CreateUrl)
	}

	var eg errgroup.Group

	eg.Go(func() error {

		return r.Run(cfg.ServerConfig.Addr)
	})
	cctx, cancel := context.WithCancel(ctx)
	eg.Go(func() error {
		for {
			select {
			case <-cctx.Done():
				return cctx.Err()
			case <-time.Tick(cfg.App.CleanupInterval):
				err := urlService.DeleteUrl(cctx)
				if err != nil {
					slog.Error(err.Error())
				}
			}
		}
	})

	notifyContext, stop := signal.NotifyContext(ctx, os.Kill, os.Interrupt)
	defer stop()

	select {
	case <-notifyContext.Done():
		cancel()
		cache.RedisCli.Close()

	}
	return eg.Wait()
}
