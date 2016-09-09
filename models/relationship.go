package models

// Relationship Post 和 Meta 的关系表
type Relationship struct {
	PostID int64 `xorm:"pk"`
	MetaID int64 `xorm:"pk"`
}
