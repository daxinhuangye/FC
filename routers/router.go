package routers

import (
	"Fcoin/controllers"

	"github.com/astaxie/beego"
)

func init() {

	beego.AutoRouter(&controllers.IndexController{})

}
