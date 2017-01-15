package posts

import (
	"database/sql"
	"dblog/dao"
	"dblog/models"
	"dblog/util/strutil"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"time"
)

type ITermDao interface {
	dao.IBaseDao
	//保存一个栏目信息
	SaveTerm(term *models.DbTerms) error
	//根据url验证在一个站点里是否有重复栏目
	//id是过滤掉修改时自己url
	RepeadByUrl(url string, siteid, id int64) (bool, error)
	//查询栏目信息
	Find(offset, rows int64, filter *map[string]interface{}) (nums int64, terms []ViewTerms, err error)
	//根据id查询一个栏目
	FindById(id int64) (ViewTerms, error)
	//根据id删除一个栏目
	DeleteById(id int64) (bool, error)
	//查询所有的栏目信息
	FindAll() ([]models.DbTerms, error)
}

func AutoTermDao() ITermDao {
	return ent_term_Dao
}

var ent_term_Dao = &termDao{}

type termDao struct {
}

//栏目的视图
type ViewTerms struct {
	models.DbTerms
	//创建人
	CreateName string
	//父级栏目
	PTermName string
}

var (
	COLUMNS = fmt.Sprintf("%s, p.term_name as parent_name, u.user_name", termTableColumn("t"))
)

//查询所有的栏目信息
func (this termDao) FindAll() ([]models.DbTerms, error) {
	db := dao.NewDB()
	rows, err := db.Query(fmt.Sprintf("select %s from %s where active=0", termTableColumn(""), "db_terms"))
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	terms := make([]models.DbTerms, 0)
	for rows.Next() {
		var (
			term_name, term_tag, description, slug            sql.NullString
			id, sort, parent_id, create_user, site_id, active sql.NullInt64
			update_time, create_time                          mysql.NullTime
		)
		err = rows.Scan(&id, &term_name, &term_tag, &sort, &description, &parent_id, &slug, &update_time, &create_time, &create_user, &site_id, &active)
		if err != nil {
			return nil, err
		}
		term := models.DbTerms{
			Id:          int64(id.Int64),
			TermName:    term_name.String,
			TermTag:     term_tag.String,
			Sort:        int(sort.Int64),
			Description: description.String,
			ParentId:    parent_id.Int64,
			Slug:        slug.String,
			UpdateTime:  update_time.Time,
			CreateTime:  create_time.Time,
			CreateUser:  create_user.Int64,
			SiteId:      site_id.Int64,
			Active:      int(active.Int64),
		}
		terms = append(terms, term)
	}
	return terms, nil
}
func (this termDao) DeleteById(id int64) (bool, error) {
	if id < 1 {
		return false, nil
	}
	exec_sql := "update db_terms set active=1 where id=?"
	params := make([]interface{}, 0)
	params = append(params, id)
	db := dao.NewDB()
	res, err := db.Exec(exec_sql, params...)
	if err != nil {
		return false, err
	}
	if count, _ := res.RowsAffected(); count > 0 {
		return true, nil
	}
	return false, nil
}
func (this termDao) FindById(id int64) (ViewTerms, error) {
	//查询表与查询条件
	from_end := "db_terms t left join db_terms p on t.parent_id=p.id left join db_user u on t.create_user=u.id where t.active=0 and t.id=?"
	//参数
	params := make([]interface{}, 0)
	params = append(params, id)
	db := dao.NewDB()
	var term ViewTerms
	rows, err := db.Query(fmt.Sprintf("select %s from %s", COLUMNS, from_end), params...)
	defer rows.Close()
	if err != nil {
		return term, err
	}
	for rows.Next() {
		term, err = viewTerms(rows)
	}
	return term, err
}

//查询栏目信息
func (this termDao) Find(offset, pagesize int64, filter *map[string]interface{}) (nums int64, terms []ViewTerms, err error) {
	if pagesize < 1 {
		return 0, nil, nil
	}
	//参数
	params := make([]interface{}, 0)
	//查询表与查询条件
	from_end := "db_terms t left join db_terms p on t.parent_id=p.id left join db_user u on t.create_user=u.id where t.active=0"
	db := dao.NewDB()
	//查询总条数
	nums = dao.Count(from_end, "t.id", db, nil)
	//设置分页查询
	params = append(params, offset)
	params = append(params, pagesize)
	rows, err := db.Query(fmt.Sprintf("select %s from %s", COLUMNS, dao.AppendPageInfo(from_end, offset, pagesize)), params...)
	defer rows.Close()
	if err != nil {
		return
	}
	//results := make([]models.DbTerms, 0)
	for rows.Next() {
		term, err := viewTerms(rows)
		if err == nil {
			terms = append(terms, term)
		}
	}
	return
}
func viewTerms(rows *sql.Rows) (ViewTerms, error) {
	var (
		term_name, term_tag, description, slug, parent_name, user_name sql.NullString
		id, sort, active, parent_id, create_user, site_id              sql.NullInt64
		update_time, create_time                                       mysql.NullTime
	)
	var terms ViewTerms
	err := rows.Scan(&id, &term_name, &term_tag, &sort, &description, &parent_id, &slug, &update_time, &create_time, &create_user, &site_id, &active, &parent_name, &user_name)
	if err != nil {
		return terms, err
	}
	terms = ViewTerms{
		DbTerms: models.DbTerms{
			Id:          int64(id.Int64),
			TermName:    term_name.String,
			TermTag:     term_tag.String,
			Sort:        int(sort.Int64),
			Description: description.String,
			ParentId:    parent_id.Int64,
			Slug:        slug.String,
			UpdateTime:  update_time.Time,
			CreateTime:  create_time.Time,
			CreateUser:  create_user.Int64,
			SiteId:      site_id.Int64,
			Active:      int(active.Int64),
		},
		PTermName:  parent_name.String,
		CreateName: user_name.String,
	}
	return terms, nil
}

//保存一个栏目信息
func (this termDao) SaveTerm(term *models.DbTerms) error {
	if term.Id < 1 {
		_, err := saveTerm(term)
		return err
	}
	_, err := updateTerm(term)
	return err
}

//查询url是否重复
func (this termDao) RepeadByUrl(url string, siteid, id int64) (bool, error) {
	if len(url) < 1 || siteid < 1 {
		return true, nil
	}
	params := make([]interface{}, 0)
	params = append(params, url)
	params = append(params, siteid)
	str_count := "select count(id) from db_terms where active=0 and slug=? and site_id=?"
	if id > 0 {
		str_count = strutil.StrAppend(str_count, " and id!=?")
		params = append(params, id)
	}
	db := dao.NewDB()
	rows, err := db.Query(str_count, params...)
	defer rows.Close()
	if err != nil {
		return true, nil
	}
	var count sql.NullInt64
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return true, nil
		}
		//如果大于0表示重复
		return count.Int64 > 0, nil
	}
	return true, nil
}

//修改一个栏目信息
func updateTerm(ter *models.DbTerms) (bool, error) {
	ter.UpdateTime = time.Now()
	str_update := fmt.Sprintf("update db_terms set term_name=?, term_tag=?, sort=?, description=?, parent_id=?, slug=?, update_time=?," +
		" create_time=?, create_user=?, site_id=?, active=? where id=?")
	db := dao.NewDB()
	res, err := db.Exec(str_update, ter.TermName, ter.TermTag, ter.Sort, ter.Description, ter.ParentId, ter.Slug, ter.UpdateTime,
		ter.CreateTime, ter.CreateUser, ter.SiteId, ter.Active, ter.Id)
	if err != nil {
		return false, err
	}
	if count, _ := res.RowsAffected(); count > 0 {
		return true, nil
	}
	return false, nil
}

//添加一个栏目信息
func saveTerm(ter *models.DbTerms) (bool, error) {
	now_time := time.Now()
	ter.CreateTime = now_time
	ter.UpdateTime = now_time
	str_insert := fmt.Sprintf("insert into db_terms (term_name, term_tag, sort, description, parent_id, slug, update_time, create_time, create_user, site_id, active) values(?,?,?,?,?,?,?,?,?,?,?)")
	db := dao.NewDB()
	res, err := db.Exec(str_insert, ter.TermName, ter.TermTag, ter.Sort, ter.Description, ter.ParentId, ter.Slug, ter.UpdateTime,
		ter.CreateTime, ter.CreateUser, ter.SiteId, ter.Active)
	if err != nil {
		return false, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return false, err
	}
	ter.Id = id
	if count, _ := res.RowsAffected(); count > 0 {
		return true, nil
	}
	return false, nil
}

//返回表里把有的列
func termTableColumn(table_alias string) string {
	if len(table_alias) > 0 {
		table_alias = strutil.StrAppend(table_alias, ".")
	}
	return fmt.Sprintf("%sid, %sterm_name, %sterm_tag, %ssort, %sdescription, %sparent_id, %sslug, %supdate_time, %screate_time, %screate_user, %ssite_id, %sactive",
		table_alias, table_alias, table_alias, table_alias, table_alias, table_alias, table_alias, table_alias, table_alias, table_alias, table_alias, table_alias)
}
