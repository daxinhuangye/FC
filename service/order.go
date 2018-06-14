package service

import (
	"Fcoin/models"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/tidwall/gjson"
)

type SymbolData struct {
	PriceDecimal  int64
	AmountDecimal int64
}

func OrderRun() {
	api := models.Fcoin{}
	result := api.Init()
	if !result {
		beego.Trace("SDK初始化失败~~")
		return
	}

	symbol_map := make(map[string]SymbolData)
	content := beego.AppConfig.String("SymbolMap")
	temp := gjson.Get(content, "data").Array()
	for _, v := range temp {
		key := v.Get("name").String()
		price_decimal := v.Get("price_decimal").Int()
		amount_decimal := v.Get("amount_decimal").Int()
		symbol_map[key] = SymbolData{PriceDecimal: price_decimal, AmountDecimal: amount_decimal}
	}

	sleep, _ := beego.AppConfig.Int("Sleep")
	amount, _ := beego.AppConfig.Float("Amount")
	symbol := beego.AppConfig.String("Symbol") + beego.AppConfig.String("Base")
	i := uint64(0)
	for {
		if price > 0 {

			currency, base := api.GetBalance()

			beego.Trace(currency, base)

			//买卖的逻辑
			if (currency >= amount) && base >= (amount*price) {

				amount_decimal := symbol_map[symbol].AmountDecimal
				price_decimal := symbol_map[symbol].PriceDecimal
				if amount_decimal > 0 && price_decimal > 0 {
					amount_decimal_str := "%." + fmt.Sprintf("%d", amount_decimal) + "f"
					price_decimal_str := "%." + fmt.Sprintf("%d", price_decimal) + "f"

					go api.CreateOrder(i, fmt.Sprintf(amount_decimal_str, amount), fmt.Sprintf(price_decimal_str, price), "sell", symbol, "limit")
					go api.CreateOrder(i, fmt.Sprintf(amount_decimal_str, amount), fmt.Sprintf(price_decimal_str, price), "buy", symbol, "limit")
					i++
				}

			} else if (currency == 0) && base == 0 {
				beego.Trace("接口错误")
			} else {
				beego.Trace("需要补仓")
			}

		}

		time.Sleep(time.Duration(sleep) * time.Second)

	}
}
