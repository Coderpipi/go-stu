package main

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc"
	"github.com/urfave/cli/v3"
	"logagent/config"
	"logagent/kafka"
	"logagent/tailhandler"
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
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "topic",
						Usage:    "指定topic名称",
						Value:    "log_collect",
						Required: false,
						OnlyOnce: true,
					},
				},
			},
		},
	}

	_ = app.Run(context.Background(), os.Args)

}

func collectCommand(ctx context.Context, command *cli.Command) error {

	if err := config.InitKafka(); err != nil {
		return err
	}

	if err := config.InitTail(); err != nil {
		return err
	}

	topic := command.String("topic")

	g := conc.WaitGroup{}

	g.Go(kafka.SendMessage)

	for {
		select {
		case <-ctx.Done():
			close(kafka.MsgChan)
			return ctx.Err()
		case msg, ok := <-tailhandler.TailHandler.Lines:
			if !ok {
				logrus.WithError(msg.Err).Error("tail log err")
				time.Sleep(time.Second)
				continue
			}

			if len(strings.TrimSpace(msg.Text)) == 0 {
				continue
			}

			logrus.WithField("msg", msg.Text).Info("read log from log file")

			producerMsg := &sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(msg.Text)}
			kafka.MsgChan <- producerMsg
		}

	}
}
