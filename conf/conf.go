package conf

import (
	"fmt"
	"github.com/spf13/viper"
)

type Cfg struct {
	Debug          bool
	Env            string
	LogLevel       string
	LogFile        string
	Path           string
	Frontend       Site
	Backend        Site
	API            API
	MySQL          map[string]MySQL
	Kafka          map[string]Kafka
	Redis          Redis
	MessageChannel string
	SSH            SSH
}

type MySQL struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

type Site struct {
	Host    string
	Port    int
	Logfile string
	Cache   bool
}

type API struct {
	Key         string
	TokenSecret string
	AesKey      string
	AesKey256   string
}

type Kafka struct {
	Topic     string
	Broker    string
	Zookeeper string
}

type Redis struct {
	Addr     string
	Password string
}

type SSH struct {
	Host     string
	Port     int
	User     string
	Password string
}

var Config *Cfg

func Init() (*Cfg, error) {
	return InitByPath("./conf")
}

func InitByPath(path string) (*Cfg, error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("toml")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("can't read config file: %s", err)
	}

	c := Cfg{}
	if err := viper.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("unable to unmarshal config: %s", err)
	}

	c.Path = path
	Config = &c
	return &c, nil
}
