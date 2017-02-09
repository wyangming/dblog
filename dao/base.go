package dao

import (
	"database/sql"
	"dblog/util/strutil"
	"fmt"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
)

type IBaseDao interface {
}

//要访问的数据库
var sql_db *sql.DB

//初始化链接
func init() {
	//&loc=Local
	mysqlurl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=Local", beego.AppConfig.String("dbuser"),
		beego.AppConfig.String("dbpass"), beego.AppConfig.String("dburl"), beego.AppConfig.String("dbport"), beego.AppConfig.String("dbname"))
	db, err := sql.Open("mysql", mysqlurl)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Please check your net or database connection info.")
		beego.Error(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Please check your net or database connection info.")
		beego.Error(err)
	}
	sql_db = db
}

//得到一个连接信息用于
func NewDB() *sql.DB {
	return sql_db
}

//查询总条数方法
//from_end代表sql语句from后的内容
//count_column代表count(column_name)，如果为空则会用count(*)代替
func Count(from_end, count_column string, db *sql.DB, params []interface{}) int64 {
	if len(from_end) < 1 {
		return int64(0)
	}
	count_sql := ""
	if len(count_column) > 1 {
		count_sql = fmt.Sprintf("%scount(%s)", count_sql, count_column)
	} else {
		count_sql = strutil.StrAppend(count_sql, "count(1)")
	}
	find_sql := fmt.Sprintf("select %s from %s", count_sql, from_end)
	fmt.Println(find_sql)
	rows, err := db.Query(find_sql, params...)
	if err != nil {
		return int64(0)
	}
	defer rows.Close()
	var nums sql.NullInt64
	for rows.Next() {
		err = rows.Scan(&nums)
		if err != nil {
			return int64(0)
		}
	}
	return nums.Int64
}

//追加分页信息
//sql没有分而立信息的sql语句
func AppendPageInfo(sql string) string {
	return strutil.StrAppend(sql, " limit ?, ?")
}
