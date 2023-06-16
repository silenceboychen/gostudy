package main

import (
	"context"
	"flag"
	"github.com/silenceboychen/gostudy/demo_3/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"strconv"
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

	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := c.SayHello(context.Background())
	if err != nil {
		log.Fatalf("stream err: %v", err)
	}
	for n := 0; n < 5; n++ {
		err := stream.Send(&helloworld.HelloRequest{Name: *name + "_" + strconv.Itoa(n)})
		if err != nil {
			log.Fatalf("client stream err: %v", err)
		}
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("server stream err: %v", err)
		}
		// 打印返回值
		log.Println(res.GetMessage())
	}

	err = stream.CloseSend()
	if err != nil {
		log.Fatalf("close stream err: %v", err)
	}
}
