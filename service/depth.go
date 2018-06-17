package service

import (
	"Fcoin/models"
	"fmt"
	_ "strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/tidwall/gjson"
)

type PriceData struct {
	Ts   uint64
	Asks float64
	Bids float64
}

var (
	wss         = "wss://api.fcoin.com/v2/ws"
	price_array = []PriceData{}
	price       = float64(0)
	depth       = ""
	sub         = ""
	api         = models.Fcoin{}
)

func DepthRun() {
	//连接服务器
	err := models.WsConn(wss)
	if err != nil {
		go models.SendMail("Fcoin wss无法连接，请及时处理~~")
		beego.Error("连接失败")
		return
	}
	beego.Error("开始监Fcoin")
	//数据监听
	go listionRead()
	go listionOrders()

	//订阅数据
	depth = fmt.Sprintf("depth.L20.%s%s", beego.AppConfig.String("Symbol"), beego.AppConfig.String("Base"))
	sub = fmt.Sprintf(`{"cmd":"sub", "args":["%s"], "id":"1"}`, depth)
	//订阅货币数据
	models.WsSend(wss, sub)
}

//监控订单 每分钟检查一次
func listionOrders() {
	for {
		//获取资金
		//currency, base := api.GetBalance()
		//beego.Trace("币：", currency, " 钱：", base)
		//获取订单列表
		content := api.GetOrders(beego.AppConfig.String("Symbol") + beego.AppConfig.String("Base"))
		data := gjson.Get(content, "data").Array()
		for k, v := range data {

		}
		time.Sleep(60 * time.Second)

	}
}

func reConn() error {
	//先关闭
	models.WsClose(wss)
	//再连接
	err := models.WsConn(wss)
	if err != nil {
		go models.SendMail("Fcoin重新连接失败~~~")
		return err
	}

	models.WsSend(wss, sub)
	return nil
}

//监听数据
func listionRead() {

	defer models.WsClose(wss)

	for {
		//获取数据
		data, err := models.WsRead(wss)
		//如果读取错误，从新连接服务器
		if err != nil {
			price = 0
			beego.Trace("交易暂停，重新连接服务器")
			//重新连接
			err = reConn()
			if err != nil {
				break
			} else {
				continue
			}

			beego.Error(err)
		}

		content := string(data)
		//beego.Trace(content)

		event := gjson.Get(content, "type").String()
		ms := gjson.Get(content, "ts").Uint()

		if event == depth {

			//卖盘
			tick := gjson.Get(content, "asks").Array()
			temp := tick[0].Array()
			asks := temp[0].Float()

			//买盘
			tick = gjson.Get(content, "bids").Array()
			temp = tick[0].Array()
			bids := temp[0].Float()

			//设置交易价格
			setPrice(ms, asks, bids)

		}

	}

}

func setPrice(ms uint64, asks, bids float64) {
	switch beego.AppConfig.String("Mode") {
	case "0":
		price = (asks + bids) / 2
	case "1":
		price = asks
	case "2":
		price = bids
	case "3":
		price = bids
	case "4":
		price = bids

	default:
		ts := ms / 1000
		length := len(price_array)
		if length == 0 {
			price_array = append(price_array, PriceData{Ts: ts, Asks: asks, Bids: bids})
		} else {
			price_data := price_array[length-1]

			if price_data.Ts == ts {
				price_array[length-1] = PriceData{Ts: ts, Asks: asks, Bids: bids}

			} else {
				price_array = append(price_array, PriceData{Ts: ts, Asks: asks, Bids: bids})
			}

		}

		size, _ := beego.AppConfig.Int("Size")
		if length >= size {
			i := 0
			price_array = append(price_array[:i], price_array[i+1:]...) // 最后面的“...”不能省略

		}

		if asks > price_array[0].Asks {
			price = asks
		} else if asks < price_array[0].Asks {
			price = bids
		} else {
			price = (asks + bids) / 2
		}

	}

}
