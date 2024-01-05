package main

import (
	"fmt"
	"goBTC"
	"goBTC/global"
	"goBTC/tasks"
	"goBTC/utils"
)

func main() {
	fmt.Println("vim-go")
	global.MysqlFlag = true
	goBTC.MustLoad("./config.yml")
	tasks.CronGetRate()
	if global.MysqlFlag {
		utils.SignalHandler("orderCheck", goBTC.Shutdown)
	}
}
