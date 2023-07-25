package http

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
)

// @Description 同步系统地址到redis
// @Author Oracle
// @Version 1.0
// @Update Oracle 2022-06-13 init
func RunPProfHttpServer(port uint64) {
	// 开启对锁调用的跟踪
	runtime.SetMutexProfileFraction(1)
	// 开启对阻塞操作的跟踪
	runtime.SetBlockProfileRate(1)

	go func() {
		var portStr string
		if port != 0 {
			portStr = fmt.Sprintf(":%d", port)
		} else {
			portStr = ":6066"
		}
		// 启动一个 http server，注意 pprof 相关的 handler 已经自动注册过了
		if err := http.ListenAndServe(portStr, nil); err != nil {
			log.Fatal(err)
		}
	}()
}
