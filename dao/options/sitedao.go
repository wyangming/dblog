package options

import (
	"database/sql"
	"dblog/dao"
	"dblog/models"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"strings"
	"time"
)

type ISiteDao interface {
	dao.IBaseDao
	//根据编号查询一个站点的信息
	FindById(id int64) (site ViewSite, err error)
	//保存一个站点信息
	Save(ent models.DbSite) (res bool, id int64, err error)
}

func AutoSiteDao() ISiteDao {
	return ent_site_dao
}

var ent_site_dao = &siteDao{}

//站点视图
type ViewSite struct {
	models.DbSite
	//创建人
	CreateName string
}
type siteDao struct {
}

//保存一个站点信息
func (this *siteDao) Save(ent models.DbSite) (res bool, id int64, err error) {
	res = false
	var save_sql string
	//所有的参数信息
	params := make([]interface{}, 0)
	//判断是否是添加
	if ent.Id < 1 {
		save_sql = fmt.Sprintf("insert into db_site (%s) values(?,?,?,?,?,?,?,?,?)", siteTableColumn(""))
		//把id列去掉
		save_sql = strings.Replace(save_sql, "id, ", "", 1)
	} else {
		//create_user, update_time, create_time, active
		save_sql = "update db_site set site_name=?, site_sub_title=?, site_desc=?, site_tag=?, site_url=?, sys_page_size=?" +
			" , update_time=? where id=?"
	}
	//共用的参数
	params = append(params, ent.SiteName)
	params = append(params, ent.SiteSubTitle)
	params = append(params, ent.SiteDesc)
	params = append(params, ent.SiteTag)
	params = append(params, ent.SiteUrl)
	params = append(params, ent.SysPageSize)
	//当前时间写添加修改不同参数的处理
	now_time := time.Now()
	if ent.Id < 1 {
		params = append(params, ent.CreateUser)
		params = append(params, now_time)
		params = append(params, now_time)
		params = append(params, ent.Active)
	} else {
		params = append(params, now_time)
		params = append(params, ent.Id)
	}
	db := dao.NewDB()
	var db_res sql.Result
	db_res, err = db.Exec(save_sql, params...)
	if err != nil {
		return
	}
	if ent.Id < 1 {
		id, err = db_res.LastInsertId()
		if err != nil {
			return
		}
	}
	if count, _ := db_res.RowsAffected(); count > 0 {
		res = true
		return
	}
	return
}

//根据编号查询一个站点的信息
func (this *siteDao) FindById(id int64) (site ViewSite, err error) {
	if id < 1 {
		return
	}
	sel_sql := fmt.Sprintf("select %s from db_site where active=0 and id=?", siteTableColumn(""))
	//参数
	params := make([]interface{}, 0)
	params = append(params, id)
	db := dao.NewDB()
	rows, err := db.Query(sel_sql, params...)
	defer rows.Close()
	if err != nil {
		return
	}
	for rows.Next() {
		var (
			id, create_user, active, sys_page_size                   sql.NullInt64
			site_name, site_sub_title, site_desc, site_url, site_tag sql.NullString
			update_time, create_time                                 mysql.NullTime
		)
		err = rows.Scan(&id, &site_name, &site_sub_title, &site_desc, &site_tag, &site_url,
			&sys_page_size, &create_user, &update_time, &create_time, &active)
		if err != nil {
			return
		}
		site.DbSite = models.DbSite{
			Id:           int64(id.Int64),
			SiteName:     site_name.String,
			SiteSubTitle: site_sub_title.String,
			SiteDesc:     site_desc.String,
			SiteTag:      site_tag.String,
			SiteUrl:      site_url.String,
			SysPageSize:  int(sys_page_size.Int64),
			CreateUser:   int64(create_user.Int64),
			UpdateTime:   update_time.Time,
			CreateTime:   create_time.Time,
			Active:       int(active.Int64),
		}
	}
	return
}

//返回表里把有的列
func siteTableColumn(table_alias string) string {
	if len(table_alias) > 0 {
		table_alias = fmt.Sprintf("%s.", table_alias)
	}
	return fmt.Sprintf("%sid, %ssite_name, %ssite_sub_title, %ssite_desc, %ssite_tag, %ssite_url, %ssys_page_size,"+
		" %screate_user, %supdate_time, %screate_time, %sactive", table_alias, table_alias, table_alias,
		table_alias, table_alias, table_alias, table_alias, table_alias, table_alias, table_alias, table_alias)
}
