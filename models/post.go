package models

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Unknwon/paginater"
	"github.com/go-xorm/xorm"

	"blog/modules/utility"
)

// Post 文章
type Post struct {
	ID           int64 `xorm:"pk autoincr"`
	Title        string
	Slug         string `xorm:"varchar(100) index unique"`
	Content      string `xorm:"text"`
	AuthorID     int64
	Type         string `xorm:"varchar(20)"` // post post_draft
	Status       string `xorm:"varchar(31)"` // 公开publish 隐藏hidden 私密private
	Password     string `xorm:"varchar(32)"`
	AllowComment string
	View         int `xorm:"notnull default(0)"`
	Like         int `xorm:"notnull default(0)"`
	CreateTime   int64
	UpdateTime   int64 `xorm:"updated"`
	PublishTime  int64 `xorm:"index"`

	// Metas []PostMeta `xorm:"-"`
	Cates     []PostMeta `xorm:"-"`
	Tags      []PostMeta `xorm:"-"`
	CateNames []string   `xorm:"-"`
	TagNames  []string   `xorm:"-"`
}

// PostDetail 文章详情
type PostDetail struct {
	Post   `xorm:"extends"`
	Author User   `xorm:"extends"`
	metas  []Meta `xorm:"extends"`
}

// TableName PostDetail 的orm映射表名
func (PostDetail) TableName() string {
	return "post"
}

// Create 创建文章
func (p *Post) Create(cates, tags string) error {

	var err error
	// 对分类名的处理
	cateIds := utility.StringSplitInt64(cates, ",")
	// 对标签名的处理
	tagNames := strings.Split(tags, ",")

	// 有自定义缩略名的话 要处理特殊字符
	if len(p.Slug) > 0 {
		p.Slug = utility.SlugNameFormat(p.Slug)
		// 处理完成要查询是否有重复
		p.Slug, err = slugNameCheck(p.Slug, p.ID)
		if err != nil {
			return err
		}
	}

	fmt.Println("------------------------")
	fmt.Println(p)
	fmt.Println("------------------------")

	// 开始事务
	s := x.NewSession()
	defer s.Close()
	s.Begin()

	if p.ID > 0 {
		_, err = s.ID(p.ID).Update(p)
	} else {
		_, err = s.InsertOne(p)
	}
	if err != nil {
		s.Rollback()
		return err
	}

	// 如果缩略名为空 则默认为id
	if len(p.Slug) == 0 {
		p.Slug = strconv.Itoa(int(p.ID))

		// 处理完成要查询是否有重复
		p.Slug, err = slugNameCheck(p.Slug, p.ID)
		if err != nil {
			return err
		}

		// 确保没问题 再插入
		_, err = s.ID(p.ID).Cols("slug").Update(p)
		if err != nil {
			s.Rollback()
			return err
		}
	}

	// 如果是更新
	if p.ID > 0 {
		// 删除分类和标签
		_, err = p.deleteMetas(s)
		if err != nil {
			s.Rollback()
			return err
		}
	}
	// 插入分类
	relationships := []Relationship{}
	for _, cateID := range cateIds {
		relationships = append(relationships, Relationship{
			MetaID: cateID,
			PostID: p.ID,
		})
	}
	// 插入标签
	postTags, err := CreateOrFindTag(tagNames)
	if err != nil {
		return err
	}
	for _, tag := range postTags {
		relationships = append(relationships, Relationship{
			MetaID: tag.ID,
			PostID: p.ID,
		})
	}

	// 需要记录条数的话 只需要在这里记下全部 meta 的id 然后 update 就好了
	if len(relationships) > 0 {
		_, err = s.Insert(&relationships)
		if err != nil {
			s.Rollback()
			return err
		}
	}

	err = s.Commit()
	if err != nil {
		s.Rollback()
		return err
	}
	return nil
}

// deleteMetas 删除文章的标签和分类
func (p *Post) deleteMetas(s *xorm.Session) (int64, error) {
	return s.Where("post_id=?", p.ID).Delete(&Relationship{})
}

// DeletePosts 删除文章
func DeletePosts(ids string) error {
	id := strings.Split(ids, ",")
	_, err := x.In("id", id).Delete(&Post{})
	return err
}

// 查询是否有相同缩略名的文章
func countSlug(s string, id int64) (int, error) {
	// fmt.Println(s)
	session := x.Where("slug=?", s)

	if id > 0 {
		session.And("id<>?", id)
	}
	count, err := session.Count(&Post{})
	return int(count), err
}

// FindPostBySlug 根据缩略名查找文章
func FindPostBySlug(slug string) (*Post, error) {
	post := &Post{}
	has, err := x.Where("slug=?", slug).
		And("status<>?", "private").
		And("type=?", "post").
		Get(post)

	if !has {
		return nil, fmt.Errorf("没有找到这篇文章: %s", slug)
	}

	err = post.FindCateAndTagByPostID()

	log.Println(post)

	return post, err
}

// FindPostByID 根据 id 查找文章
func FindPostByID(id int64) (*Post, error) {
	post := &Post{}
	has, err := x.Where("id=?", id).And("type=?", "post").Get(post)
	if !has {
		return nil, fmt.Errorf("没有找到这篇文章: %d", id)
	}
	if err != nil {
		return post, err
	}
	err = post.FindCateAndTagByPostID()

	return post, err
}

// FindPosts 查询所有文章带分页
func FindPosts(page, limit int) (*[]PostDetail, error) {
	if page > 0 {
		page = page - 1
	}
	offset := page * limit
	posts := &[]PostDetail{}
	err := x.Where("status=?", "publish").
		Join("LEFT", "user", "post.author_id = user.id").
		And("type=?", "post").
		Limit(limit, offset).
		Desc("publish_time").
		Find(posts)

	return posts, err
}

// FindPostsDetail 查询所有文章详情带分页
func FindPostsDetail(page, limit int, isAdmin bool) ([]*PostDetail, error) {
	var err error
	if page > 0 {
		page = page - 1
	}
	offset := page * limit
	posts := make([]*PostDetail, 0, limit)

	session := x.NewSession()
	if !isAdmin {
		session.Where("status=?", "publish")
	}
	err = session.Join("LEFT", "user", "post.author_id = user.id").
		And("type=?", "post").
		Limit(limit, offset).
		Desc("publish_time").
		Find(&posts)

	if err != nil {
		return posts, err
	}

	for i := range posts {
		err := posts[i].FindCateAndTagByPostID()
		if err != nil {
			return posts, err
		}
	}

	return posts, err
}

// PostsCount 文章数统计
func PostsCount(isAdmin bool) (int, error) {
	var (
		count int64
		err   error
	)
	if isAdmin {
		count, err = x.Count(&Post{})
	} else {
		count, err = x.Where("status=?", "publish").Count(&Post{})
	}

	return int(count), err
}

// PostsPagination 文章分页
func PostsPagination(page, limit int, isAdmin bool) (*paginater.Paginater, error) {
	count, err := PostsCount(isAdmin)
	if err != nil {
		return nil, err
	}

	p := paginater.New(count, limit, page, 10)

	return p, nil
}

// LatestPosts 最近文章
func LatestPosts(limit int) (*[]Post, error) {
	posts := &[]Post{}
	err := x.Where("status=?", "publish").
		And("type=?", "post").
		Limit(limit).
		Desc("publish_time").
		Find(posts)

	return posts, err
}

// Like 文章点赞
func Like(id int) error {
	if _, err := x.Exec("UPDATE `post` SET like = like + 1 WHERE id = ?", id); err != nil {
		return err
	}

	return nil
}

// View 文章浏览次数自增1
func View(id int) error {
	_, err := x.Where("id=?", id).Incr("view", 1).Update(&Post{})
	return err
}

func slugNameCheck(s string, id int64) (string, error) {
	count := 1

	temp := s
	for i, err := countSlug(s, id); i > 0; i, err = countSlug(s, id) {
		// fmt.Println(i)
		// fmt.Println(count)
		if err != nil {
			return s, err
		}
		s = fmt.Sprintf("%s-%d", temp, count)
		count++
	}

	return s, nil
}

// FindCateAndTagByPostID 查询文章的 tag 和 category
func (p *Post) FindCateAndTagByPostID() error {
	// log.Println(p.ID)
	metas, err := FindMetasByPostID(p.ID)
	// log.Println(metas)
	if err != nil {
		return err
	}

	for _, meta := range *metas {
		if meta.Type == "category" {
			p.Cates = append(p.Cates, meta)
			p.CateNames = append(p.CateNames, meta.Name)
		} else if meta.Type == "tag" {
			p.Tags = append(p.Tags, meta)
			p.TagNames = append(p.TagNames, meta.Name)
		}
	}

	return nil
}
