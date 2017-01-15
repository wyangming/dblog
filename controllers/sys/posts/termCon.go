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

type TermControllers struct {
	sys.SysBaseController
}

//跳转到栏目管理页面
func (this *TermControllers) Get() {
	this.TplName = "sys/posts/term.html"
	postService := posts.AutoTermService()
	pageCur, _ := this.GetInt("paged")
	offset := int64(0)
	pageSize, _ := beego.AppConfig.Int64("sys_page_size")
	if pageCur > 1 {
		offset = (int64(pageCur) - 1) * pageSize
	}
	nums, terms, err := postService.Find(offset, pageSize, nil)
	if err != nil {
		beego.Error(err)
		return
	}
	//数据库是从0开始计数的
	page := pagination.NewPaginator(offset+1, pageSize, nums)
	this.Data["terms"] = terms
	this.Data["page"] = page
}

//跳转到栏目添加页面
func (this *TermControllers) ToAdd() {
	this.TplName = "sys/posts/term_to_add.html"
	id, _ := this.GetInt64("id")
	if id < 1 {
		return
	}
	postService := posts.AutoTermService()
	term, err := postService.FindById(id)
	if err != nil {
		this.Data["msg"] = this.Tr("page_sql_err_msg")
		beego.Error(err)
		return
	}
	this.Data["term"] = term
}

//保存栏目
func (this *TermControllers) Save() {
	this.TplName = "sys/posts/term_to_add.html"
	//栏目名称
	term_name := this.GetString("term_name")
	if len(term_name) < 1 {
		this.Data["msg"] = this.Tr("control_sys_term_add_msg_name_nil")
		return
	}
	//栏目访问url
	term_url := this.GetString("term_url")
	if len(term_url) < 1 {
		this.Data["msg"] = this.Tr("control_sys_term_add_msg_url_nil")
		return
	}
	postService := posts.AutoTermService()
	//栏目标签
	term_tags := this.GetString("term_tags")
	//栏目描述
	term_des := this.Input().Get("term_des")
	terms := &models.DbTerms{}
	//栏目编号
	term_id, _ := this.GetInt64("id")
	//操作类型
	action_type := this.Tr("page_btn_add")
	if term_id > 0 {
		view_terms, _ := postService.FindById(term_id)
		*terms = view_terms.DbTerms
		action_type = this.Tr("page_btn_save")
	} else {
		terms = &models.DbTerms{}
		//当添加时的处理
		terms.CreateUser = web.GetCurUserId(this.Controller)
		terms.SiteId = web.GetSySiteId(this.Controller)
	}
	terms.TermName = term_name
	terms.TermTag = term_tags
	terms.Description = term_des
	terms.Slug = term_url
	//父级栏目
	p_term_id, _ := this.GetInt64("p_term_id")
	if p_term_id > 0 {
		terms.ParentId = p_term_id
	}
	//检查url是否重复
	if url_bool, _ := postService.RepeadByUrl(terms.Slug, terms.SiteId, terms.Id); url_bool {
		this.Data["msg"] = fmt.Sprintf("%s%s%s", this.Tr("page_sys_term_add_lable_slug"),
			term_url, this.Tr("control_sys_result_repeat"))
		return
	}
	err := postService.SaveTerm(terms)
	if err != nil {
		beego.Error(err)
		this.Data["msg"] = fmt.Sprintf("%s%s%s%s", this.Tr("page_sys_term"), term_name, action_type, this.Tr("control_sys_result_faile"))
		return
	}
	this.Data["msg"] = fmt.Sprintf("%s%s%s%s", this.Tr("page_sys_term"), term_name, action_type, this.Tr("control_sys_result_success"))
}

//删除栏目
func (this *TermControllers) Remove() {
	//判断参数是否为空
	id, _ := this.GetInt64("id")
	if id < 1 {
		this.DbJsonMsg("page_sql_err_msg", false)
		this.ServeJSON()
	}
	//判断参数是否合法
	postService := posts.AutoTermService()
	view_terms, _ := postService.FindById(id)
	if len(view_terms.CreateName) < 1 {
		this.DbJsonMsg("page_sql_err_msg", false)
		this.ServeJSON()
	}
	res_int, err := postService.DeleteById(id)
	//数据出现错误
	if err != nil || res_int == 1 {
		beego.Error(err)
		this.DbJsonMsg("page_sql_err_msg", false)
		this.ServeJSON()
	}
	//栏目下面还有文章
	if res_int == 2 {
		this.DbJsonMsg("control_sys_term_remove_post_notnull", false)
		this.ServeJSON()
	}
	this.DbJsonMsg("control_sys_term_remove_suc", true)
	this.ServeJSON()
}
