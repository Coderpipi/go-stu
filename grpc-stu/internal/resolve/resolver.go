package resolve

import (
	"log"

	"google.golang.org/grpc/resolver"
)

var serverAddress = []string{"localhost:50051", "localhost:50052", "localhost:50053"}

type Builder struct {
}

func (b Builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &simpleResolver{cc: cc}
	r.start()
	return r, nil
}

func (b Builder) Scheme() string {
	return "lbw"
}

type simpleResolver struct {
	cc resolver.ClientConn
}

func (s simpleResolver) start() {
	addrs := make([]resolver.Address, len(serverAddress))
	for i, address := range serverAddress {
		addrs[i] = resolver.Address{Addr: address}
	}

	if err := s.cc.UpdateState(resolver.State{Addresses: addrs}); err != nil {
		log.Fatal(err)
	}
}

func (s simpleResolver) ResolveNow(options resolver.ResolveNowOptions) {
}

func (s simpleResolver) Close() {

}
