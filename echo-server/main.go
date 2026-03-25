package main

import (
	"context"
	"flag"
	"fmt"
	"grpcv2/echo"
	"grpcv2/echo-server/server"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var (
	port = flag.Int("port", 50056, "The server port")
)

func getOptions() []grpc.ServerOption {
	var opts []grpc.ServerOption
	//opts = append(opts, server.GetTlsOpt())
	opts = append(opts, server.GetMTlsOpt())
	opts = append(opts, grpc.UnaryInterceptor(server.UnaryInterceptor))
	opts = append(opts, grpc.StreamInterceptor(server.StreamInterceptor))
	opts = append(opts, server.GetKeepaliveOpt()...)
	return opts
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		fmt.Printf("Failed to listen: %v", err)
		return
	}

	s := grpc.NewServer(getOptions()...)
	echo.RegisterEchoServiceServer(s, &server.EchoServer{})

	h := health.NewServer()
	// 设置默认状态为 SERVING
	h.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(s, h)

	log.Printf("server listening at:%v\n", lis.Addr())

	go func() {
		if err := s.Serve(lis); err != nil {
			fmt.Printf("Failed to serve: %v", err)
		}
	}()

	nameServer := server.NewNameServer("localhost:60051")
	serviceName := "myecho"
	addr := fmt.Sprintf("localhost:%d", *port)
	go func() {
		nameServer.RegisterName(serviceName, addr)
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	<-ctx.Done()

}
