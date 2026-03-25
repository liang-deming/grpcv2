package server

import (
	"context"
	"fmt"
	"grpcv2/name"
	"log"
)

type NameServer struct {
	name.UnimplementedNameServer
}

func (NameServer) Register(ctx context.Context, in *name.NameRequest) (*name.NameResponse, error) {
	for _, address := range in.Address {
		Register(in.ServiceName, address)
	}
	log.Println(GetByServiceName(in.ServiceName))
	return &name.NameResponse{ServiceName: in.ServiceName}, nil
}

func (NameServer) GetAddress(ctx context.Context, in *name.NameRequest) (*name.NameResponse, error) {
	addr := GetByServiceName(in.ServiceName)
	fmt.Println(in.ServiceName)
	log.Println(addr)
	return &name.NameResponse{ServiceName: in.ServiceName, Address: addr}, nil
}
