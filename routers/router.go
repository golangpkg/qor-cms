package routers

import (
	"github.com/golangpkg/qor-cms/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//beego.Router("/", &controllers.MainController{})
	userInfoController := &controllers.UserInfoController{}
	beego.Router("/", userInfoController, "get:LoginIndex")
	beego.Router("/auth/login", userInfoController, "get:LoginIndex")
	beego.Router("/auth/login", userInfoController, "post:Login")
	beego.Router("/auth/logout", userInfoController, "get:Logout")
	beego.Router("/admin/common/kindeditor/upload", &controllers.FileUploadController{}, "post:Upload")

}
