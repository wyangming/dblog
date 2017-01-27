package main

//show processlist;可以显示当前用户的链接信息mysql

//待完美的功能
//栏目添加修改时没有用js验证
//文章添加修改时没有用js验证，预览功能也没有做

//

import (
	"dblog/controllers"
	"dblog/controllers/sys"
	"dblog/controllers/sys/options"
	"dblog/controllers/sys/posts"
	"dblog/controllers/web"
	_ "dblog/dao"
	_ "dblog/routers"
	//"fmt"
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

func init() {
}

func main() {
	i18n_reg()
	template_reg()
	router_reg()
	filter_reg()
	beego.Run()
}

//链接注册
func router_reg() {
	main := &web.MainController{}
	//文章页面信息
	beego.Router("/:id([0-9]+).html", main, "*:Posts")
	//栏目页面
	beego.Router("/:id", main, "*:Term")

	//登录页面
	beego.Router("/login", &controllers.LoginController{})
	//后台主页面
	beego.Router("/sys", &sys.InitController{})

	//使用iframe时用这个
	//beego.Router("/sys/load", &sys.InitController{}, "*:Load")

	//栏目管理
	termCon := &posts.TermControllers{}
	//跳转到栏目管理页面
	beego.Router("/sys/posts/term", termCon)
	//跳转到添加栏目页面
	beego.Router("/sys/posts/term/toadd", termCon, "*:ToAdd")
	//保存栏目信息
	beego.Router("/sys/posts/term/save", termCon, "*:Save")
	//删除一个栏目信息
	beego.Router("/sys/posts/term/remove", termCon, "*:Remove")

	//文章管理
	postsCon := &posts.PostsControllers{}
	//跳转到文章管理页面
	beego.Router("/sys/posts/posts", postsCon)
	//跳转到文章添加页面
	beego.Router("/sys/posts/posts/toadd", postsCon, "*:ToAdd")
	//保存文章信息
	beego.Router("/sys/posts/posts/save", postsCon, "*:Save")
	//文章发布
	beego.Router("/sys/posts/posts/release", postsCon, "*:Release")
	//删除文章
	beego.Router("/sys/posts/posts/del", postsCon, "*:Delete")

	//系统设置
	siteCon := &options.SiteControllers{}
	//跳转到站点设置页面
	beego.Router("/sys/options/site", siteCon)
	beego.Router("/sys/options/site/save", siteCon, "*:Save")
}

//过滤器
func filter_reg() {
	//后台过滤器，检测是否有权限
	beego.InsertFilter("/sys", beego.BeforeRouter, sys.SysAuthor)
	beego.InsertFilter("/sys/**", beego.BeforeRouter, sys.SysAuthor)
}

//国际化
func i18n_reg() {
	//使用国际化时，需要导入beego的i18n包：go get github.com/beego/i18n
	//注册国际化
	i18n.SetMessage("zh-CN", "conf/locale_zh-CN.ini")
}

//模板注册
func template_reg() {
	//经过注册后就可以在模版里使用i18n国际化
	beego.AddFuncMap("i18n", i18n.Tr)

	//测试用的的模板
	//beego.AddFuncMap("my_TmpMethod", TmpMethod)
}

/*

func TmpMethod(p ...interface{}) interface{} {
	fmt.Println(p)
	maps := make(map[string]interface{})
	maps["name"] = "dubing"
	maps["age"] = 27
	return maps
}
*/
