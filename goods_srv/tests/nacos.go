package main

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"time"
)

func main()  {
	sc := []constant.ServerConfig{
		{
			IpAddr: "127.0.0.1",
			Port: 8848,
		},
	}
	cc := constant.ClientConfig{
		NamespaceId:         "9f8f0f46-f183-425a-bb52-487c48262e1f", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}
	configClient,err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs":sc,
		"clientConfig": cc,
	})
	if err != nil{
		panic(err)
	}
	content,err := configClient.GetConfig(vo.ConfigParam{
		DataId:   "config.json",
		Group:    "DEFAULT_GROUP",
	})
	if err != nil{
		panic(err)
	}
	fmt.Println(content)
	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: "config.json",
		Group:  "DEFAULT_GROUP",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("group" + group)
		},
	})
	if err != nil{
		panic(err)
	}
	time.Sleep( 3000 * time.Second)
}