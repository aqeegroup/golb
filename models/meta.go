package models

import "strings"

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

// Create 创建一个分类或者标签
func (m *Meta) Create() error {
	var err error
	if m.ID > 0 {
		_, err = x.Where("id=?", m.ID).Update(m)
	} else {
		_, err = x.InsertOne(m)
	}

	return err
}

// CateNameExist 判断分类 name 是否已经存在
func CateNameExist(name string) (bool, error) {
	count, err := x.Where("type=?", "post").And("name=?", name).Count(&Meta{})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CateSlugExist 判断分类 slug 是否已经存在
func CateSlugExist(slug string) (bool, error) {
	count, err := x.Where("type=?", "post").And("slug=?", slug).Count(&Meta{})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// TagNameExist 判断Tag name 是否已经存在
func TagNameExist(name string) (bool, error) {
	count, err := x.Where("type=?", "tag").And("name=?", name).Count(&Meta{})
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

// DeleteMetas 根据 id 删除标签 或者 分类
func DeleteMetas(ids string) (int64, error) {
	id := strings.Split(ids, ",")
	return x.In("id", id).Delete(&Meta{})
}
