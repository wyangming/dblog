package web

import (
	"dblog/controllers"
	"dblog/service/posts"
	"dblog/util/pagination"
	"github.com/astaxie/beego"
)

//前台页面基本的控制类
type WebaseController struct {
	siteid int64
	controllers.BaseController
}

//处理所有访问前端页面之前的信息
func (this *WebaseController) Prepare() {
	//golang不会执行父级里的同名方法，如果想需要显示调用
	//this.BaseController.Prepare()
	//fmt.Println("this is WebaseController Prepare method.")
}

//主页面控制器
type MainController struct {
	WebaseController
}

//首页
func (this *MainController) Get() {
	postService := posts.AutoPostsService()
	filter := make(map[string]interface{})
	filter["post_active"] = 0
	//filter["hasum"] = true
	//filter["hascon"] = true
	nums, posts, err := postService.Find(0, 20, filter)
	if err != nil {
		beego.Error(err)
		return
	}
	//数据库是从0开始计数的
	page := pagination.NewPaginator(1, 20, nums)
	this.Data["posts"] = posts
	this.Data["page"] = page
	this.TplName = "index.html"
}

//网站文章
func (this *MainController) Posts() {
	this.Data["msg"] = "Posts test " + this.GetString(":id")
	this.TplName = "post.html"
}

//网站栏目
func (this *MainController) Term() {
	this.Data["msg"] = "Term test " + this.GetString(":id")
	this.TplName = "index.html"
}
