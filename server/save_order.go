package server

import (
	"bytes"
	"goBTC/db/brc20_market"
	"goBTC/global"
	"goBTC/models"
	"goBTC/utils/logutils"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func SaveOrder(txHash string, ordiInfo *models.OrdiInfo) {
	log := global.LOG
	// 保存铭文订单信息
	inscriberID := txHash + "i0"
	order := &brc20_market.Order{
		TxHash:          txHash,
		InscribeID:      &inscriberID,
		InscribeContent: ordiInfo.Body,
		ContentType:     &ordiInfo.ContentType,
		Tick:            ordiInfo.Tick,
		State:           1,
		Number:          ordiInfo.Amount,
		ServerFee:       ordiInfo.ServiceFee,
		GasFee:          &ordiInfo.GasFee,
		GasFeeTotal:     &ordiInfo.GasFeeTotal,
		To:              &ordiInfo.To,
	}
	err := order.Create()
	if err != nil {
		logutils.LogErrorf(log, "Create pending order txhash[%s], orderInfo[%+v] error: %+v", txHash, order, err)
		return
	}
}

func GetInscriptionInfoByOrdinals(order *brc20_market.Order) {
	log := global.LOG
	// 发送GET请求
	resp, err := http.Get("https://ordinals.com/inscription/" + *order.InscribeID) // 替换为你要抓取的URL
	if err != nil {
		logutils.LogErrorf(log, "Error fetching URL: %+v", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logutils.LogErrorf(log, "Error reading response body: %+v", err)
		return
	}
	// fmt.Printf("body: %+v\n", string(bodyBytes))

	// 解析HTML字符串为goquery.Document对象
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bodyBytes))
	if err != nil {
		logutils.LogErrorf(log, "goquery.NewDocumentFromReader error: %+v", err)
		return
	}

	// 遍历 dl 标签内的 dt 和 dd 对
	doc.Find("dl").Each(func(i int, s *goquery.Selection) {
		s.Children().EachWithBreak(func(j int, child *goquery.Selection) bool {
			if child.Is("dt") {
				dtText := strings.TrimSpace(child.Text())
				// 找到目标 dt 后，获取下一个节点（即对应的 dd）
				dd := child.NextFiltered("dd")
				if dd.Length() < 0 {
					return true
				}
				ddContent := strings.TrimSpace(dd.Text())
				switch dtText {
				case "content type":
					order.ContentType = &ddContent
				case "genesis height":
					height, _ := strconv.ParseInt(ddContent, 0, 64)
					order.BlockHeight = &height
				}
			}
			return true // 继续遍历
		})
	})
}
