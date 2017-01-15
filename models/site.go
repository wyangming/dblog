package models

import (
	"time"
)

//站点信息
type DbSite struct {
	Id int64
	//站点名称
	SiteName string `orm:"null"`
	//站点标签
	SiteTag string `orm:"null"`
	//站点绑定的url
	SiteUrl string `orm:"null"`
	//创建人编号
	CreateUser int64
	//最后更新时间
	UpdateTime time.Time `orm:"auto_now"`
	//创建时间
	CreateTime time.Time `orm:"auto_now_add"`
	//状态：0正常，1删除
	Active int `orm:"default(0)"` //1删除，0正常
}
