package models

import (
	"fmt"
)

// Post 文章
type Post struct {
	ID           int `xorm:"pk autoincr"`
	Title        string
	Slug         string `xorm:"varchar(100) index unique"`
	Content      string `xorm:"text"`
	AuthorID     int
	Type         string `xorm:"varchar(20)"` // post post_draft
	Status       string `xorm:"varchar(31)"` // 公开publish 隐藏hidden 私密private
	Password     string `xorm:"varchar(32)"`
	AllowComment string
	View         int `xorm:"notnull default(0)"`
	Like         int `xorm:"notnull default(0)"`
	CreateTime   int
	UpdateTime   int `xorm:"updated"`

	metas []Meta `xorm:"-"`
}

// CreatePost 创建文章
func CreatePost(p *Post) {
	x.Insert(p)
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
		Desc("id").
		Find(posts)

	return posts, err
}

// RecentPosts 最近文章 - 这个并没有什么卵用
func RecentPosts(limit int) (*[]Post, error) {

	posts := &[]Post{}
	err := x.Where("status=?", "publish").
		And("type=?", "post").
		Limit(limit).
		Desc("id").
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
