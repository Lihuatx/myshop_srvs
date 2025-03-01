package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"myshop_srvs/smartChat_srv/handler"
	"myshop_srvs/smartChat_srv/proto"
	"regexp"
)

var chat proto.ChatServiceClient
var conn *grpc.ClientConn

func TestchatserverProcessmessage() {
	ctx := context.Background()
	req := &proto.ChatRequest{
		UserId:   1,
		Question: "我的订单号5现在状态是什么？",
	}
	rsp, err := chat.ProcessMessage(ctx, req)
	if err != nil {
		fmt.Printf("ProcessMessage failed: %v", err)
	}
	if rsp.Answer == "" {
		fmt.Printf("ProcessMessage failed: answer is empty. Request: %+v", req)
	} else {
		fmt.Printf("ProcessMessage success: %+v", rsp)
	}
}

func Test_classifyByKeywords() {
	testCases := []struct {
		question string
		want     string
	}{
		{"我的订单号5现在状态是什么？", "order_query"},
		{"退货", "refund"},
		{"物流信息", "logistics"},
		{"投诉", "complaint"},
		{"其他问题", "other"},
	}
	for _, tc := range testCases {
		got := handler.ClassifyByKeywords(tc.question)
		if got != tc.want {
			fmt.Printf("ClassifyByKeywords(%s) = %s, want %s", tc.question, got, tc.want)
		}

	}
}

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:59671", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	chat = proto.NewChatServiceClient(conn)
}

var orderSNRegex = regexp.MustCompile(`([A-Za-z0-9]{1,10})`)

func ExtractOrderSN(question string) string {
	matches := orderSNRegex.FindStringSubmatch(question)
	if len(matches) > 0 {
		return matches[0]
	}
	return ""
}

func main() {
	Init()
	//Test_classifyByKeywords()
	TestchatserverProcessmessage()

}
