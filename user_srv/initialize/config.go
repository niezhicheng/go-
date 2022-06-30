package initialize

import (
	"awesomeProject8/user_srv/global"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func InitConfig()  {
	v := viper.New()
	path := "./config.yaml"
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil{
		panic(err)
	}
	if err := v.Unmarshal(&global.ServerConfig); err != nil{
		panic(err)
	}
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		_ = v.ReadInConfig()
		_ = v.Unmarshal(&global.ServerConfig)
	})

}
