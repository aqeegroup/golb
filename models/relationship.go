package models

// Relationship Post 和 Meta 的关系表
type Relationship struct {
	PostID int `xorm:"pk"`
	MetaID int `xorm:"pk"`
}
