package user

import (
	"dblog/dao/user"
	"dblog/models"
	"dblog/service"
)

//服务接口
type IUserService interface {
	service.IBaseService
	Login(username, pwd string) (*models.DbUser, bool, error)
}

//得到接口实例
func AutoUserService() IUserService {
	return ent_user_service
}

//实现接口内容
var ent_user_service = &userService{}

type userService struct {
}

func (this *userService) Login(username, pwd string) (*models.DbUser, bool, error) {
	userDao := user.AutoUserDao()
	user, err := userDao.FindByName(username)
	if user.Id < 1 {
		return user, false, err
	}
	if username != user.UserName {
		return user, false, err
	}
	if pwd != user.Pwd {
		return user, false, err
	}
	return user, true, nil
}
