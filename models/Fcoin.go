package models

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"strconv"
	_ "time"
	"tsEngine/tsCrypto"
	"tsEngine/tsTime"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/tidwall/gjson"
)

type Fcoin struct {
	AccessKey string
	SecretKey string
	DiffTime  uint64
}

const (
	HOST = "https://api.fcoin.com/v2"
)

func (this *Fcoin) Init() bool {
	s_time := this.getTime()
	if s_time == 0 {
		return false
	}
	ts := tsTime.CurrMs()
	this.DiffTime = uint64(this.getTime()) - ts
	return true
}

func (this *Fcoin) GetSymbols() string {
	api := HOST + "/public/symbols"

	curl := httplib.Get(api)

	//设置超时时间 2秒链接，3秒读数据
	//curl.SetTimeout(5*time.Second, 5*time.Second)

	//获取请求的内容
	temp, err := this.request(curl, 3)
	if err != nil {
		return ""
	}
	content := string(temp)
	status := gjson.Get(content, "status").Int()
	if status != 0 {
		return ""
	}

	return content
}

//获取资产
func (this *Fcoin) GetBalance() (float64, float64) {

	api := HOST + "/accounts/balance"
	ts := tsTime.CurrMs() + this.DiffTime

	str := "GET" + api + fmt.Sprintf("%d", ts)
	sg := this.getSha1(tsCrypto.Base64Encode(str))

	curl := httplib.Get(api)
	curl.Header("FC-ACCESS-KEY", beego.AppConfig.String("AccessKey"))
	curl.Header("FC-ACCESS-SIGNATURE", sg)
	curl.Header("FC-ACCESS-TIMESTAMP", fmt.Sprintf("%d", ts))
	//设置超时时间 2秒链接，3秒读数据
	//curl.SetTimeout(5*time.Second, 5*time.Second)

	//获取请求的内容
	temp, err := this.request(curl, 3)
	if err != nil {
		beego.Error(err)
		return 0, 0
	}
	content := string(temp)
	data := gjson.Get(content, "data").Array()

	temp1 := ""
	temp2 := ""
	for _, v := range data {
		currency := v.Get("currency").String()
		if currency == beego.AppConfig.String("Symbol") {
			temp1 = v.Get("available").String()
		}
		if currency == beego.AppConfig.String("Base") {
			temp2 = v.Get("available").String()
		}
	}

	currency, _ := strconv.ParseFloat(temp1, 64)
	base, _ := strconv.ParseFloat(temp2, 64)

	return currency, base
}

//创建订单
func (this *Fcoin) CreateOrder(i uint64, amount, price, side, symbol, types string) {

	if side == "buy" {
		beego.Trace("交易编号：:", i, "买入：", amount, "个 单价：", price)
	} else {
		beego.Trace("交易编号：:", i, "卖出：", amount, "个 单价：", price)
	}

	var postData struct {
		Amount string `json:"amount"`
		Price  string `json:"price"`
		Side   string `json:"side"`
		Type   string `json:"type"`
		Symbol string `json:"symbol"`
	}
	postData.Amount = amount
	postData.Price = price
	postData.Side = side
	postData.Type = types
	postData.Symbol = symbol

	api := HOST + "/orders"
	ts := tsTime.CurrMs() + this.DiffTime

	str := fmt.Sprintf("POST%s%damount=%s&price=%s&side=%s&symbol=%s&type=%s", api, ts, amount, price, side, symbol, types)
	sg := this.getSha1(tsCrypto.Base64Encode(str))

	curl := httplib.Post(api)
	curl.Header("FC-ACCESS-KEY", beego.AppConfig.String("AccessKey"))
	curl.Header("FC-ACCESS-SIGNATURE", sg)
	curl.Header("FC-ACCESS-TIMESTAMP", fmt.Sprintf("%d", ts))
	//设置超时时间 2秒链接，3秒读数据
	curl.JSONBody(postData)

	//curl.SetTimeout(5*time.Second, 5*time.Second)

	//获取请求的内容
	temp, err := this.request(curl, 1)
	if err != nil {
		beego.Trace("交易错误:", err)
		return
	}

	content := string(temp)

	status := gjson.Get(content, "status").Int()

	if status == 0 {
		if side == "buy" {
			beego.Trace("交易成功:", "买入：", amount, "个 单价：", price)
		} else {
			beego.Trace("交易成功:", "卖出：", amount, "个 单价：", price)
		}
	} else {
		beego.Trace(content)
	}

}

//获取订单数
func (this *Fcoin) GetOrders(symbol string) string {

	api := HOST + "/orders?limit=1000&states=submitted&symbol=" + symbol
	ts := tsTime.CurrMs() + this.DiffTime

	str := fmt.Sprintf("GET%s%d", api, ts)
	sg := this.getSha1(tsCrypto.Base64Encode(str))

	curl := httplib.Get(api)
	curl.Header("FC-ACCESS-KEY", beego.AppConfig.String("AccessKey"))
	curl.Header("FC-ACCESS-SIGNATURE", sg)
	curl.Header("FC-ACCESS-TIMESTAMP", fmt.Sprintf("%d", ts))
	//设置超时时间 2秒链接，3秒读数据
	//curl.JSONBody(postData)

	//curl.SetTimeout(5*time.Second, 5*time.Second)

	//获取请求的内容
	temp, err := this.request(curl, 3)
	if err != nil {
		beego.Trace("订单列表接口无法调用:", err)
		return ""
	}

	content := string(temp)

	status := gjson.Get(content, "status").Int()

	if status == 0 {
		return ""
	}
	return content
}

func (this *Fcoin) getSha1(data string) string {

	hmh := hmac.New(sha1.New, []byte(beego.AppConfig.String("SecretKey")))

	hmh.Write([]byte(data))

	hex_data := tsCrypto.Base64Encode(string(hmh.Sum(nil)))

	hmh.Reset()

	return hex_data
}

func (this *Fcoin) getTime() int64 {
	api := HOST + "/public/server-time"

	curl := httplib.Get(api)

	//设置超时时间 2秒链接，3秒读数据
	//curl.SetTimeout(5*time.Second, 5*time.Second)
	temp, err := this.request(curl, 3)
	if err != nil {
		beego.Error(err)
		return 0
	}

	content := string(temp)
	status := gjson.Get(content, "status").Int()
	if status != 0 {
		return 0
	}
	ts := gjson.Get(content, "data").Int()
	return ts
}

func (this *Fcoin) request(curl *httplib.BeegoHTTPRequest, num int) (string, error) {

	temp, err := curl.Bytes()

	if err != nil {
		i := 0
		for i < num {
			temp, err = curl.Bytes()
			if err == nil {
				break
			}
			i++
		}
	}

	if err != nil {
		return "", err
	}

	content := string(temp)
	return content, err
}
