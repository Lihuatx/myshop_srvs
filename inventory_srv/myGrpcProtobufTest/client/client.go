package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"myshop_srvs/inventory_srv/myGrpcProtobufTest/proto"
)

func main() {
	conn, err := grpc.NewClient("127.0.0.1:8899", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	defer conn.Close()

	c := proto.NewGreeterClient(conn)
	r, err := c.SayHello(context.Background(), &proto.HelloRequest{Name: "汪洋"})
	if err != nil {
		return
	}
	fmt.Println(r.Answer)
}
