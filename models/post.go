package models

import (
	"blog/modules/utility"
	"fmt"
	"strconv"
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

	metas []Meta `xorm:"-"`
}

// Create 创建文章
func (p *Post) Create() error {

	// 有自定义缩略名的话 要处理特殊字符
	if len(p.Slug) > 0 {
		p.Slug = utility.SlugNameFormat(p.Slug)
		// 处理完成要查询是否有重复
		var err error
		p.Slug, err = slugNameCheck(p.Slug)
		if err != nil {
			return err
		}
	}

	// 开始事务
	s := x.NewSession()
	defer s.Close()
	s.Begin()

	_, err := s.InsertOne(p)
	if err != nil {
		s.Rollback()
		return err
	}

	// 如果缩略名为空 则默认为id
	if len(p.Slug) == 0 {
		p.Slug = strconv.Itoa(int(p.ID))

		_, err = s.ID(p.ID).Update(p)
		if err != nil {
			s.Rollback()
			return err
		}
	}

	err = s.Commit()
	if err != nil {
		return err
	}
	return nil
}

// 查询是否有相同缩略名的文章
func countSlug(s string, id ...int) (int, error) {
	fmt.Println(s)
	session := x.Where("slug=?", s)

	if len(id) >= 1 {
		session.And("id<>?", id[0])
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
	post.Slug = string(post.ID)
	fmt.Println(post.Slug)

	return post, err
}

// FindPosts 查询所有文章带分页
func FindPosts(page, limit int) (*[]Post, error) {
	if page > 0 {
		page = page - 1
	}
	offset := page * limit
	posts := &[]Post{}
	err := x.Where("status=?", "publish").
		And("type=?", "post").
		Limit(limit, offset).
		Desc("publish_time").
		Find(posts)

	return posts, err
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
	if _, err := x.Exec("UPDATE `post` SET view = view + 1 WHERE id = ?", id); err != nil {
		return err
	}

	return nil
}

// CountPost 统计全部文章数目
func CountPost() (int64, error) {
	return x.Count(&Post{})
}

func slugNameCheck(s string) (string, error) {
	count := 1

	temp := s
	for i, err := countSlug(s); i > 0; i, err = countSlug(s) {
		fmt.Println(i)
		fmt.Println(count)
		if err != nil {
			return s, err
		}
		s = fmt.Sprintf("%s-%d", temp, count)
		count++
	}

	return s, nil
}
