package options

import (
	"dblog/controllers/sys"
	"dblog/models"
	service "dblog/service/options"
	"dblog/util/web"
	"fmt"
	"github.com/astaxie/beego"
)

type SiteControllers struct {
	sys.SysBaseController
}

//跳转后台设置页面
func (this *SiteControllers) Get() {
	sys_siteid, _ := this.Ctx.Input.Session(web.SYS_SITEID).(int64)
	this.TplName = "sys/options/site.html"
	site_service := service.AutoSiteService()
	site, err := site_service.FindById(sys_siteid)
	if err != nil {
		beego.Error(err)
	}
	this.Data["site"] = site
}

//保存
func (this *SiteControllers) Save() {
	this.TplName = "sys/options/site.html"
	site_name := this.GetString("site_name")
	if len(site_name) < 1 {
		this.Data["msg"] = this.Tr("control_sys_options_site_msg_name_nil")
		return
	}
	site_url := this.GetString("site_url")
	if len(site_url) < 1 {
		this.Data["msg"] = this.Tr("control_sys_options_site_msg_domain_nil")
		return
	}
	sys_page_size, page_err := this.GetInt32("sys_page_size")
	if page_err != nil {
		sys_page_size = 10
	}
	id, _ := this.GetInt64("id")
	site_service := service.AutoSiteService()
	var site models.DbSite
	if id < 1 {
		site = models.DbSite{
			Active: 0,
		}
		site.CreateUser = web.GetCurUserId(this.Controller)
	} else {
		site_view, err := site_service.FindById(id)
		if err != nil {
			beego.Error(err)
		}
		site = site_view.DbSite
	}
	site.SiteName = site_name
	site.SiteUrl = site_url
	site.SysPageSize = int(sys_page_size)
	site.SiteSubTitle = this.GetString("site_sub_title")
	site.SiteTag = this.GetString("site_tag")
	site.SiteDesc = this.GetString("site_desc")

	res, res_id, err := site_service.Save(site)
	if err != nil {
		beego.Error(err)
	}
	//操作类型
	action_type := this.Tr("page_btn_update")
	if id < 1 {
		action_type = this.Tr("page_btn_add")
		site.Id = res_id
	}
	//操作结果
	action_res := this.Tr("control_sys_result_faile")
	if res {
		action_res = this.Tr("control_sys_result_success")
	}
	this.Data["site"] = site
	this.Data["msg"] = fmt.Sprintf("%s%s%s%s", this.Tr("control_sys_options_site"),
		site_name, action_type, action_res)
}
