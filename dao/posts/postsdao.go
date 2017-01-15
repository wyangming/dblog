package posts

import (
	"database/sql"
	"dblog/dao"
	"dblog/models"
	"dblog/util/strutil"
	"dblog/util/uuid"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"strings"
	"time"
)

type IPostsDao interface {
	dao.IBaseDao
	//查询在一个栏目下是否还有文章
	FindHasByTermId(termId int64) (bool, error)
	//保存文章信息
	SavePost(info map[string]interface{}) (bool, error)
	//查询文章信息
	Find(offset, rows int64, filter *map[string]interface{}) (nums int64, terms []ViewPosts, err error)
	//根据编号查询
	FindById(id int64) (ViewPosts, error)
}

func AutoPostsDao() IPostsDao {
	return ent_posts_Dao
}

var ent_posts_Dao = &postsDao{}

type postsDao struct {
}

//文章的视图
type ViewPosts struct {
	models.DbPosts
	//创建人
	CreateName string
	//栏目
	PTermName string
	//栏目编号
	TermId int64
	//副栏目编号名称信息
	TtermF map[int64]string
}

//根据编号查询
func (this postsDao) FindById(id int64) (term ViewPosts, err error) {
	if id < 1 {
		return
	}
	columns := fmt.Sprintf("%s, tp.term_id, t.term_name, u.user_name", postsTableColumn("p"))
	from_end := "db_posts p left join db_term_posts tp on p.id=tp.post_id left join db_terms t on tp.term_id=t.id left join db_user u on p.create_user=u.id where 1=1" +
		" and p.id=?"
	//参数
	params := make([]interface{}, 0)
	params = append(params, id)
	db := dao.NewDB()
	from_end = fmt.Sprintf("%s and p.active!=1 and tp.active!=1 and t.active!=1 and u.active!=1", from_end)
	rows, err := db.Query(fmt.Sprintf("select %s from %s", columns, from_end), params...)
	defer rows.Close()
	if err != nil {
		return
	}
	for rows.Next() {
		term, err = viewPost(rows, true, true)
		if err != nil {
			return
		}
	}
	//查询副栏目信息
	sel_sql := "select t.id,t.term_name from db_terms t left join db_term_posts tp on t.id=tp.term_id where tp.post_id=? and t.active!=1 and tp.main_posts=1"
	rows_f, err := db.Query(sel_sql, params...)
	defer rows_f.Close()
	if err != nil {
		return
	}
	term_f := make(map[int64]string, 0)
	for rows_f.Next() {
		var id sql.NullInt64
		var term_name sql.NullString
		err = rows_f.Scan(&id, &term_name)
		if err != nil {
			continue
		}
		term_f[id.Int64] = term_name.String
	}
	term.TtermF = term_f
	return
}

//查询文章信息
func (this postsDao) Find(offset, pagesize int64, filter *map[string]interface{}) (nums int64, terms []ViewPosts, err error) {
	if pagesize < 1 {
		return
	}
	columns := fmt.Sprintf("%s, tp.term_id, t.term_name, u.user_name", postsTableColumn("p"))
	//去掉文章的两个内容
	columns = strings.Replace(columns, ", p.summary, p.html_content, p.text_content", "", 1)
	from_end := "db_posts p left join db_term_posts tp on p.id=tp.post_id left join db_terms t on tp.term_id=t.id left join db_user u on p.create_user=u.id where 1=1"
	//参数
	params := make([]interface{}, 0)
	db := dao.NewDB()
	//查询总条数
	nums = dao.Count(from_end, "p.id", db, nil)
	//设置分页查询
	params = append(params, offset)
	params = append(params, pagesize)
	from_end = fmt.Sprintf("%s and p.active!=1 and tp.active!=1 and t.active!=1 and u.active!=1", from_end)
	rows, err := db.Query(fmt.Sprintf("select %s from %s", columns, dao.AppendPageInfo(from_end, offset, pagesize)), params...)
	defer rows.Close()
	if err != nil {
		return
	}
	for rows.Next() {
		term, err := viewPost(rows, false, false)
		if err == nil {
			terms = append(terms, term)
		}
	}
	//查询副栏目的编号
	return
}

//返回单个文章内容
//hasCon是否有查询文章的内容
func viewPost(rows *sql.Rows, hasCon, hasSum bool) (ViewPosts, error) {
	var (
		title, shot_title, tags, source_url, author, summary, html_content, text_content, term_name, user_name sql.NullString
		id, create_user, site_id, active, term_id                                                              sql.NullInt64
		release_time, update_time, create_time                                                                 mysql.NullTime
		terms                                                                                                  ViewPosts
		err                                                                                                    error
	)
	//有内容
	if hasCon {
		//有内容，有摘要
		if hasSum {
			err = rows.Scan(&id, &title, &shot_title, &tags, &source_url, &author, &summary, &html_content, &text_content,
				&release_time, &update_time, &create_time, &create_user, &site_id, &active, &term_id, &term_name, &user_name)
		} else {
			err = rows.Scan(&id, &title, &shot_title, &tags, &source_url, &author, &html_content, &text_content,
				&release_time, &update_time, &create_time, &create_user, &site_id, &active, &term_id, &term_name, &user_name)
		}
	} else {
		//无内容，有摘要
		if hasSum {
			err = rows.Scan(&id, &title, &shot_title, &tags, &source_url, &author, &summary,
				&release_time, &update_time, &create_time, &create_user, &site_id, &active, &term_id, &term_name, &user_name)
		} else {
			err = rows.Scan(&id, &title, &shot_title, &tags, &source_url, &author,
				&release_time, &update_time, &create_time, &create_user, &site_id, &active, &term_id, &term_name, &user_name)
		}
	}
	if err != nil {
		return terms, err
	}
	terms = ViewPosts{
		DbPosts: models.DbPosts{
			Id:          int64(id.Int64),
			Title:       title.String,
			ShotTitle:   shot_title.String,
			Tags:        tags.String,
			SourceUrl:   source_url.String,
			Author:      author.String,
			Summary:     summary.String,
			HtmlContent: html_content.String,
			TextContent: text_content.String,
			ReleaseTime: release_time.Time,
			UpdateTime:  update_time.Time,
			CreateTime:  create_time.Time,
			CreateUser:  int64(create_user.Int64),
			SiteId:      int64(site_id.Int64),
			Active:      int(active.Int64),
		},
		CreateName: user_name.String,
		PTermName:  term_name.String,
		TermId:     int64(term_id.Int64),
	}
	return terms, nil
}

//查询在一个栏目下是否还有文章
func (this postsDao) FindHasByTermId(termId int64) (bool, error) {
	if termId < 1 {
		return false, nil
	}
	sql_count := "select count(p.id) from db_posts p left join db_term_posts tp on p.id=tp.post_id where p.active=0 and tp.active=0 and tp.term_id=?"
	//参数
	params := make([]interface{}, 0)
	params = append(params, termId)
	db := dao.NewDB()
	var count sql.NullInt64
	rows, err := db.Query(sql_count, params...)
	defer rows.Close()
	if err != nil {
		return false, err
	}
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return false, nil
		}
		return count.Int64 > 0, nil
	}
	return false, nil
}

//保存文章信息
func (this postsDao) SavePost(info map[string]interface{}) (bool, error) {
	id := info["id"].(int64)
	post, _ := info["post"].(*models.DbPosts)
	term_id := info["term_id"].(int64)
	now_time := time.Now()
	post.UpdateTime = now_time
	var err error
	//连接数据信息
	db := dao.NewDB()
	tx, _ := db.Begin()
	//添加时的判断
	if id < 1 {
		post.CreateTime = now_time
		id, err = insert(tx, post)
		if err != nil || id < 1 {
			tx.Rollback()
			return false, err
		}
	} else {
		//修改时的判断
		res, err := update(tx, post)
		if err != nil || !res {
			tx.Rollback()
			return false, err
		}
	}
	//栏目与文章的关系
	res, err := saveTermPosts(tx, id, term_id)
	if !res || err != nil {
		tx.Rollback()
		return false, err
	}
	tx.Commit()
	return true, nil
}

//保存文章与栏目的关系
func saveTermPosts(db *sql.Tx, postid, mainTermId int64, fTermIds ...int64) (bool, error) {
	//修改原有的关系数据状态为删除
	update_sql := "update db_term_posts set active=1 where post_id=?"
	db.Exec(update_sql, postid)
	value_tmp := "(?,?,?,?,?)"
	insert_sql := fmt.Sprintf("insert into db_term_posts(%s) values%s", dbTermPostsTableColum(""), value_tmp)
	//参数
	add_params := make([]interface{}, 0)
	postTermParams(&add_params, postid, mainTermId, true)
	//判断是否有副栏目
	if len(fTermIds) > 0 {
		for term_id := range fTermIds {
			insert_sql = fmt.Sprintf("%s, %s", insert_sql, value_tmp)
			postTermParams(&add_params, postid, int64(term_id), false)
		}
	}
	res, err := db.Exec(insert_sql, add_params...)
	if err != nil {
		return false, err
	}
	if count, _ := res.RowsAffected(); count > 0 {
		return true, nil
	}
	return false, nil
}

//添加栏目文章关系参数
func postTermParams(params *[]interface{}, postid, termsid int64, mainTerm bool) {
	id := uuid.Rand().Hex()
	*params = append(*params, id)
	*params = append(*params, termsid)
	*params = append(*params, postid)
	if mainTerm {
		*params = append(*params, 0)
	} else {
		*params = append(*params, 1)
	}
	*params = append(*params, 0)
}

//修改一个文章信息
func update(db *sql.Tx, post *models.DbPosts) (bool, error) {
	update_sql := "update db_posts set title=?, shot_title=?, tags=?, source_url=?, author=?, summary=?, html_content=?, text_content==?," +
		" release_time=?, update_time=?, create_time=?, create_user=?, site_id=? where id=? and active!=1"
	res, err := db.Exec(update_sql, post.Title, post.ShotTitle, post.Tags, post.SourceUrl,
		post.Author, post.Summary, post.HtmlContent, post.TextContent, post.ReleaseTime, post.UpdateTime, post.CreateTime,
		post.CreateUser, post.SiteId, post.Id)
	if err != nil {
		return false, err
	}
	if count, _ := res.RowsAffected(); count > 0 {
		return false, nil
	}
	return true, nil
}

//保存一个文章信息
func insert(db *sql.Tx, post *models.DbPosts) (int64, error) {
	insert_column := strings.Replace(postsTableColumn(""), "id, ", "", 1)
	insert_sql := fmt.Sprintf("insert into db_posts (%s) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)", insert_column)
	res, err := db.Exec(insert_sql, post.Title, post.ShotTitle, post.Tags, post.SourceUrl,
		post.Author, post.Summary, post.HtmlContent, post.TextContent, post.ReleaseTime, post.UpdateTime, post.CreateTime,
		post.CreateUser, post.SiteId, post.Active)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	if count, _ := res.RowsAffected(); count > 0 {
		return id, nil
	}
	return 0, nil
}

//返回表里把有的列
func postsTableColumn(table_alias string) string {
	if len(table_alias) > 0 {
		table_alias = strutil.StrAppend(table_alias, ".")
	}
	return fmt.Sprintf("%sid, %stitle, %sshot_title, %stags, %ssource_url, %sauthor, %ssummary, %shtml_content, %stext_content,"+
		" %srelease_time, %supdate_time, %screate_time, %screate_user, %ssite_id, %sactive", table_alias, table_alias, table_alias,
		table_alias, table_alias, table_alias, table_alias, table_alias, table_alias, table_alias, table_alias, table_alias,
		table_alias, table_alias, table_alias)
}

//返回栏目关系表的列
func dbTermPostsTableColum(table_alias string) string {
	if len(table_alias) > 0 {
		table_alias = strutil.StrAppend(table_alias, ".")
	}
	return fmt.Sprintf("%sid, %sterm_id, %spost_id, %smain_posts, %sactive", table_alias, table_alias, table_alias,
		table_alias, table_alias)
}
