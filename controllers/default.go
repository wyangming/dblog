package controllers

import (
//"github.com/astaxie/beego"
)

type MainController struct {
	BaseController
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.html"
}

//网站文章
func (this *MainController) Posts() {
	this.Data["msg"] = "Posts test " + this.GetString(":id")
	this.TplName = "index.html"
}

//网站栏目
func (this *MainController) Term() {
	this.Data["msg"] = "Term test " + this.GetString(":id")
	this.TplName = "index.html"
}

func (this *MainController) TmpMethod() string {
	return "sdfsdf-aaa"
}
