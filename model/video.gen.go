// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameVideo = "video"

// Video mapped from table <video>
type Video struct {
	Avid        int       `gorm:"column:avid;not null;primaryKey" json:"avid"`
	LastUpdated time.Time `gorm:"column:last_updated;not null;default:CURRENT_TIMESTAMP" json:"last_updated"`
}

// TableName Video's table name
func (*Video) TableName() string {
	return TableNameVideo
}
