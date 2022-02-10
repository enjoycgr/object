package main

import (
	"context"
	"flag"
	"log"

	"object/app/objectv2/rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type Config struct {
	zrpc.RpcServerConf
}

var cfgFile = flag.String("f", "./hello.yaml", "cfg file")

func main() {
	flag.Parse()

	var cfg Config
	conf.MustLoad(*cfgFile, &cfg)

	srv, err := zrpc.NewServer(cfg.RpcServerConf, func(s *grpc.Server) {
		pb.RegisterGreeterServer(s, &Hello{})
	})
	if err != nil {
		log.Fatal(err)
	}
	srv.Start()
}

type Hello struct{}

func (h *Hello) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "hello " + in.Name}, nil
}
