// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameComment = "comment"

// Comment mapped from table <comment>
type Comment struct {
	Rpid    int       `gorm:"column:rpid;not null" json:"rpid"`
	From    int       `gorm:"column:from;not null" json:"from"`
	To      int       `gorm:"column:to;not null" json:"to"`
	Avid    int       `gorm:"column:avid;not null" json:"avid"`
	Ctime   time.Time `gorm:"column:ctime;not null" json:"ctime"`
	Like    int       `gorm:"column:like;not null" json:"like"` // 获赞数
	ID      int       `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Content string    `gorm:"column:content;not null" json:"content"`
}

// TableName Comment's table name
func (*Comment) TableName() string {
	return TableNameComment
}
