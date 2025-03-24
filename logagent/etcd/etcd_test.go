package etcd

import (
	"context"
	_ "embed"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	c *clientv3.Client

	//go:embed conf.json
	CollectConf string
)

func TestMain(m *testing.M) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:12379", "localhost:22379", "localhost:32379"},
		DialTimeout: time.Second * 5,
	})

	if err != nil {
		panic(err)
	}
	c = client
	m.Run()
	c.Close()
}

func TestPut(t *testing.T) {
	_, err := c.Put(context.Background(), "collect_log_conf", CollectConf)
	if err != nil {
		t.Error(err)
	}
}
