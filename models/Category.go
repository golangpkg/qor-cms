package models

import "github.com/jinzhu/gorm"

type Category struct {
	gorm.Model
	//Id         int64  `orm:"auto"`
    Name    string //用户名
    //Description    string  //描述
}