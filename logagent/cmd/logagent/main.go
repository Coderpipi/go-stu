package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
	"logagent/common"
	"logagent/config"
	"logagent/etcd"
	"logagent/kafka"
	"logagent/tailhandler"
	"logagent/util"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	app := &cli.Command{
		Name:        "logagent",
		Usage:       "日志收集程序",
		Description: "日志收集， 使用 kafka + tail 进行日志收集",
		Commands: []*cli.Command{
			{
				Name:   "collect",
				Usage:  "监听文件，开始收集日志",
				Action: collectCommand,
			},
		},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	common.SG.Go(func() {
		if err := app.Run(ctx, os.Args); err != nil &&
			!errors.Is(err, context.Canceled) {
			logrus.WithError(err).Error("run app error")
		}
	})

	<-ctx.Done()
	common.SG.Wait()
}

func collectCommand(ctx context.Context, command *cli.Command) error {

	if err := config.InitKafka(); err != nil {
		return err
	}

	common.SG.Go(func() {
		kafka.SendMessage(ctx)
	})

	localIP, err := util.GetOutBoundIP()
	if err != nil {
		return err
	}

	collectKey := fmt.Sprintf(config.Cfg.CollectKey, localIP)

	if err := config.InitEtcd(); err != nil {
		return err
	}

	conf, err := etcd.GetConf(collectKey)
	if err != nil {
		return err
	}

	// watch etcd conf
	common.SG.Go(func() {
		etcd.WatchConf(ctx, collectKey)
	})

	if err := tailhandler.InitTailHandler(ctx, conf); err != nil {
		return err
	}

	return nil
}
