package main

import (
	_ "github.com/golangpkg/qor-cms/routers"
	"github.com/qor/admin"
	"net/http"
	"github.com/astaxie/beego"
	"github.com/golangpkg/qor-cms/conf/auth"
	"github.com/golangpkg/qor-cms/conf/db"
	"github.com/golangpkg/qor-cms/models"
	"github.com/qor/auth/auth_identity"
)

func main() {
	//开启session。配置文件 配置下sessionon = true即可。
	beego.BConfig.WebConfig.Session.SessionOn = true

	//DB, _ := gorm.Open("sqlite3", "demo.db") //for sqlite3
	DB := db.DB
	DB.AutoMigrate(&models.Category{}, &models.Article{}, &auth_identity.AuthIdentity{},
		&models.IndexSlider{}, &models.JobCompany{}, &models.Job{}, &models.Coin100rank{})

	// Initalize
	Admin := admin.New(&admin.AdminConfig{SiteName: "qor-cms", DB: DB, Auth: auth.AdminAuth{}})

	//初始化文章管理
	models.InitArticleUi(Admin)

	//初始化首页轮播图
	models.InitIndexSliderUi(Admin)

	//初始化招聘
	models.InitJobUi(Admin)

	//初始化Coin100rank
	models.InitCoin100rankUi(Admin)

	// 启动服务
	mux := http.NewServeMux()
	Admin.MountTo("/admin", mux)
	beego.Handler("/admin/*", mux)
	beego.Run()
}
