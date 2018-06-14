package main

import (
	_ "Fcoin/models"
	_ "Fcoin/routers"
	"Fcoin/service"

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

	//启动行情服务
	go service.DepthRun()
	go service.OrderRun()

	//框架启动
	beego.Run()
}
