package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/golangpkg/qor-cms/models"
)

var (
	apiToken = beego.AppConfig.String("apiToken")
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

//使用http 接口调用发布。
func (c *MainController) GetApiPublish() {
	defer c.ServeJSON()
	token := c.GetString("token", "")

	if token == "" || apiToken != token {
		c.Data["json"] = "error token"
		return
	}
	c.Data["json"] = "ok"
	logs.Info("############### api  publish5page ###############")
	//生成html代码。
	models.GenArticleAndCategoryList(5)
}
