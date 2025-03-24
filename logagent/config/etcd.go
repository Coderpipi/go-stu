package config

import (
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"logagent/etcd"
)

func InitEtcd() error {
	cli, err := clientv3.New(clientv3.Config{Endpoints: Cfg.Etcd.Addr})
	if err != nil {
		return err
	}

	etcd.Cli = cli
	logrus.Info("etcd init success...")
	return nil
}
