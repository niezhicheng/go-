package main

import (
	"awesomeProject8/user_srv/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
)

var userClient proto.UserClient
var conn *grpc.ClientConn
func Init()  {
	var err error
	conn,err = grpc.Dial("127.0.0.1:50051",grpc.WithInsecure())
	if err != nil{
		panic(err)
	}
	userClient = proto.NewUserClient(conn)
}

func TextGetUserList()  {
	var err error
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 10,
	})
	if err != nil{
		panic(err)
	}
	for i, i2 := range rsp.Data{
		fmt.Println(i,i2.NickName)
		checkrsp,err := userClient.CheckPassword(context.Background(),&proto.PasswordCheckInfo{
			Password: "admin123",
			EncryptedPassword: i2.Password,
		})
		if err != nil{
			panic(err)
		}
		fmt.Println(checkrsp.Success)
	}


	//defer conn.Close()
	//c := proto.NewUserClient(conn)
	//r,err := c.GetUserList(context.Background(),&proto.PageInfo{
	//	Pn:    1,
	//	PSize: 10,
	//})
	//if err != nil{
	//	panic(err)
	//}
	////fmt.Println(r.Data)
	//for i, datum := range r.Data {
	//	fmt.Println(i,datum)
	//}
}

func main()  {
	Init()
	TextGetUserList()
	defer conn.Close()
}