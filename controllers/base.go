package controllers

import (
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

type BaseController struct {
	beego.Controller
	i18n.Locale
}

func (this *BaseController) Prepare() {
	//一般用浏览器header里Accept-Language来判断地区
	lang := this.GetString("lang")
	if len(lang) < 1 {
		lang = "zh-CN"
	} else {
		lang = "zh-CN"
	}
	this.Lang = lang
	//为了让i18n在棋牌中使用，存储一下地区信息
	this.Data["Lang"] = lang
}
