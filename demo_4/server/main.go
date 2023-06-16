package main

import (
	"flag"
	"fmt"
	"github.com/silenceboychen/gostudy/demo_3/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
	"strconv"
)

var (
	port = flag.Int("port", 8080, "The server port")
)

type server struct {
	helloworld.UnimplementedHelloServer
}

func (s *server) SayHello(stream helloworld.Hello_SayHelloServer) error {
	n := 0
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		err = stream.Send(&helloworld.HelloReply{
			Message: "server stream: " + res.GetName() + "_" + strconv.Itoa(n),
		})
		if err != nil {
			return err
		}
		n++
		log.Printf("client stream: %s", res.GetName())
	}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	helloworld.RegisterHelloServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
