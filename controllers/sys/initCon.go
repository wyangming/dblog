package sys

import ()

type InitController struct {
	SysBaseController
}

//跳转到主贾面
func (this *InitController) Get() {
	this.TplName = "sys/main.html"
}

//跳转到后台初始化页面
func (this *InitController) Load() {
	this.TplName = "sys/load.html"
}
