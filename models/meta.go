package models

import (
	"strings"

	"github.com/go-xorm/xorm"
	"github.com/jiayx/go/stringutil"
)

// Meta 文章的属性 - 分类、标签
type Meta struct {
	ID          int64  `xorm:"pk autoincr"`
	Name        string `xorm:"varchar(100)"`
	Slug        string `xorm:"varchar(100)"`       // 分类缩略名用于创建友好的链接形式, 建议使用字母, 数字, 下划线和横杠
	Type        string `xorm:"varchar(100) index"` // 类型 category 和 tag
	Description string `xorm:"varchar(255)"`       // 描述
	Count       int    // 计数
	Order       int    // 排序
	ParentID    int64  // 父级id

	Children []Meta `xorm:"-"` // 子分类
}

// CreateOrUpdate 创建一个分类或者标签
func (m *Meta) CreateOrUpdate() error {
	var err error
	if m.ID > 0 {
		_, err = x.Where("id=?", m.ID).Update(m)
	} else {
		_, err = x.InsertOne(m)
	}

	return err
}

// CateNameExist 判断分类 name 是否已经存在
func (m *Meta) CateNameExist() (bool, error) {

	s := x.Where("type=?", "category").And("name=?", m.Name)
	if m.ID > 0 {
		s.And("id<>?", m.ID)
	}
	count, err := s.Count(&Meta{})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CateSlugExist 判断分类 slug 是否已经存在
func (m *Meta) CateSlugExist() (bool, error) {
	s := x.Where("type=?", "category").And("slug=?", m.Slug)
	if m.ID > 0 {
		s.Where("id<>?", m.ID)
	}
	count, err := s.Count(&Meta{})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// FindAllCates 查询全部分类
func FindAllCates() (*[]Meta, error) {
	cates := &[]Meta{}
	err := x.Where("type=?", "category").Find(cates)
	return cates, err
}

// FindAllTags 查询全部标签
func FindAllTags() (*[]Meta, error) {
	tags := &[]Meta{}
	err := x.Where("type=?", "tag").Find(tags)
	return tags, err
}

// DeleteMetas 根据 id 删除标签或分类
func DeleteMetas(ids string) (int64, error) {
	id := strings.Split(ids, ",")

	return x.In("id", id).Delete(&Meta{})
}

// DeleteMetasByPostID 根据 postId 删除标签和分类 - 估计没用
func DeleteMetasByPostID(s *xorm.Session, postID int64) (int64, error) {
	return s.Where("post_id=?", postID).Delete(&Relationship{})
}

// CatesCount 统计分类个数
func CatesCount() (int, error) {
	count, err := x.Where("type=?", "category").Count(&Meta{})
	return int(count), err
}

// TagsCount 统计标签个数
func TagsCount() (int, error) {
	count, err := x.Where("type=?", "tag").Count(&Meta{})
	return int(count), err

}

// TagNameExist 判断Tag name 是否已经存在
func TagNameExist(name string) (bool, error) {
	count, err := x.Where("type=?", "tag").And("name=?", name).Count(&Meta{})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CreateOrFindTag 查询 tag 或者插入 新增
func CreateOrFindTag(tagNames []string) ([]*Meta, error) {
	// 已经有的标签
	hasTags := []*Meta{}
	err := x.Select("id, name").In("name", tagNames).Find(&hasTags)
	if err != nil {
		return nil, err
	}

	var hasTagNames []string
	for _, tag := range hasTags {
		hasTagNames = append(hasTagNames, tag.Name)
	}

	for _, name := range tagNames {
		if !stringutil.InArray(name, hasTagNames) {
			tag := &Meta{
				Name: name,
				Slug: name,
				Type: "tag",
			}

			// 这里不采用一次性插入 是为了插入之后能拿到自增 ID
			if _, err := x.Insert(tag); err != nil {
				return nil, err
			}

			hasTags = append(hasTags, tag)
		}
	}

	return hasTags, err
}
