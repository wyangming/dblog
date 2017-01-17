package user

import (
	"database/sql"
	"dblog/dao"
	"dblog/models"
	"dblog/util/strutil"
	"fmt"
	"github.com/go-sql-driver/mysql"
)

type IUserDao interface {
	dao.IBaseDao
	FindByName(username string) (*models.DbUser, error)
}

func AutoUserDao() IUserDao {
	return ent_user_Dao
}

var ent_user_Dao = &userDao{}

type userDao struct {
}

//返回表里把有的列
func tableColumn(table_alias string) string {
	if len(table_alias) > 0 {
		table_alias = strutil.StrAppend(table_alias, ".")
	}
	return fmt.Sprintf("%sid, %spwd, %suser_name, %supdate_time, %screate_time, %ssite_id, %sactive",
		table_alias, table_alias, table_alias, table_alias, table_alias, table_alias, table_alias)
}

//根据用户名查的用户
func (this *userDao) FindByName(username string) (*models.DbUser, error) {
	user := &models.DbUser{}
	if len(username) < 1 {
		return user, nil
	}
	db := dao.NewDB()
	rows, err := db.Query(fmt.Sprintf("select %s from db_user where active=0 and user_name=?", tableColumn("")), username)
	defer rows.Close()
	if err != nil {
		return user, err
	}
	for rows.Next() {
		var (
			pwd, user_name           sql.NullString
			update_time, create_time mysql.NullTime
			id, active, site_id      sql.NullInt64
		)
		err = rows.Scan(&id, &pwd, &user_name, &update_time, &create_time, &site_id, &active)
		if err != nil {
			fmt.Println(err)
			return user, err
		}
		if id.Valid {
			user.Id = id.Int64
		}
		if pwd.Valid {
			user.Pwd = pwd.String
		}
		if user_name.Valid {
			user.UserName = user_name.String
		}
		if update_time.Valid {
			user.UpdateTime = update_time.Time
		}
		if create_time.Valid {
			user.CreateTime = create_time.Time
		}
		if site_id.Valid {
			user.SiteId = site_id.Int64
		}
		if active.Valid {
			user.Active = int(active.Int64)
		}
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return user, err
	}
	return user, nil
}
