package main

import (
	"flag"
	"fmt"
	"grpcv2/name"
	"grpcv2/name-server/server"
	"log"
	"net"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 60051, "")
)

func main() {
	//testdata()

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	name.RegisterNameServer(s, &server.NameServer{})
	log.Printf("server listening at : %v", lis.Addr())
	err = s.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}

}

func testdata() {
	// 最终在 NameServer 的内存里形成这样的数据：
	server.Register("myecho", "localhost:50051")
	alldata := server.GetAllData()
	fmt.Println(alldata)
	fmt.Println(server.GetByServiceName("myecho"))
}
