package main

import (
	"context"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var c *clientv3.Client

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

func TestEtcdOperation(t *testing.T) {
	ctx := context.Background()
	t.Run("test put", func(t *testing.T) {
		_, err := c.Put(ctx, "foo", "bar")
		if err != nil {
			t.Error("put foo to etcd fail")
		}
	})
	t.Run("test get", func(t *testing.T) {
		get, err := c.Get(ctx, "foo")
		if err != nil {
			t.Error("get foo from etcd fail")
		}
		for _, kv := range get.Kvs {
			t.Logf("key: %s, v: %s", string(kv.Key), string(kv.Value))
		}
	})
}

func TestEtcdWatch(t *testing.T) {
	ctx := context.Background()
	watch := c.Watch(ctx, "foo")
	go func() {
		time.Sleep(time.Second)
		c.Put(ctx, "foo", "bar")
		c.Put(ctx, "foo", "bar1")
		c.Delete(ctx, "foo")
	}()
	for w := range watch {
		for _, event := range w.Events {
			t.Logf("type:%s, key: %s, value: %s", event.Type, event.Kv.Key, event.Kv.Value)
		}
	}
}
