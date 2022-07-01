package config

type MysqlConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int `mapstructrue:"port" json:"port"`
	Name string `mapstructrue:"db" json:"name"`
	User string `mapstructrue:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int `mapstructrue:"port" json:"port"`
}

type ServerConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
	MysqlInfo MysqlConfig `mapstructure:"mysqlinfo" json:"mysqlinfo"`
	ConsulInfo ConsulConfig `mapstructure:"consulinfo" json:"consulinfo"`
}


