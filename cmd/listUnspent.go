package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/rpcclient"
)

func main() {
	// 配置RPC连接
	connCfg := &rpcclient.ConnConfig{
		// Host:         "167.235.193.148:8332",
		Host:         "128.140.73.158:18443",
		User:         "btc",
		Pass:         "btc2021",
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	// 创建RPC客户端
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatalf("Error creating new RPC client: %v", err)
	}
	defer client.Shutdown()

	fmt.Printf("test1\n")
	// 指定要扫描的地址
	// address := "bc1pedhx8x04v487ed4qwshnl4fnz6dyj363yfa7xs85qj2sruprp7lqmnjvly"
	address := "tb1pg0uc7ujx6rplw4wj73etg505jh49k63s7wc3kyngf73ze7ffue4skru6ld"

	// 构建scantxoutset请求参数
	scanParams := []interface{}{
		"start",
		[]string{fmt.Sprintf("addr(%s)", address)},
	}

	// 将参数转换为[]json.RawMessage
	rawParams := make([]json.RawMessage, len(scanParams))
	for i, param := range scanParams {
		rawParam, err := json.Marshal(param)
		if err != nil {
			log.Fatalf("Failed to marshal param: %v", err)
		}
		rawParams[i] = rawParam
	}

	// 调用scantxoutset RPC方法
	result, err := client.RawRequest("scantxoutset", rawParams)
	if err != nil {
		log.Fatalf("Failed to call scantxoutset: %v", err)
	}

	// 输出结果
	fmt.Printf("Scantxoutset result: %s\n", result)
}
