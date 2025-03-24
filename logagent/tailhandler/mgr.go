package tailhandler

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
	"logagent/common"
)

var (
	tailHandlerMgr *mgr
)

type (
	// Mgr handler管理者
	mgr struct {
		tailHandlerMap map[string]*TailHandler
		collectEntries []*common.CollectEntry
		confChan       chan []*common.CollectEntry
	}
)

func (m *mgr) Watch(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(m.confChan)
			return
		case newConf := <-m.confChan:
			newTailHandlerMap := make(map[string]*TailHandler, len(m.tailHandlerMap))
			logrus.Infof("conf: %v", newConf)
			for _, c := range newConf {
				// 如果handler已经存在且配置没变，保留
				if m.exist(c) {
					newTailHandlerMap[c.Path] = m.tailHandlerMap[c.Path]
					delete(m.tailHandlerMap, c.Path)
					continue
				}
				// 否则新建handler
				handler, err := NewTailHandler(c.Path, c.Topic)
				if err != nil {
					logrus.WithField("path", c.Path).
						WithField("topic", c.Topic).
						WithError(err).Error("new tail handler err")
					continue
				}

				newTailHandlerMap[handler.path] = handler
				common.SG.Go(func() {
					handler.Start(ctx)
				})

			}

			// 关闭掉那些曾经有但现在没有的 handler
			for _, v := range m.tailHandlerMap {
				v.Close()
			}

			clear(m.tailHandlerMap)
			m.tailHandlerMap = newTailHandlerMap
		}
	}
}

func (m *mgr) exist(conf *common.CollectEntry) bool {
	v, ok := m.tailHandlerMap[conf.Path]
	if !ok {
		return false
	}
	return v.topic == conf.Topic
}

func SendNewConf(conf []*common.CollectEntry) {
	tailHandlerMgr.confChan <- conf
}

func InitMgr() {
	sync.OnceFunc(func() {
		tailHandlerMgr = &mgr{
			tailHandlerMap: make(map[string]*TailHandler),
			confChan:       make(chan []*common.CollectEntry),
		}
	})()
}
