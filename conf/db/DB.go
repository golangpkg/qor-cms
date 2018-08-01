package db

import (
	"fmt"
	"github.com/astaxie/beego"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
)

//声明全局变量
var (
	url_format = "%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local"
	//数据库注册。
	dbhost     = beego.AppConfig.String("dbhost")
	dbport     = beego.AppConfig.String("dbport")
	dbuser     = beego.AppConfig.String("dbuser")
	dbpassword = beego.AppConfig.String("dbpassword")
	db         = beego.AppConfig.String("db")
	// Set up the database
	mysql_url = fmt.Sprintf(url_format, dbuser, dbpassword, dbhost, dbport, db)
	// 注册数据库，可以是sqlite3 也可以是 mysql 换下驱动就可以了。
	DB, _ = gorm.Open("mysql", mysql_url)
	//DB, _ = gorm.Open("sqlite3", "demo.db")
)

func init() {
	fmt.Println("	mysql_url: ", mysql_url)
}
