package client

import (
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/minchenzz/brc20tool/pkg/btcapi/mempool"
)

type BTCClient struct {
	Client        *rpcclient.Client
	Params        *chaincfg.Params
	MempoolClient *mempool.MempoolClient
}

type Node struct {
	Ip       string
	Port     int
	User     string
	Password string
	Net      string
}

func NewBTCClient(conf *Node) (*BTCClient, error) {
	var (
		url     string
		isHttps bool
	)
	if conf.Port != 0 {
		url = fmt.Sprintf("%s:%d", conf.Ip, conf.Port)
	} else {
		url = conf.Ip
	}

	// 某些https节点配置需要做一些特殊处理
	if strings.HasPrefix(url, "https://") {
		isHttps = true
		url = strings.TrimPrefix(url, "https://")
	}

	connCfg := &rpcclient.ConnConfig{
		Host:         url,
		User:         conf.User,
		Pass:         conf.Password,
		HTTPPostMode: true,     // Bitcoin core only supports HTTP POST mode
		DisableTLS:   !isHttps, // Bitcoin core does not provide TLS by default
	}
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, err
	}

	btcClient := &BTCClient{
		Client: client,
	}
	upperNet := strings.ToUpper(conf.Net)
	switch upperNet {
	case MainNet:
		btcClient.Params = &chaincfg.MainNetParams
	case TestNet:
		btcClient.Params = &chaincfg.TestNet3Params
	case RegNet:
		btcClient.Params = &chaincfg.RegressionNetParams
	default:
		btcClient.Params = &chaincfg.Params{}
	}

	btcClient.MempoolClient = mempool.NewClient(btcClient.Params)

	return btcClient, nil
}
