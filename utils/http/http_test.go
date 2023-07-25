package http

import (
	"fmt"
	"testing"
	"time"
)

func TestHttpGet(t *testing.T) {
	// code, body, err := HttpGetWithTimeout("http://192.168.5.245:8801/v1/kline/index-latest?symbol=BTC_USDC&period=1week", 1 * time.Second)
	code, body, err := HttpGetWithTimeout("https://ordiscan.com/inscription/18212135", 1*time.Second)
	//code, body, err := HttpGet("http://192.168.5.245:8801/v1/kline/index-latest?symbol=BTC_USDC&period=1week")
	fmt.Println(code, body, err)
}
