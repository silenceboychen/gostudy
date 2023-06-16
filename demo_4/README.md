# grpc系列课程（四）：双向流式rpc

> 开发环境：
>
> 系统： ubuntu20.04
>
> go版本： 1.19
>
> 编辑器： goland

``grpc``：远程过程调用，使用场景很多，也是比较流行的技术之一。使用go开发grpc服务，除了必须的go语言开发环境之外，还需要安装grpc相关命令。

## grpc环境配置

### protoc安装

```bash
$ sudo apt install -y protobuf-compiler
$ protoc --version

libprotoc 3.6.1
```

如果是其他系统电脑，安装protoc可参考文档：[Protocol Buffer Compiler Installation](https://grpc.io/docs/protoc-installation/)

## protocol编译插件安装

```bash
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

安装完成后可以在bin目录下看到相关指令：

```bash
$ ls $GOPATH/bin

protoc-gen-go  protoc-gen-go-grpc
```

## 项目开发

[项目源码地址](https://github.com/silenceboychen/gostudy/tree/main/demo_4)

### 项目目录结构

```
├── demo_4
│   ├── client
│   │   └── main.go
│   ├── go.mod
│   ├── go.sum
│   ├── helloworld
│   │   ├── helloworld_grpc.pb.go
│   │   ├── helloworld.pb.go
│   │   └── helloworld.proto
│   ├── README.md
│   └── server
│       └── main.go
```

### 项目创建

```bash
$ mkdir demo_4 && cd demo_4
$ go mod init
```

### 安装grpc依赖

```bash
$ go get -u google.golang.org/grpc
```

### 编写proto文件

流式``rpc``使用``stream``关键字定义

```protobuf
# helloworld/helloworld.proto

syntax = "proto3";

option go_package = "/helloworld";
package helloworld;

service Hello {
  rpc SayHello (stream HelloRequest) returns (stream HelloReply) {}
}

message HelloRequest {
  string name = 1;
}

message HelloReply {
  string message = 1;
}
```

### 生成go代码

```bash
$ protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    helloworld/helloworld.proto
```

命令执行成功之后会在helloworld目录下生成两个文件： ``helloworld_grpc.pb.go``和``helloworld.pb.go``，**注意：** 不要手动编辑这两个文件。

### 编写服务端代码

``flag``用法可参考官方文档： https://pkg.go.dev/flag

```go
// server/main.go

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
```

### 编写客户端代码

```go
// client/main.go

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
```

### 项目运行

开启两个终端，分别运行服务端代码和客户端代码，服务端代码要先运行。

**服务端**

```bash
$ go run server/main.go

2023/06/16 20:51:39 server listening at [::]:8080
2023/06/16 20:51:42 client stream: world_0
2023/06/16 20:51:42 client stream: world_1
2023/06/16 20:51:42 client stream: world_2
2023/06/16 20:51:42 client stream: world_3
2023/06/16 20:51:42 client stream: world_4
```

**客户端**

```bash
$ go run client/main.go

2023/06/16 20:51:42 server stream: world_0_0
2023/06/16 20:51:42 server stream: world_1_1
2023/06/16 20:51:42 server stream: world_2_2
2023/06/16 20:51:42 server stream: world_3_3
2023/06/16 20:51:42 server stream: world_4_4
```

下一篇将向大家介绍grpc调试工具。