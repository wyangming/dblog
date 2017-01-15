package sys

import (
	"dblog/controllers"
	"dblog/util/web"
	"github.com/astaxie/beego/context"
)

var (
	//后台认证
	SysAuthor = func(ctx *context.Context) {
		cur_userid, _ := ctx.Input.Session(web.CURR_USERID).(int64)
		if cur_userid < 1 && ctx.Request.RequestURI != web.LOGIN_PAGE_URL {
			ctx.Redirect(302, "/login")
		}
		sys_siteid, _ := ctx.Input.Session(web.SYS_SITEID).(int64)
		if sys_siteid < 1 && ctx.Request.RequestURI != web.LOGIN_PAGE_URL {
			ctx.Redirect(302, "/login")
		}
	}
)

type SysBaseController struct {
	controllers.BaseController
}

//这里只有msg提示信息，跟res返回结果两种
func (this *SysBaseController) DbJsonMsg(i18n_str string, bolean bool) {
	res := make(map[string]interface{})
	res["msg"] = this.Tr(i18n_str)
	res["res"] = bolean
	this.Data["json"] = res
}
