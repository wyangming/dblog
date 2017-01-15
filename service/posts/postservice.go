package posts

import (
	"dblog/dao/posts"
	"dblog/service"
)

//服务接口
type IPostService interface {
	service.IBaseService
	//保存文章信息
	SavePost(info map[string]interface{}) (bool, error)
	//查询文章信息
	Find(offset, rows int64, filter *map[string]interface{}) (nums int64, terms []posts.ViewPosts, err error)
	//根据编号查询
	FindById(id int64) (posts.ViewPosts, error)
}

//得到接口实例
func AutoPostsService() IPostService {
	return ent_posts_service
}

//实现接口内容
var ent_posts_service = &postService{}

type postService struct {
}

//保存文章信息
func (this postService) SavePost(info map[string]interface{}) (bool, error) {
	postsDao := posts.AutoPostsDao()
	return postsDao.SavePost(info)
}

//查询文章信息
func (this postService) Find(offset, rows int64, filter *map[string]interface{}) (nums int64, terms []posts.ViewPosts, err error) {
	postsDao := posts.AutoPostsDao()
	return postsDao.Find(offset, rows, filter)
}

//根据编号查询
func (this postService) FindById(id int64) (posts.ViewPosts, error) {
	if id < 1 {
		return posts.ViewPosts{}, nil
	}
	postsDao := posts.AutoPostsDao()
	return postsDao.FindById(id)
}
