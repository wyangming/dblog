package main

//show processlist;可以显示当前用户的链接信息mysql

//待完美的功能
//栏目添加修改时没有用js验证

import (
	"dblog/controllers"
	"dblog/controllers/sys"
	"dblog/controllers/sys/posts"
	_ "dblog/dao"
	_ "dblog/routers"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

func init() {
	fmt.Println("dblog begin start...")
}

func main() {
	i18n_reg()
	template_reg()
	router_reg()
	filter_reg()
	beego.Run()
	fmt.Println("dblog over start...")
}

//链接注册
func router_reg() {
	fmt.Println("router regin begin...")
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

	fmt.Println("router regin over...")
}

//过滤器
func filter_reg() {
	fmt.Println("filter regin begin...")
	//后台过滤器，检测是否有权限
	beego.InsertFilter("/sys", beego.BeforeRouter, sys.SysAuthor)
	beego.InsertFilter("/sys/**", beego.BeforeRouter, sys.SysAuthor)
	fmt.Println("filter regin over...")
}

//国际化
func i18n_reg() {
	fmt.Println("i18n regin begin...")
	//使用国际化时，需要导入beego的i18n包：go get github.com/beego/i18n
	//注册国际化
	i18n.SetMessage("zh-CN", "conf/locale_zh-CN.ini")
	fmt.Println("i18n regin over...")
}

//模板注册
func template_reg() {
	fmt.Println("template regin begin...")
	//经过注册后就可以在模版里使用i18n国际化
	beego.AddFuncMap("i18n", i18n.Tr)

	//测试用的的模板
	//beego.AddFuncMap("my_TmpMethod", TmpMethod)
	fmt.Println("template regin over...")
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
