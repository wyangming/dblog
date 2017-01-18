package posts

import (
	"dblog/dao/posts"
	"dblog/models"
	"dblog/service"
)

//服务接口
type ITermService interface {
	service.IBaseService
	SaveTerm(term *models.DbTerms) error
	//根据url验证在一个站点里是否有重复栏目
	//id是过滤掉修改时自己url
	RepeadByUrl(url string, siteid, id int64) (bool, error)
	//查询栏目信息
	Find(offset, rows int64, filter map[string]interface{}) (nums int64, terms []posts.ViewTerms, err error)
	//根据id查询一个栏目
	FindById(id int64) (posts.ViewTerms, error)
	//根据id删除一个栏目信息
	//return int 0删除success 1删除fail 2栏目下有文章
	DeleteById(id int64) (int, error)
	//查询所有的栏目信息
	FindAll() ([]models.DbTerms, error)
}

//得到接口实例
func AutoTermService() ITermService {
	return ent_term_service
}

//实现接口内容
var ent_term_service = &termService{}

type termService struct {
}

func (this *termService) DeleteById(id int64) (int, error) {
	if id < 1 {
		return 1, nil
	}
	postsDao := posts.AutoPostsDao()
	has_posts, err := postsDao.FindHasByTermId(id)
	if err != nil {
		return 1, err
	}
	if has_posts {
		return 2, nil
	}
	terDao := posts.AutoTermDao()
	del_bool, err := terDao.DeleteById(id)
	if err != nil {
		return 1, err
	}
	if del_bool {
		return 0, nil
	}
	return 1, nil
}
func (this *termService) SaveTerm(term *models.DbTerms) error {
	terDao := posts.AutoTermDao()
	return terDao.SaveTerm(term)

}
func (this *termService) RepeadByUrl(url string, siteid, id int64) (bool, error) {
	terDao := posts.AutoTermDao()
	return terDao.RepeadByUrl(url, siteid, id)
}

//查询栏目信息
func (this *termService) Find(offset, rows int64, filter map[string]interface{}) (nums int64, terms []posts.ViewTerms, err error) {
	terDao := posts.AutoTermDao()
	nums, terms, err = terDao.Find(offset, rows, filter)
	return
}

//根据id查询一个栏目
func (this *termService) FindById(id int64) (posts.ViewTerms, error) {
	terDao := posts.AutoTermDao()
	return terDao.FindById(id)
}

//查询所有的栏目信息
func (this *termService) FindAll() ([]models.DbTerms, error) {
	terDao := posts.AutoTermDao()
	return terDao.FindAll()
}
