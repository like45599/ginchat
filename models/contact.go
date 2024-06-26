// @Author Bing
// @Date 2024-04-25 21:52:00
// @Desc
package models

import (
	"gorm.io/gorm"
)

// 人员关系
type Contact struct {
	gorm.Model
	OwnerId  uint //谁的关系信息
	TargetId uint //对应的谁
	Type     int  //对应的类型 0 1 2
	Desc     string
}

func (table *Contact) TableName() string {
	return "Contact"
}
