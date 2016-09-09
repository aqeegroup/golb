package models

import "fmt"

// Post 文章
type Post struct {
	ID           int64 `xorm:"pk autoincr"`
	Title        string
	Slug         string `xorm:"varchar(100) index"`
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
func (p *Post) Create() (int64, error) {

	s := x.NewSession()
	defer s.Close()
	s.Begin()

	id, err := x.InsertOne(p)
	if err != nil {
		s.Rollback()
		return 0, err
	}
	if len(p.Slug) == 0 {
		p.Slug = string(id)
		id, err = x.Update(p)
		if err != nil {
			s.Rollback()
			return 0, err
		}
	}

	err = s.Commit()
	if err != nil {
		return 0, err
	}
	return id, nil
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

// RecentPosts 最近文章
func RecentPosts(limit int) (*[]Post, error) {

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
