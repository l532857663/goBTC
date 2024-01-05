package models

import (
	"goBTC/db"
	"goBTC/elastic"
	"time"
)

type Server struct {
	Mysql       db.Mysql              `mapstructure:"mysql"        json:"mysql"        yaml:"mysql"`        // 数据库配置
	Zap         Zap                   `mapstructure:"zap"          json:"zap"          yaml:"zap"`          // 日志配置
	Https       Https                 `mapstructure:"https"        json:"https"        yaml:"https"`        // 网络配置
	ElasticConf elastic.ElasticConfig `mapstructure:"elastic_conf" json:"elastic_conf" yaml:"elastic_conf"` // elastic配置
	ChainNode   ChainNode             `mapstructure:"chain_node"   json:"chain_node"   yaml:"chain_node"`   // 链节点配置
	CronTasks   CronTasks             `mapstructure:"cron_tasks"   json:"cron_tasks"   yaml:"cron_tasks"`   // 定时任务
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

type ChainNode struct {
	Ip          string            `mapstructure:"ip"              yaml:"ip"`           // IP地址 或 域名
	Port        int               `mapstructure:"port"            yaml:"port"`         // 端口号
	User        string            `mapstructure:"user"            yaml:"user"`         // 用户名
	Password    string            `mapstructure:"password"        yaml:"password"`     // 密码
	Net         string            `mapstructure:"net"             yaml:"net"`          // 网络类型
	ChainConfig map[string]string `mapstructure:"chain_config"    yaml:"chain_config"` // 链配置
	NodeJsUrl   string            `mapstructure:"node_js_url"      yaml:"node_js_url"` // 请求nodejs链接
}

type CronTasks struct {
	GetPriceService string        `mapstructure:"get_price_service" yaml:"get_price_service"`
	GetTransferInfo time.Duration `mapstructure:"get_transfer_info" yaml:"get_transfer_info"`
}
