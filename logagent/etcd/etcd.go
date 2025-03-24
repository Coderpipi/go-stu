package etcd

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"logagent/common"
	"logagent/tailhandler"
)

var (
	Cli *clientv3.Client
)

func GetConf(key string) ([]*common.CollectEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	res, err := Cli.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("get conf from etcd by %s, err: %w", key, err)
	}

	if len(res.Kvs) == 0 || res.Kvs[0] == nil {
		logrus.Warnf("get conf from etcd not found, key:%s", key)
		return []*common.CollectEntry{}, nil
	}

	v := res.Kvs[0].Value

	collectEntries := make([]*common.CollectEntry, 0)

	err = json.Unmarshal(v, &collectEntries)
	if err != nil {
		return nil, err
	}

	return collectEntries, nil
}

func WatchConf(ctx context.Context, confKey string) {
	logrus.Info("starting watch the conf...")

	wCh := Cli.Watch(ctx, confKey)
	newConf := make([]*common.CollectEntry, 0)
	for wr := range wCh {
		for _, evt := range wr.Events {
			logrus.Infof("type: %s, key: %s, val:%s", evt.Type, evt.Kv.Key, evt.Kv.Value)

			err := json.Unmarshal(evt.Kv.Value, &newConf)
			if err != nil {
				logrus.WithError(err).Error("get new conf failed")
				continue
			}

			// 通知 tailhandler 启用新配置
			tailhandler.SendNewConf(newConf)
		}
	}

}
