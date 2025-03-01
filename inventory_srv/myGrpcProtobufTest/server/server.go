package main

import (
	"context"
	"google.golang.org/grpc"
	"myshop_srvs/inventory_srv/myGrpcProtobufTest/proto"
	"net"
)

type Server struct {
	proto.UnimplementedGreeterServer
}

func (s Server) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{
		Answer: "GetYourName:" + req.Name,
	}, nil
}

func main() {
	g := grpc.NewServer()
	proto.RegisterGreeterServer(g, &Server{})
	lis, err := net.Listen("tcp", "0.0.0.0:8899")
	if err != nil {
		panic(err)
	}
	err = g.Serve(lis)
	if err != nil {
		panic(err)
	}
}
