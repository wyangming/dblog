package models

import (
	"time"
)

//权限信息表
type DbAuthority struct {
	Id int64
	//权限名称
	AuthName string `orm:"null"`
	//权限url
	AuthUrl string `orm:"null"`
	//父级权限编号
	ParentId int64
	//最后更新时间
	UpdateTime time.Time `orm:"auto_now"`
	//创建时间
	CreateTime time.Time `orm:"auto_now_add"`
	//创建人编号
	CreateUser int64
	//0权限类型分组即父权限，1正常url跳转链接，2返回数据的链接
	AuthType int `orm:"default(0)"`
	//权限说明
	Remark string
	//状态：1删除，0正常
	Active int `orm:"default(0)"`
}

//用户信息表
type DbUser struct {
	Id int64
	//密码
	Pwd string
	//用户名
	UserName string
	//最后更新时间
	UpdateTime time.Time `orm:"auto_now"`
	//创建时间
	CreateTime time.Time `orm:"auto_now_add"`
	//管理的站点id
	SiteId int64
	//状态：1删除，0正常
	Active int `orm:"default(0)"`
}

//角色表
type DbRole struct {
	Id int64
	//角色名
	RoleName string `orm:"null"`
	//最后更新时间
	UpdateTime time.Time `orm:"auto_now"`
	//创建时间
	CreateTime time.Time `orm:"auto_now_add"`
	//创建人编号
	CreateUser int64
	//状态：1删除，0正常
	Active int `orm:"default(0)"`
}

//用户角色表
type DbUserRole struct {
	Id string `orm:"pk"`
	//用户编号
	UserId int64
	//角色编号
	RoleId int64
	//最后更新时间
	UpdateTime time.Time `orm:"auto_now"`
	//创建时间
	CreateTime time.Time `orm:"auto_now_add"`
	//创建人编号
	CreateUser int64
	//0主要角色，1兼职角色
	MainRole int `orm:"default(0)"`
	//状态：1删除，0正常
	Active int `orm:"default(0)"`
}

//角色权限表
type DbRoleAuth struct {
	Id string `orm:"pk"`
	//角色编号
	RoleId int64
	//权限编号
	AuthorId int64
	//创建时间
	CreateTime time.Time `orm:"auto_now_add"`
	//创建人编号
	CreateUser int64
	//状态：1删除，0正常
	Active int `orm:"default(0)"` //1删除，0正常
}
