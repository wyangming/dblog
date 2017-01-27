package options

import (
	"dblog/dao/options"
	"dblog/models"
	"dblog/service"
)

//服务接口
type ISiteService interface {
	service.IBaseService
	//根据id查询站点
	FindById(id int64) (options.ViewSite, error)
	//保存一个站点信息
	Save(ent models.DbSite) (bool, int64, error)
}

//得到接口实例
func AutoSiteService() ISiteService {
	return ent_site_service
}

//实现接口内容
var ent_site_service = &siteService{}

type siteService struct {
}

//根据编号查询
func (this *siteService) FindById(id int64) (options.ViewSite, error) {
	if id < 1 {
		return options.ViewSite{}, nil
	}
	siteDao := options.AutoSiteDao()
	return siteDao.FindById(id)
}

//保存一个站点信息
func (this *siteService) Save(ent models.DbSite) (bool, int64, error) {
	siteDao := options.AutoSiteDao()
	return siteDao.Save(ent)
}

//得到一个站点的分页大小
func SitePageSize(id int64) int64 {
	if id < 1 {
		//默认分页大小为10
		return 10
	}
	siteDao := options.AutoSiteDao()
	site, _ := siteDao.FindById(id)
	return int64(site.SysPageSize)
}
