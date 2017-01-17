package posts

import (
	"dblog/controllers/sys"
	"dblog/models"
	"dblog/service/posts"
	"dblog/util/pagination"
	"dblog/util/web"
	"fmt"
	"github.com/astaxie/beego"
)

type PostsControllers struct {
	sys.SysBaseController
}

//跳转到文章管理页面
func (this *PostsControllers) Get() {
	postService := posts.AutoPostsService()
	pageCur, _ := this.GetInt("paged")
	offset := int64(0)
	pageSize, _ := beego.AppConfig.Int64("sys_page_size")
	if pageCur > 1 {
		offset = (int64(pageCur) - 1) * pageSize
	}
	nums, posts, err := postService.Find(offset, pageSize, nil)
	if err != nil {
		beego.Error(err)
		return
	}
	//数据库是从0开始计数的
	page := pagination.NewPaginator(offset+1, pageSize, nums)
	this.Data["posts"] = posts
	this.Data["page"] = page
	this.TplName = "sys/posts/posts.html"
}

//跳转到添加文章页面
func (this *PostsControllers) ToAdd() {
	addPageInfo(this)
	id, _ := this.GetInt64("id")
	if id < 1 {
		return
	}
	postService := posts.AutoPostsService()
	post, err := postService.FindById(id)
	if err != nil {
		beego.Error(err)
	}
	this.Data["post"] = post
}

//跳转到添加页面所需要的信息
func addPageInfo(this *PostsControllers) {
	this.TplName = "sys/posts/posts_to_add.html"
	//所有的栏目信息
	termService := posts.AutoTermService()
	terms, _ := termService.FindAll()
	this.Data["terms"] = terms
	//空文章信息
	postService := posts.AutoPostsService()
	post, _ := postService.FindById(int64(0))
	this.Data["post"] = post
}

//保存一个文章信息
func (this *PostsControllers) Save() {
	//service
	postService := posts.AutoPostsService()
	id, _ := this.GetInt("id")
	posts := &models.DbPosts{}
	action_type := this.Tr("page_btn_update")
	if id < 1 {
		//保存为草稿
		posts.Active = 2
		//当添加时的处理
		posts.CreateUser = web.GetCurUserId(this.Controller)
		posts.SiteId = web.GetSySiteId(this.Controller)
		action_type = this.Tr("page_btn_add")
	} else {
		postView, _ := postService.FindById(int64(id))
		posts = &postView.DbPosts
		posts.Id = int64(id)
	}
	//处理文章信息
	posts.Title = this.GetString("post_title")
	posts.ShotTitle = this.GetString("post_stitle")
	posts.Tags = this.GetString("post_tags")
	posts.SourceUrl = this.GetString("post_source_url")
	posts.Author = this.GetString("post_author")
	posts.Summary = this.GetString("post_summary")
	post_content := this.GetString("post_content")
	posts.HtmlContent = post_content
	posts.TextContent = this.GetString("text_content")
	//保存到处理map里
	save_info := make(map[string]interface{})
	save_info["post"] = posts
	term_id, _ := this.GetInt("term_id")
	save_info["term_id"] = int64(term_id)
	save_info["id"] = int64(id)
	//保存内容
	res, err := postService.SavePost(save_info)
	if err != nil {
		beego.Error(err)
	}
	//处理提示信息
	action_str := this.Tr("page_sys_post")
	if !res {
		this.Data["msg"] = fmt.Sprintf("%s%s%s%s", action_str, posts.Title, action_type, this.Tr("control_sys_result_faile"))
		addPageInfo(this)
		return
	}
	this.Data["msg"] = fmt.Sprintf("%s%s%s%s", action_str, posts.Title, action_type, this.Tr("control_sys_result_success"))
	addPageInfo(this)
}

//发布文章
func (this *PostsControllers) Release() {
	res, err := updateStatus("replase", this)
	if err != nil {
		beego.Error(err)
	}
	if !res {
		this.DbJsonMsg("page_sql_err_msg", false)
		this.ServeJSON()
	}
	this.DbJsonMsg("", true)
	this.ServeJSON()
}

//修改状态
func updateStatus(type_str string, this *PostsControllers) (bool, error) {
	id, _ := this.GetInt("id")
	if id < 1 {
		return false, nil
	}
	postService := posts.AutoPostsService()
	return postService.UpdateStatus(int64(id), type_str)
}

//删除文章
func (this *PostsControllers) Delete() {
	res, err := updateStatus("del", this)
	if err != nil {
		beego.Error(err)
	}
	if !res {
		this.DbJsonMsg("page_sys_post_msg_del_err", false)
		this.ServeJSON()
	}
	this.DbJsonMsg("", true)
	this.ServeJSON()
}
