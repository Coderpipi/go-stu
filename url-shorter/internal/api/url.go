package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"url-shorter/internal/model"
)

type (
	IUrlService interface {
		CreateUrl(ctx context.Context, req *model.CreateUrlRequest) (*model.CreateUrlResponse, error)
		GetUrl(ctx context.Context, shortCode string) (string, error)
		DeleteUrl(ctx context.Context) error
	}

	UrlHandler struct {
		urlServ IUrlService
	}
)

func NewUrlHandler(urlServ IUrlService) *UrlHandler {
	return &UrlHandler{
		urlServ: urlServ,
	}
}

func (h *UrlHandler) CreateUrl(c *gin.Context) {
	req := &model.CreateUrlRequest{}

	if err := c.ShouldBind(req); err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "bad request",
		})
		return
	}

	// 调用业务函数
	resp, err := h.urlServ.CreateUrl(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 200,
		"msg":  "ok",
		"data": resp,
	})

}

func (h *UrlHandler) RedirectUrl(c *gin.Context) {
	shortCode := c.Param("code")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "code is empty",
		})
		return
	}

	url, err := h.urlServ.GetUrl(c.Request.Context(), shortCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "internal server error",
		})
		return
	}

	c.Redirect(http.StatusPermanentRedirect, url)
}
