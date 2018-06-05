package models

import (
	"github.com/jinzhu/gorm"
)

type Article struct {
	gorm.Model
	//Id         int64  `orm:"auto"`
    Title    string //标题
    Content    string `gorm:"type:text"` //内容
	Category     Category
    CategoryId    int64 `form:"category"`//分类
    Url    string  //地址
	//publish2.Schedule
}