package models

import (
	//"dblog/util/uuid"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	//_ "github.com/mattn/go-sqlite3"
	//"strconv"
	//"strings"
	//"os"
	//"path"
	"time"
)

//初始化数据库配置
const (
	// 设置数据库名称
	_SQLITE3_DRIVER = "sqlite3"
	_MYSQL_DRIVER   = "mysql"
)

func init() {
	//initDataBase()
}
func initDataBase() {
	dbtype := beego.AppConfig.String("dbtype")

	//orm.DefaultTimeLoc = time.Local
	orm.DefaultTimeLoc = time.UTC

	orm.RegisterModel(new(DbSite), new(DbUser), new(DbRole), new(DbUserRole), new(DbAuthority),
		new(DbRoleAuth), new(DbTerms), new(DbPosts), new(DbTermPosts))

	// 开启 ORM 调试模式
	orm.Debug = true

	//mysql数据库处理
	if dbtype == "mysql" {
		orm.RegisterDriver(_MYSQL_DRIVER, orm.DRMySQL) //&loc=Local
		mysqlurl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=Local", beego.AppConfig.String("dbuser"),
			beego.AppConfig.String("dbpass"), beego.AppConfig.String("dburl"), beego.AppConfig.String("dbport"), beego.AppConfig.String("dbname"))
		orm.RegisterDataBase("default", _MYSQL_DRIVER, mysqlurl, 30, 60)
	} else { //sqlite3数据库处理
		//sqlite3还有个时间问题
		/*
			dbname := beego.AppConfig.String("dbname")
			if !IsExist(dbname) {
				os.MkdirAll(path.Dir(dbname), os.ModePerm)
				os.Create(dbname)
			}
			orm.RegisterDriver(_SQLITE3_DRIVER, orm.DRSqlite)
			orm.RegisterDataBase("default", _SQLITE3_DRIVER, dbname, 30)
		*/
	}
	// 自动建表
	orm.RunSyncdb("default", false, true)
	//初始化数据
	initDate()
}

//判断文件是否存在
/*
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
*/

//初始化数据
func initDate() {
	con := orm.NewOrm()
	qs := con.QueryTable("db_site")
	site := &DbSite{}
	err := qs.One(site)
	if site.Id > 0 {
		return
	}
	//如果有数据就不需要初始化
	//站点信息
	site.SiteName = "dblog官网"
	site.SiteTag = "dblog官网,dblog"
	site.SiteUrl = "localhost"
	site.Active = 0
	_, err = con.Insert(site)
	if err != nil {
		beego.Error(err)
		return
	}
	//默认用户信息
	user := &DbUser{
		Pwd:      "admin",
		UserName: "admin",
		SiteId:   site.Id,
		Active:   0,
	}
	_, err = con.Insert(user)
	if err != nil {
		beego.Error(err)
		return
	}
}
