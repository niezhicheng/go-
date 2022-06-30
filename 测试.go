package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func main()  {
	v := viper.New()
	v.SetConfigFile("config.yaml")
	v.ReadInConfig()
	fmt.Println(v.Get("name"))
}
