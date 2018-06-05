package routers

import (
	"github.com/golangpkg/qor-cms/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
