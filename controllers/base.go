package controllers

import (
	"tsEngine/tsTime"

	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
	Code    int
	Msg     string
	Result  interface{}
	AdminId int64
}

var langTypes []string // Languages that are supported.

func init() {
	beego.Trace("初始化控制器")

}

func (this *BaseController) Display(tpl string, layout bool) {

	this.Data["Version"] = beego.AppConfig.String("Version")

	if beego.AppConfig.String("runmode") == "dev" {
		this.Data["Version"] = tsTime.CurrSe()
	}

	this.Data["Appname"] = beego.AppConfig.String("AppName")
	this.Data["Website"] = beego.AppConfig.String("WebSite")
	this.Data["Weburl"] = beego.AppConfig.String("WebUrl")
	this.Data["Email"] = beego.AppConfig.String("Email")
	if layout {
		this.Layout = "layout/main.html"
	}
	this.TplName = tpl + ".html"
}

//json 输出
func (this *BaseController) TraceJson() {
	this.Data["json"] = &map[string]interface{}{"Code": this.Code, "Msg": this.Msg, "Data": this.Result}
	this.ServeJSON()
	this.StopRun()
}
