package distributed_cache_stu

import (
	"errors"
	"log"
	"sync"
)

type (
	Getter interface {
		Get(key string) ([]byte, error)
	}

	GetterFunc func(key string) ([]byte, error)

	Group struct {
		name      string
		getter    Getter
		mainCache cache
		peers     PeerPicker
	}
)

var (
	mu     sync.RWMutex
	groups = map[string]*Group{}
)

func (g GetterFunc) Get(key string) ([]byte, error) {
	return g(key)
}

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil getter")
	}

	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}

	g.peers = peers
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, errors.New("key is required")
	}

	if v, ok := g.mainCache.Get(key); ok {
		log.Println("[Cache] hit")
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	if g.peers == nil {
		log.Println("get from locally")
		return g.getLocally(key)
	}

	if peer, ok := g.peers.PickPeer(key); ok {
		log.Println("get from peer....")
		if value, err := g.getFromPeer(peer, key); err != nil {
			return value, nil
		}
		log.Println("get from peer failed")

	}

	log.Println("get from locally")
	return g.getLocally(key)
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}

	return ByteView{b: bytes}, nil
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.Add(key, value)
}
