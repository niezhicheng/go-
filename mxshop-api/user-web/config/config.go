package config

type UserSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type TengxunSmsConfig struct {
	Apikey string `mapstructure:"apikey" json:"apikey"`
	ApiSecrect string `mapstructure:"apisecrect" json:"apisecrect"`

}

type RedisConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int `mapstructure:"port" json:"port"`
	Expire int `mapstructure:"expire" json:"expire"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int `mapstructure:"port" json:"port"`
}



type ServerConfig struct {
	Name string `mapstructure:"name" json:"name"`
	Port int `mapstructure:"port" json:"port"`
	UserSrvInfo UserSrvConfig `mapstructure:"user_srv" json:"user_srv"`
	JWTInfo JWTConfig `mapstructure:"jwt" json:"jwt"`
	TengXuninfo TengxunSmsConfig `mapstructure:"tengxuninfo" json:"tengxuninfo"`
	RedisInfo RedisConfig `mapstructure:"redisinfo" json:"redisinfo"`
	ConsulInfo ConsulConfig `mapstructure:"consul" json:"consul"`
}

type NacOsCOnfig struct {
	Nacosinfo NacosConfig `mapstructure:"nacos" json:"nacos"`
}


type NacosConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int `mapstructure:"port" json:"port"`
	NameSpace string `mapstructure:"namespace" json:"namespace"`
	User string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
	DataId string `mapstructure:"dataid" json:"dataid"`
	Group string `mapstructure:"group" json:"group"`
}

