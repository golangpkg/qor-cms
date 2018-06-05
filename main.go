package main

import (
	"fmt"
	"net/http"
	//   "github.com/qor/qor"
	"github.com/qor/admin"
	"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/golangpkg/qor-cms/models"

	"github.com/qor/auth"
	"github.com/qor/auth/auth_identity"
	"github.com/qor/auth/providers/password"
	"github.com/qor/session/manager"
	"github.com/qor/redirect_back"
)

func main() {
	// Set up the database
	mysql_url := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local","root", "mariadb","127.0.0.1", "3306", "qor_cms")
	gormDB, _ := gorm.Open("mysql", mysql_url)
	//  DB, _ := gorm.Open("sqlite3", "demo.db") //for sqlite3
	gormDB.AutoMigrate(&models.Category{}, &models.Article{},&auth_identity.AuthIdentity{})

	var RedirectBack = redirect_back.New(&redirect_back.Config{
		SessionManager:  manager.SessionManager,
		IgnoredPrefixes: []string{"/auth"},
	})
	// Initialize Auth with configuration
	Auth := auth.New(&auth.Config{
		DB: gormDB,
		Redirector: auth.Redirector{RedirectBack},
	})
	// Register Auth providers
	// Allow use username/password
	Auth.RegisterProvider(password.New(&password.Config{}))

	// Initalize
	Admin := admin.New(&admin.AdminConfig{DB: gormDB})

	// Create resources from GORM-backend model
	Admin.AddResource(&models.Category{})
	Admin.AddResource(&models.Article{})

	mux := http.NewServeMux()
	// Mount admin to the mux
	Admin.MountTo("/admin", mux)
	// Mount Auth to Router
	mux.Handle("/auth/", Auth.NewServeMux())
	http.ListenAndServe(":8080", manager.SessionManager.Middleware(RedirectBack.Middleware(mux)))
	//使用beego启动
	//beego.Handler("/admin/*", mux)
	//beego.Run()
}