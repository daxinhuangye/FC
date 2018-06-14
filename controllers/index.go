package controllers

import (
	"github.com/astaxie/beego"
)

type IndexController struct {
	BaseController
}

//类似构造函数
func (this *IndexController) Prepare() {

}

//默认网站首页
func (this *IndexController) Get() {
	beego.Trace("test")
	this.Display("index", false)
}
