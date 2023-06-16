package main

import (
	"context"
	"flag"
	"github.com/silenceboychen/gostudy/demo_2/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"time"
)

var (
	addr = flag.String("addr", "localhost:8080", "the address to connect to")
	name = flag.String("name", "world", "Name to greet")
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := helloworld.NewHelloClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := c.SayHello(ctx, &helloworld.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not call: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("stream error: %v", err)
		}

		log.Printf("%s", res.Message)
	}
}
