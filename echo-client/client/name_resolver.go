package client

import (
	"context"
	"grpcv2/name"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

const (
	MyScheme      = "myscheme"
	MyServiceName = "myecho"
)

// var addrs = []string{"localhost:50051", "localhost:50052", "localhost:50053"}
var nameServer *NameServer

func GetNameResolver(ns *NameServer) grpc.DialOption {
	nameServer = ns
	return grpc.WithResolvers(&MyResolverBuilder{})
}

type MyResolverBuilder struct {
}

func (*MyResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &MyResolver{
		target:     target,
		cc:         cc,
		addrsStore: map[string][]string{MyServiceName: nameServer.getAddressByServiceName(MyServiceName)},
	}
	r.start()
	return r, nil
}

func (*MyResolverBuilder) Scheme() string {
	return MyScheme
}

type MyResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *MyResolver) start() {
	log.Println("Resolver start")
	addrStrs := r.addrsStore[r.target.Endpoint()]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{
			Addr: s,
		}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}

func (r *MyResolver) ResolveNow(o resolver.ResolveNowOptions) {
	log.Println("Resolve Now")
	log.Println(r.cc)
	r.addrsStore = map[string][]string{MyServiceName: {"localhost:50054", "localhost:50055", "localhost:50056"}}
	r.start()
	log.Println(r.cc)
}

func (r *MyResolver) Close() {}

type NameServer struct {
	conn *grpc.ClientConn
}

func NewNameServer(addr string) *NameServer {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	return &NameServer{
		conn: conn,
	}
}

func (ns *NameServer) getAddressByServiceName(serviceName string) []string {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	client := name.NewNameClient(ns.conn)
	in := &name.NameRequest{
		ServiceName: serviceName,
	}
	res, err := client.GetAddress(context.Background(), in)
	if err != nil {
		log.Println(err)
		return []string{}
	}
	log.Println(res.Address)
	return res.Address
}
