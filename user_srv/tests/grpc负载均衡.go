package main

import (
	"awesomeProject8/user_srv/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc/credentials/insecure"
)

func main()  {
	conn,err := grpc.Dial("consul://172.16.39.133:8500/user-srv?wait=14s",grpc.WithTransportCredentials(insecure.NewCredentials()),grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	if err != nil{
		panic(err)
	}
	defer conn.Close()
	usersrcclient := proto.NewUserClient(conn)
	for i := 0; i < 10; i++ {
		res,err := usersrcclient.GetUserList(context.Background(),&proto.PageInfo{
			Pn:    1,
			PSize: 10,
		})
		if err != nil{
			panic(err)
		}
		for i, i2 := range res.Data {
			fmt.Println(i,i2.Mobile)
		}
	}
}
