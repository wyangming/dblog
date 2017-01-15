package models

import (
	"time"
)

//文章信息
type DbPosts struct {
	Id int64
	//标题
	Title string `orm:"null"`
	//短标题
	ShotTitle string `orm:"null"`
	//文章标签
	Tags string `orm:"null"`
	//文章来源
	SourceUrl string `orm:"null"`
	//作者
	Author string `orm:"null"`
	//摘要
	Summary string `orm:"null"`
	//html内容
	HtmlContent string `orm:"type(text);null"`
	//html过滤后的内容
	TextContent string `orm:"type(text);null"`
	//发布时间
	ReleaseTime time.Time `orm:"null"`
	//最后更新时间
	UpdateTime time.Time `orm:"auto_now"`
	//创建时间
	CreateTime time.Time `orm:"auto_now_add"`
	//创建人编号
	CreateUser int64
	//站点编号，归属于那个站点
	SiteId int64
	//状态，0正常，1删除，2草稿没有发布的
	Active int `orm:"default(0)"`
}

//栏目列表信息
type DbTermPosts struct {
	Id string `orm:"pk"`
	//栏目编号
	TermId int64
	//文章编号
	PostId int64
	//0主栏目关系，1副栏目关系
	MainPosts int `orm:"default(0)"`
	//状态：0正常，1删除
	Active int `orm:"default(0)"`
}

//栏目信息
type DbTerms struct {
	Id int64
	//栏目标签
	TermName string `orm:"null"`
	//栏目标签
	TermTag string `orm:"null"`
	//相对于父栏目来说的排序，默认为0
	Sort int `orm:"default(0)"`
	//栏目描述
	Description string `orm:"null"`
	//父级栏目编号
	ParentId int64
	//用于生成栏目的url字符串
	Slug string
	//最后更新时间
	UpdateTime time.Time `orm:"auto_now"`
	//创建时间
	CreateTime time.Time `orm:"auto_now_add"`
	//创建人编号
	CreateUser int64
	//站点编号，归属于那个站点
	SiteId int64 //归属于那个站点
	//状态：0正常，1删除
	Active int `orm:"default(0)"`
}
