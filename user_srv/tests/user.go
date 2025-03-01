package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"myshop_srvs/user_srv/proto"
	"time"
)

var userClient proto.UserClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	userClient = proto.NewUserClient(conn)
}

func TestGetUserList() {
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 2,
	})
	if err != nil {
		panic(err)
	}

	for _, user := range rsp.Data {
		fmt.Println(user.Mobile, user.NickName, user.PassWord)
		checkRsp, err := userClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password:          "admin123",
			EncryptedPassword: user.PassWord,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(checkRsp.Success)
	}
}

func testGetUserByMobile() {
	for ii := 0; ii < 10; ii++ {
		rsp, err := userClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
			Mobile: fmt.Sprintf("1878222222%d", ii),
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(rsp.Mobile, rsp.NickName, rsp.PassWord)
	}
}

func testGetUserById() {
	for ii := 1; ii < 11; ii++ {
		rsp, err := userClient.GetUserById(context.Background(), &proto.IdRequest{
			Id: int32(ii),
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(rsp.Mobile, rsp.NickName, rsp.PassWord)
	}
}

func testCreateUser() {
	rsp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: "yancey10",
		Mobile:   "15819934030",
		PassWord: "123456",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Mobile, rsp.NickName, rsp.PassWord)
}

func testUpdateUser() {
	_, err := userClient.UpdateUser(context.Background(), &proto.UpdateUserInfo{
		Id:       10,
		NickName: "yancey11",
		Gender:   "male",
		BirthDay: uint64(time.Now().Unix()),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Update success!!")
}

func main() {
	Init()

	//TestGetUserList()
	//testGetUserByMobile()
	testGetUserById()
	//testCreateUser()
	//testUpdateUser()

	conn.Close()
}
