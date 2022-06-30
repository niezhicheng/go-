package initalize

import (
	"awesomeProject8/mxshop-api/user-web/global"
	"encoding/json"
	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitConfig()  {
	v := viper.New()
	//s, _ := os.Getwd()
	path := "./config.yaml"
	//fmt.Println(path)
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil{
	}
	//zap.S().Info("这是",global.ServerConfig.UserSrvInfo)
	if err := v.Unmarshal(global.NacosConfig); err != nil{
		panic(err)
	}
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.NacosConfig)
	})


	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Nacosinfo.Host,
			Port: uint64(global.NacosConfig.Nacosinfo.Port),
		},
	}
	cc := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Nacosinfo.NameSpace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
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
		DataId:   global.NacosConfig.Nacosinfo.DataId,
		Group:    global.NacosConfig.Nacosinfo.Group,
	})
	if err != nil{
		panic(err)
	}
	json.Unmarshal([]byte(content),&global.ServerConfig)
	zap.S().Info(global.ServerConfig.ConsulInfo.Host)
	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: global.NacosConfig.Nacosinfo.DataId,
		Group:  global.NacosConfig.Nacosinfo.Group,
		OnChange: func(namespace, group, dataId, data string) {
		},
	})
	if err != nil{
		panic(err)
	}
}
