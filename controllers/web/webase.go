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
	findPosts(this, true)
	pubInfo(this)
	this.TplName = "index.html"
}
func pubInfo(this *MainController) {
	postService := posts.AutoPostsService()
	//文章总数
	archsNum := postService.PostsNum()
	this.Data["postsNum"] = archsNum
	//栏目处理
	termService := posts.AutoTermService()
	nums, terms, err := termService.Find(0, 100, nil)
	if err != nil {
		beego.Error(err)
	}
	this.Data["terms"] = terms
	this.Data["termNum"] = nums
}
func findPosts(this *MainController, isIndex bool) {
	//首页文章处理
	postService := posts.AutoPostsService()
	filter := make(map[string]interface{})
	filter["post_active"] = 0
	if isIndex {
		filter["hasum"] = true
	} else {
		slug := this.GetString(":id")
		filter["hasum"] = true
		filter["slug"] = slug
	}
	//filter["hascon"] = true
	pageCur, _ := this.GetInt("paged")
	offset := int64(0)
	pageSize := int64(10)
	if pageCur > 1 {
		offset = (int64(pageCur) - 1) * pageSize
	}
	nums, archs, err := postService.Find(offset, pageSize, filter)
	//数据库是从0开始计数的
	page := pagination.NewPaginator(offset+1, pageSize, nums)
	if err != nil {
		beego.Error(err)
		return
	}
	this.Data["posts"] = archs
	this.Data["page"] = page
}

//网站文章
func (this *MainController) Posts() {
	this.TplName = "post.html"
	pubInfo(this)
	id, err := this.GetInt64(":id")
	if err != nil {
		beego.Error(err)
		return
	}
	postService := posts.AutoPostsService()
	post, err := postService.FindById(id)
	if err != nil {
		beego.Error(err)
		return
	}
	this.Data["post"] = post
}

//网站栏目
func (this *MainController) Term() {
	this.TplName = "term.html"
	pubInfo(this)
	findPosts(this, false)
	sulg := this.GetString(":id")
	//栏目处理
	termService := posts.AutoTermService()
	term, err := termService.FindBySlug(sulg)
	if err != nil {
		beego.Error(err)
	}
	this.Data["term"] = term
}
