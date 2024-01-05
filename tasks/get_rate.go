package tasks

import (
	"encoding/json"
	"goBTC/db/brc20_market"
	"goBTC/global"
	"goBTC/utils/http"
	"goBTC/utils/logutils"

	"github.com/robfig/cron/v3"
)

type TasksService struct {
	isRunning bool
	cronCtx   *cron.Cron
}

// @Description 定时获取汇率存储数据库
// Joker 2024-01-03 init
func CronGetRate() {
	tasksCtx := &TasksService{
		isRunning: true,
		cronCtx:   cron.New(cron.WithSeconds()),
	}

	// 保存汇率
	getPriceCron := global.CONFIG.CronTasks.GetPriceService
	_, err := tasksCtx.cronCtx.AddFunc(getPriceCron, func() {
		tasksCtx.SaveRate4DB()
	})
	if err != nil {
		logutils.LogErrorf(global.LOG, "Add task [CronGetRate] error: %+v", err)
		return
	}

	// 启动cron任务控制
	tasksCtx.cronCtx.Start()
	logutils.LogInfof(global.LOG, "Register CronGetRate task success")
}

// Joker 2024-01-03 init
func (t *TasksService) SaveRate4DB() {
	// 查询汇率
	rate := GetRate()
	if rate == nil {
		logutils.LogErrorf(global.LOG, "SaveRate4DB not get rate")
		return
	}
	// 保存汇率
	r := brc20_market.Rate{
		Pair:  rate.Symbol,
		Price: rate.Price,
	}
	err := r.SaveRate()
	if err != nil {
		logutils.LogErrorf(global.LOG, "SaveRate4DB SaveRate error: %+v", err)
		return
	}
}

var (
	BnbRateHost = []string{
		"https://api.binance.com",
		"https://api1.binance.com",
		"https://data-api.binance.vision",
	}
	BnbRateUri = "/api/v3/ticker/price?symbol=BTCUSDT"
)

type RateResp struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func GetRate() *RateResp {
	var r *RateResp
	for _, host := range BnbRateHost {
		theUrl := host + BnbRateUri
		code, body, err := http.HttpWithTimeoutDisableKeepAlives(theUrl, "GET")
		if code != 200 || err != nil {
			logutils.LogErrorf(global.LOG, "GetRate [%s] ask http error: %+v", theUrl, err)
			continue
		}
		logutils.LogInfof(global.LOG, "code: %+v, body: %+v", code, string(body))
		err = json.Unmarshal(body, &r)
		if err != nil {
			logutils.LogErrorf(global.LOG, "GetRate [%s] json unmarshal error: %+v", theUrl, err)
			continue
		}
		break
	}
	return r
}
