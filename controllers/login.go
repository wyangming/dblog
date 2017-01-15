package controllers

import (
	"dblog/service/user"
	"dblog/util/web"
	"strings"
)

type LoginController struct {
	BaseController
}

//用户登录
func (this *LoginController) Get() {
	cur_userId := web.GetCurUserId(this.Controller)
	sys_siteId := web.GetSySiteId(this.Controller)
	if cur_userId > 0 && sys_siteId > 0 {
		this.Ctx.Redirect(302, "/sys")
	}
	this.TplName = "login.html"
}

//用户登录认证
func (this *LoginController) Post() {
	this.TplName = "login.html"
	username := this.Input().Get("username")
	pwd := this.Input().Get("pwd")
	if len(strings.Trim(username, " ")) < 1 {
		this.Data["msg"] = this.Tr("control_login_msg_user_nil")
		return
	}
	if len(strings.Trim(pwd, " ")) < 1 {
		this.Data["msg"] = this.Tr("control_login_msg_pwd_nil")
		return
	}
	userDao := user.AutoUserService()
	user, login_result, err := userDao.Login(username, pwd)
	if !login_result || err != nil {
		this.Data["msg"] = this.Tr("control_login_msg")
		return
	}
	web.SetCurUserId(this.Controller, user.Id)
	web.SetSySiteId(this.Controller, user.SiteId)
	this.Ctx.Redirect(302, "/sys")
}
