package server

import (
	"context"
	"fmt"
	"grpcv2/echo"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EchoServer struct {
	echo.UnimplementedEchoServiceServer
}

func (EchoServer) UnaryEcho(ctx context.Context, req *echo.EchoRequest) (*echo.EchoResponse, error) {
	return &echo.EchoResponse{Message: "server got your message: " + req.Message}, nil
}

func (EchoServer) ServerStreamingEcho(req *echo.EchoRequest, stream echo.EchoService_ServerStreamingEchoServer) error {
	fmt.Printf("server recv: %v\n", req.Message)
	filepath := "echo-server/file/server.jpg"
	file, err := os.Open(filepath)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to open file: %v", err)
	}
	defer file.Close()

	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if n == 0 {
			break
		}
		stream.Send(&echo.EchoResponse{
			Bytes:   buf[:n],
			Message: "server send image",
		})
	}
	return nil
}

func (EchoServer) ClientStreamingEcho(stream echo.EchoService_ClientStreamingEchoServer) error {
	fmt.Println("server start to recv client streaming")
	filepath := "echo-server/file/" + strconv.FormatInt(time.Now().UnixMilli(), 10) + ".jpg"
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to create file: %v", err)
	}
	defer file.Close()
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
			break
		}
		file.Write(req.Bytes[:len(req.Bytes)])
		fmt.Printf(
			"server recv\nMessage: %v\nTimestamp: %v\nLength: %v\n", req.Message, req.Timestamp, req.Length)
	}
	err = stream.SendAndClose(&echo.EchoResponse{Message: "server got your image"})
	return err
}

func (EchoServer) BidirectionalStreamingEcho(stream echo.EchoService_BidirectionalStreamingEchoServer) error {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		filepath := "echo-server/file/" + strconv.FormatInt(time.Now().UnixMilli(), 10) + ".jpg"
		file, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			file.Write(req.Bytes[:len(req.Bytes)])
			fmt.Printf("server recv\nMessage: %v\nTimestamp: %v\nLength: %v\n", req.Message, req.Timestamp, req.Length)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		filepath := "echo-server/file/server.jpg"
		file, err := os.Open(filepath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		buf := make([]byte, 1024)
		for {
			n, err := file.Read(buf)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			if n == 0 {
				break
			}
			stream.Send(&echo.EchoResponse{
				Message: "server send image",
				Bytes:   buf[:n],
			})
		}
	}()

	wg.Wait()
	return nil
}
