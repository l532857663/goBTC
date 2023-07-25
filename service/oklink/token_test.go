package oklink

import (
	"fmt"
	"goBTC/global"
	"testing"
	"utils/logutils"
	"utils/model"
)

func init() {
	// 初始化zap日志库
	zap := model.Zap{
		Level:         "info",
		Format:        "console",
		Prefix:        "[account_info]",
		Director:      "",
		ShowLine:      true,
		EncodeLevel:   "LowercaseColorLevelEncoder",
		StacktraceKey: "stacktrace",
		LogInConsole:  true,
	}
	global.LOG = logutils.Log("", zap)
}

func Test_GetBalanceByAddress(t *testing.T) {
	p := &Platform{}
	p.HttpHeader = map[string]string{
		"Content-Type":  "application/json",
		"Ok-Access-Key": "21140049-abba-4a24-9550-3a37fe4a69c6",
	}
	symbol := "ETH"
	address := "0x5FF539ed0135d01A5bd681DeeF7a8604f343A66f"
	// address := "0xfD4e2816e05E2150a8Ceb9C438Ed78d6317ff746"
	protocolType := "token_20"
	page := "1"
	limit := "10"
	resp, err := p.GetBalanceByAddress(symbol, address, protocolType, page, limit)
	if err != nil {
		t.Errorf("GetBalanceByAddress error: %v", err)
	}
	if resp == nil {
		t.Errorf("GetBalanceByAddress error: response is nil")
	}
	fmt.Printf("wch---- resp: %v\n", resp)
}
