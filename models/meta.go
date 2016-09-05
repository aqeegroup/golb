package models

// Meta 文章的属性 - 分类、标签
type Meta struct {
	ID          int    `xorm:"pk autoincr"`
	Name        string `xorm:"varchar(100)"`
	Slug        string `xorm:"varchar(100)"`       // 分类缩略名用于创建友好的链接形式, 建议使用字母, 数字, 下划线和横杠
	Type        string `xorm:"varchar(100) index"` // 类型 category 和 tag
	Description string `xorm:"varchar(255)"`       // 描述
	Count       int    // 计数
	Order       int    // 排序
	ParentID    int    // 父级id

	Children []Meta `xorm:"-"` // 子分类
}
