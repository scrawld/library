package config

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Basic struct {
		Debug bool   `yaml:"debug"`
		Port  string `yaml:"port"`
	} `yaml:"basic"`

	DB map[string]struct {
		Master Mysql `yaml:"master"`
		Slave  Mysql `yaml:"slave"`
	} `yaml:"db"`

	Redis struct {
		Addr         string `yaml:"addr"`         // 服务器地址:端口
		Username     string `yaml:"username"`     // 用户名
		Password     string `yaml:"password"`     // 密码
		DB           int    `yaml:"db"`           // redis数据库
		TlsProtocols bool   `yaml:"tlsProtocols"` // tls是否启动
	} `yaml:"redis"`

	Rabbitmq struct {
		Host         string `yaml:"host"`
		Port         string `yaml:"port"`
		Username     string `yaml:"username"`
		Password     string `yaml:"password"`
		Vhost        string `yaml:"vhost"`
		TlsProtocols bool   `yaml:"tlsProtocols"`
	} `yaml:"rabbitmq"`

	// 内部服务地址配置
	InnerServer map[string]string `yaml:"inner-server"`
}

// mysql
type Mysql struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbName"`
}

var config = &ServerConfig{}

func Get() *ServerConfig {
	return config
}

// ReadInConfig 读配置文件
func ReadInConfig(path string) error {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("viper read error: %s", err)
	}
	v.WatchConfig() // 监视配置文件修改

	// 配置文件更新时
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(config); err != nil {
			log.Printf("config change: unmarshal error: %s", err)
		}
	})
	if err := v.Unmarshal(config); err != nil {
		return fmt.Errorf("unmarshal error: %s", err)
	}
	return nil
}
