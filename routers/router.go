package routers

import (
	"dblog/controllers/web"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &web.MainController{})
}
