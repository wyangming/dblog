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
	Find(offset, rows int64, filter map[string]interface{}) (nums int64, terms []posts.ViewPosts, err error)
	//根据编号查询
	FindById(id int64) (posts.ViewPosts, error)
	//修改文章状态
	//type_str操作类型replase发布 del删除
	UpdateStatus(id int64, type_str string) (bool, error)
	//查询发布文章的总数
	PostsNum() int64
}

//得到接口实例
func AutoPostsService() IPostService {
	return ent_posts_service
}

//实现接口内容
var ent_posts_service = &postService{}

type postService struct {
}

//查询发布文章的总数
func (this *postService) PostsNum() int64 {
	postsDao := posts.AutoPostsDao()
	return postsDao.PostsNum()
}

//保存文章信息
func (this *postService) SavePost(info map[string]interface{}) (bool, error) {
	postsDao := posts.AutoPostsDao()
	return postsDao.SavePost(info)
}

//查询文章信息
func (this *postService) Find(offset, rows int64, filter map[string]interface{}) (nums int64, terms []posts.ViewPosts, err error) {
	postsDao := posts.AutoPostsDao()
	return postsDao.Find(offset, rows, filter)
}

//根据编号查询
func (this *postService) FindById(id int64) (posts.ViewPosts, error) {
	if id < 1 {
		return posts.ViewPosts{}, nil
	}
	postsDao := posts.AutoPostsDao()
	return postsDao.FindById(id)
}

//修改文章状态
func (this *postService) UpdateStatus(id int64, type_str string) (bool, error) {
	if id < 1 || len(type_str) < 1 {
		return false, nil
	}
	var (
		bol = false
		err error
	)
	postsDao := posts.AutoPostsDao()
	switch type_str {
	case "replase":
		bol, err = postsDao.UpdateSinglePro(id, "active", 0)
		break
	case "del":
		bol, err = postsDao.UpdateSinglePro(id, "active", 1)
		break
	}
	return bol, err
}
