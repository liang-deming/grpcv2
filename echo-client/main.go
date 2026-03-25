package main

import (
	"flag"
	"fmt"
	"grpcv2/echo"
	"grpcv2/echo-client/client"
	"grpcv2/echo-client/client_pool"
	"log"
	"time"

	"google.golang.org/grpc"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func getOptions() []grpc.DialOption {
	var opts []grpc.DialOption
	//opts = append(opts, client.GetTlsOpt())
	opts = append(opts, client.GetMTlsOpt())
	opts = append(opts, grpc.WithUnaryInterceptor(client.UnaryInterceptor))
	opts = append(opts, grpc.WithStreamInterceptor(client.StreamInterceptor))
	opts = append(opts, client.GetAuth(client.FetchToken()))
	opts = append(opts, client.GetNameResolver(client.NewNameServer("localhost:60051")))
	opts = append(opts, client.GetKeepaliveOpt())
	return opts
}

func main() {
	//flag.Parse()
	//conn, err := grpc.Dial(*addr, getOptions()...)
	//根据 协议 + 服务名 通过名称服务解析，访问服务器
	//conn, err := grpc.Dial(fmt.Sprintf("%s:///%s", client.MyScheme, client.MyServiceName), getOptions()...)

	pool, err := client_pool.GetPool(fmt.Sprintf("%s:///%s", client.MyScheme, client.MyServiceName), getOptions()...)
	if err != nil {
		log.Fatal(err)
	}
	conn := pool.Get()
	defer pool.Put(conn)

	c := echo.NewEchoServiceClient(conn)
	client.CallUnary(c)
	time.Sleep(8 * time.Second)
	client.CallUnary(c)
	time.Sleep(8 * time.Second)
	client.CallUnary(c)
	time.Sleep(8 * time.Second)
	client.CallUnary(c)
	time.Sleep(8 * time.Second)
	client.CallUnary(c)
	time.Sleep(8 * time.Second)
	//client.CallServerStream(c)
	//client.CallClientStream(c)
	//client.CallBidirectional(c)

}
