package models

import (
	"goBTC/db"
	"goBTC/elastic"
)

type Server struct {
	Mysql       db.Mysql              `mapstructure:"mysql"        json:"mysql"        yaml:"mysql"`        // 数据库配置
	Zap         Zap                   `mapstructure:"zap"          json:"zap"          yaml:"zap"`          // 日志配置
	Https       Https                 `mapstructure:"https"        json:"https"        yaml:"https"`        // 网络配置
	ElasticConf elastic.ElasticConfig `mapstructure:"elastic_conf" json:"elastic_conf" yaml:"elastic_conf"` // elastic配置
}

type Zap struct {
	Level         string `mapstructure:"level"          json:"level"         yaml:"level"`          // 日志级别
	Format        string `mapstructure:"format"         json:"format"        yaml:"format"`         // 输出方式
	Prefix        string `mapstructure:"prefix"         json:"prefix"        yaml:"prefix"`         // 前缀
	Director      string `mapstructure:"director"       json:"director"      yaml:"director"`       // 目录
	LinkName      string `mapstructure:"link-name"      json:"linkName"      yaml:"link-name"`      // 文件名
	ShowLine      bool   `mapstructure:"show-line"      json:"showLine"      yaml:"showLine"`       // 是否展示行号
	EncodeLevel   string `mapstructure:"encode-level"   json:"encodeLevel"   yaml:"encode-level"`   // 日志编码类型
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktraceKey" yaml:"stacktrace-key"` // 堆栈跟踪
	LogInConsole  bool   `mapstructure:"log-in-console" json:"logInConsole"  yaml:"log-in-console"` // 是否在工作台输出日志
}

type Https struct {
	IsHttps bool   `mapstructure:"is_https" json:"is_https" yaml:"is_https"` // 是否开启HTTPS
	CaCert  string `mapstructure:"ca_cert"  json:"ca_cert"  yaml:"ca_cert"`  // 根证书路径
}
