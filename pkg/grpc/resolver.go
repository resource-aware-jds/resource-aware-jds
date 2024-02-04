package grpc

import (
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"google.golang.org/grpc/resolver"
)

func ProvideRAJDSGRPCResolver() RAJDSGRPCResolver {
	res := RAJDSGRPCResolver{
		AddressMapper: make(map[string][]string),
	}
	resolver.Register(&res)
	return res
}

type RAJDSGRPCResolver struct {
	AddressMapper map[string][]string
}

func (r *RAJDSGRPCResolver) AddHost(name string, value string) {
	r.AddressMapper[name] = []string{value}
}

func (r *RAJDSGRPCResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	res := &customResolver{
		target:     target,
		cc:         cc,
		addrsStore: r.AddressMapper,
	}
	res.start()
	return res, nil
}

func (r *RAJDSGRPCResolver) Scheme() string {
	return cert.DefaultScheme
}

type customResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *customResolver) start() {
	addrStrs := r.addrsStore[r.target.URL.Hostname()]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs}) //nolint:errcheck
}
func (*customResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*customResolver) Close()                                  {}
