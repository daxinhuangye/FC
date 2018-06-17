package main

import (
	_ "Fcoin/models"
	_ "Fcoin/routers"
	"Fcoin/service"
	_ "fmt"
	"net"

	"tsEngine/tsCrypto"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

//默认启动
func main() {

	//log记录设置
	beego.SetLogger("file", `{"filename":"./logs/logs.log"}`)

	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	} else {
		beego.SetLevel(beego.LevelInformational)
	}
	beego.SetLogFuncCall(true)

	//验证使用权
	check := false

	mac := beego.AppConfig.String("Mac")

	interfaces, err := net.Interfaces()
	if err != nil {
		beego.Trace("获取mac错误")
	}

	for _, inter := range interfaces {

		temp := inter.HardwareAddr.String() + "4394363415800817141"
		mmac := tsCrypto.GetMd5([]byte(temp))
		if mac == mmac {
			check = true
			break
		}
	}
	if !check {
		beego.Trace("无权使用，请联系管理员")
		return
	}

	//启动行情服务
	go service.DepthRun()
	go service.OrderRun()

	//框架启动
	beego.Run()

}
