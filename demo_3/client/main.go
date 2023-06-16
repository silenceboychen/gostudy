package main

import (
	"context"
	"flag"
	"github.com/silenceboychen/gostudy/demo_3/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := c.SayHello(ctx)
	if err != nil {
		log.Fatalf("Upload list err: %v", err)
	}
	for n := 0; n < 5; n++ {
		//向流中发送消息
		err := stream.Send(&helloworld.HelloRequest{Name: "stream client rpc " + strconv.Itoa(n)})
		if err != nil {
			log.Fatalf("stream request err: %v", err)
		}
	}
	//关闭流并获取返回的消息
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("SayHello get response err: %v", err)
	}
	log.Println(res)
}
