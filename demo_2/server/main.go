package main

import (
	"flag"
	"fmt"
	"github.com/silenceboychen/gostudy/demo_2/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

var (
	port = flag.Int("port", 8080, "The server port")
)

type server struct {
	helloworld.UnimplementedHelloServer
}

func (s *server) SayHello(in *helloworld.HelloRequest, stream helloworld.Hello_SayHelloServer) error {
	log.Printf("Received: %v", in.GetName())
	for i := 0; i < 5; i++ {
		stream.Send(&helloworld.HelloReply{Message: fmt.Sprintf("hello %s---%d", in.GetName(), i)})
	}
	return nil
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
