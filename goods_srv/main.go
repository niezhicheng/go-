package main

import (
	"awesomeProject8/goods_srv/global"
	"awesomeProject8/goods_srv/handler"
	"awesomeProject8/goods_srv/initialize"
	"awesomeProject8/goods_srv/proto"
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"syscall"
)


func main()  {
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	IP := flag.String("ip","0.0.0.0","ip地址")
	Port := flag.Int("port",50051,"端口")
	fmt.Println("这是",*Port)
	flag.Parse()
	fmt.Println("ip: ",IP)
	server := grpc.NewServer()
	//proto.RegisterUserServer(server,&handler.UserServer{})
	proto.RegisterGoodsServer(server,&handler.GoodsServer{})
	lis,err := net.Listen("tcp",fmt.Sprintf("%s:%d",*IP,*Port))
	if err != nil{
		panic(err)
	}
	//注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server,health.NewServer())
	cfg := api.DefaultConfig()
	//zap.S().Errorw("这是ip地址",global.ServerConfig.ConsulInfo.Host)
	//cfg.Address = fmt.Sprintf("%s:%d",global.ServerConfig.ConsulInfo.Host,global.ServerConfig.ConsulInfo.Port)
	cfg.Address = "172.16.39.133:8500"
	address := fmt.Sprintf("%s:%d",global.ServerConfig.ConsulInfo.Host,global.ServerConfig.ConsulInfo.Port)
	fmt.Println("这是地址",address)
	fmt.Println(global.ServerConfig.ConsulInfo.Port)
	//cfg.Address = address
	fmt.Println(cfg.Address,"这边")
	client,err := api.NewClient(cfg)
	if err != nil{
		panic(err)
	}
	check := &api.AgentServiceCheck{
		GRPC: fmt.Sprintf("172.16.39.133:%d",*Port),
		Timeout: "5s",
		Interval: "5s",
		DeregisterCriticalServiceAfter: "10s",
	}
	//
	//
	regiterrations := new(api.AgentServiceRegistration)
	regiterrations.Name = "goods-srv"
	regiterrations.Port = *Port
	zap.S().Info("启动端口时",*Port)
	regiterrations.ID = fmt.Sprintf("%s",uuid.NewV4())
	//regiterrations.Port = *Port
	regiterrations.Tags = []string{"xaiomi","goods-srv"}
	regiterrations.Address = "172.16.39.133"
	regiterrations.Check = check
	sss := client.Agent().ServiceRegister(regiterrations)
	//client.Agent().ServiceRegister()
	if sss != nil{
		panic(sss)
	}
	//fmt.Println()
	go func() {
		err = server.Serve(lis)
		if err != nil{
			panic(err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit,syscall.SIGINT,syscall.SIGTERM)
	<-quit
	if err = client.Agent().ServiceDeregister(regiterrations.ID);err != nil{
		panic(err)
	}
}