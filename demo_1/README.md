# grpc系列课程（一）：单项rpc

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

[项目源码地址](https://github.com/silenceboychen/gostudy/tree/main/demo_1)

### 项目目录结构

```
├── demo_1
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
$ mkdir demo_1 && cd demo_1
$ go mod init
```

### 安装grpc依赖

```bash
$ go get -u google.golang.org/grpc
```

### 编写proto文件

```protobuf
# helloworld/helloworld.proto

syntax = "proto3";

option go_package = "/helloworld";
package helloworld;

service Hello {
  rpc SayHello (HelloRequest) returns (HelloReply) {}
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
	"context"
	"flag"
	"fmt"
	"github.com/silenceboychen/gostudy/demo_1/helloworld"
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

func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &helloworld.HelloReply{Message: "Hello " + in.GetName()}, nil
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
	"github.com/silenceboychen/gostudy/demo_1/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	r, err := c.SayHello(ctx, &helloworld.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
```

### 项目运行

开启两个终端，分别运行服务端代码和客户端代码，服务端代码要先运行。

**服务端**

```bash
$ go run server/main.go

2023/06/16 18:56:57 server listening at [::]:8080
2023/06/16 18:57:02 Received: world
```

**客户端**

```bash
$ go run client/main.go

2023/06/16 18:57:02 Greeting: Hello world
```

以上便是go实现单项rpc的所有内容，如果一切顺利，恭喜你，以及成功入门grpc，下一篇将向大家介绍服务端流式rpc。