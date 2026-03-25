package server

import (
	"context"
	"grpcv2/name"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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

func (ns *NameServer) RegisterName(serviceName, addr string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	client := name.NewNameClient(ns.conn)
	in := &name.NameRequest{
		ServiceName: serviceName,
		Address:     []string{addr},
	}
	_, err := client.Register(context.Background(), in)
	if err != nil {
		log.Println(err)
	}

}
