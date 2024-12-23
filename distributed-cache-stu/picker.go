package distributed_cache_stu

import (
	"sync"

	"distributed-cache-stu/consistenthash"
)

const (
	defaultReplicas = 50
)

type (
	GrpcPool struct {
		selfAddr    string
		mu          sync.Mutex
		peers       *consistenthash.Map
		grpcGetters map[string]*grpcGetter
	}
)

func NewGrpcPool(addr string) *GrpcPool {
	return &GrpcPool{selfAddr: addr}
}

func (g *GrpcPool) Set(peers ...string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.peers = consistenthash.New(defaultReplicas, nil)
	g.peers.Add(peers...)

	g.grpcGetters = make(map[string]*grpcGetter, len(peers))
	for _, peer := range peers {
		g.grpcGetters[peer] = &grpcGetter{Addr: peer}
	}
}

func (g *GrpcPool) PickPeer(key string) (PeerGetter, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if peer := g.peers.Get(key); peer != "" && peer != g.selfAddr {
		return g.grpcGetters[peer], true
	}

	return nil, false
}

var _ PeerPicker = (*GrpcPool)(nil)
