package loadbalancer

import (
	"log"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/metadata"
)

const Name = "ab_testing"

func NewBuilder(groups map[string]string, defaultAddr string) balancer.Builder {
	return base.NewBalancerBuilder(Name, &pickerBuilder{
		groups:      groups,
		defaultAddr: defaultAddr,
	}, base.Config{HealthCheck: true})
}

type pickerBuilder struct {
	groups      map[string]string
	defaultAddr string
}

func (p pickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	scs := make(map[string]balancer.SubConn)

	for sc, inf := range info.ReadySCs {
		scs[inf.Address.Addr] = sc
	}

	return &picker{
		groups:      p.groups,
		defaultAddr: p.defaultAddr,
		subConns:    scs,
	}
}

type picker struct {
	groups      map[string]string
	defaultAddr string
	subConns    map[string]balancer.SubConn
}

func (p picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	// 从请求获取元数据
	md, ok := metadata.FromOutgoingContext(info.Ctx)
	if !ok {
		log.Println("unable to get metadata from context, using default address")
		return p.defaultConn()
	}
	// 从请求头中获取数据
	name := md.Get("user-group")
	if len(name) < 1 {
		log.Println("user group not specified in metadata, using default address")
		return p.defaultConn()
	}
	// 检查用户组对应的链接
	addr, ok := p.groups[name[0]]
	if !ok {
		log.Println("group not in list, using default address")
		return p.defaultConn()
	}

	// 将地址映射到子链接
	subConn, ok := p.subConns[addr]
	if !ok {
		log.Println("address is not in list of address, using default address")
		return p.defaultConn()
	}
	// 返回选择的子链接
	return balancer.PickResult{SubConn: subConn}, nil
}

func (p picker) defaultConn() (balancer.PickResult, error) {
	conn, ok := p.subConns[p.defaultAddr]
	if !ok {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	return balancer.PickResult{SubConn: conn}, nil
}
