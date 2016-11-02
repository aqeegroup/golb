package models

import "strings"

// Relationship Post 和 Meta 的关系表
type Relationship struct {
	PostID int64 `xorm:"pk"`
	MetaID int64 `xorm:"pk"`
}

// PostMeta 文章 id 与标签的对应
type PostMeta struct {
	Relationship `xorm:"extends"`
	Meta         `xorm:"extends"`
}

// TableName 查询结构体指定表名
func (PostMeta) TableName() string {
	return "relationship"
}

// FindMetasByPostID 查询指定文章的属性 - 包括分类和标签
func FindMetasByPostID(id int64) (*[]PostMeta, error) {
	postMetas := &[]PostMeta{}
	err := x.Where("post_id=?", id).
		Join("LEFT", "meta", "relationship.meta_id = meta.id").
		Find(postMetas)

	return postMetas, err
}

// DeleteRelationship 删除关系
func DeleteRelationship(ids string, typ string) (int64, error) {
	id := strings.Split(ids, ",")
	cond := typ + "_id"

	return x.In(cond, id).Delete(&Relationship{})
}
