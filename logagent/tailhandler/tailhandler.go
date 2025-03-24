package tailhandler

import (
	"context"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"logagent/common"
	"logagent/kafka"
)

type (
	TailHandler struct {
		path    string
		topic   string
		t       *tail.Tail
		closeCh chan struct{}
	}
)

func InitTailHandler(ctx context.Context, conf []*common.CollectEntry) error {
	InitMgr()

	for _, v := range conf {
		handler, err := NewTailHandler(v.Path, v.Topic)
		if err != nil {
			continue
		}
		tailHandlerMgr.tailHandlerMap[handler.path] = handler
		common.SG.Go(func() {
			handler.Start(ctx)
		})
	}

	common.SG.Go(func() {
		tailHandlerMgr.Watch(ctx)
	})

	return nil
}

func NewTailHandler(path, topic string) (*TailHandler, error) {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}

	tailHandler := &TailHandler{
		path:    path,
		topic:   topic,
		closeCh: make(chan struct{}, 1),
	}

	t, err := tail.TailFile(path, config)
	if err != nil {
		logrus.WithError(err).WithField("path", path).Error("new tail handler error")
		return nil, err
	}

	tailHandler.t = t

	return tailHandler, nil
}

func (t *TailHandler) Start(ctx context.Context) {
	logrus.Infof("start collect log, path: %s, subscribe topic: %s", t.path, t.topic)

	for {
		select {
		case <-ctx.Done():
			logrus.WithError(ctx.Err()).Error("context canceld")
			return
		case <-t.closeCh:
			logrus.WithFields(logrus.Fields{
				"path":  t.path,
				"topic": t.topic,
			}).Info("handler closed")
			return
		case msg, ok := <-t.t.Lines:
			if !ok {
				logrus.WithError(msg.Err).Error("tail log err")
				time.Sleep(time.Second)
				continue
			}

			if len(strings.TrimSpace(msg.Text)) == 0 {
				continue
			}

			logrus.WithField("msg", msg.Text).Info("read log from log file")

			producerMsg := &sarama.ProducerMessage{Topic: t.topic, Value: sarama.StringEncoder(msg.Text)}
			kafka.MsgChan <- producerMsg

		}

	}

}

func (t *TailHandler) Close() {
	close(t.closeCh)
}
