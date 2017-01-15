package web

import (
	"github.com/astaxie/beego"
	"net/url"
)

const (
	//登录页面路径
	LOGIN_PAGE_URL string = "login"

	//sys后台使用的变量
	//当前登录人编号
	CURR_USERID string = "curr_userid"
	//当前登录人管理的网站编号
	SYS_SITEID string = "sys_siteid"
)

//内部的方法
func retInt64(control beego.Controller, key string) int64 {
	userid := control.GetSession(key)
	str, ok := userid.(int64)
	if ok {
		return str
	}
	return 0
}

//外部的方法
//后台当前登录人的Id
func SetCurUserId(control beego.Controller, userid int64) {
	control.SetSession(CURR_USERID, userid)
}
func GetCurUserId(control beego.Controller) int64 {
	return retInt64(control, CURR_USERID)
}

//后台站点的编号
func SetSySiteId(control beego.Controller, siteId int64) {
	control.SetSession(SYS_SITEID, siteId)
}
func GetSySiteId(control beego.Controller) int64 {
	return retInt64(control, SYS_SITEID)
}

//把当前的url存到数到this.Data["DURL"]
func SetUrl(control beego.Controller) {
	link, _ := url.ParseRequestURI(control.Ctx.Request.URL.String())
	control.Data["DURL"] = link
}
