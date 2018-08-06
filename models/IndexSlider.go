package models

import (
	"github.com/jinzhu/gorm"
	"github.com/golangpkg/qor-cms/conf/db"
	"path/filepath"
	"github.com/qor/admin"
	"github.com/astaxie/beego/logs"
)

type IndexSlider struct {
	gorm.Model
	//Id         int64  `orm:"auto"`
	Image string //图片地址
	Url   string //链接地址
}

func InitIndexSliderUi(adminConf *admin.Admin) {
	//增加首页轮播图。
	indexSlider := adminConf.AddResource(&IndexSlider{}, &admin.Config{Name: "首页轮播图", Menu: []string{"资源管理"}})
	indexSlider.Meta(&admin.Meta{Name: "Image", Label: "图片（860X300）地址", Type: "kindimage"})
	indexSlider.Meta(&admin.Meta{Name: "Url", Label: "链接地址"})
	indexSlider.Action(&admin.Action{
		Name:  "publish",
		Label: "发布",
		Handler: func(actionArgument *admin.ActionArgument) error {
			logs.Info("############### publish ###############")
			//生成html代码。
			GenIndexSlider()
			return nil
		},
		Modes: []string{"collection"},
	})
}

func GenIndexSlider() {
	var indexSliders []IndexSlider
	db.DB.Find(&indexSliders)

	data := make(map[string]interface{})
	data["IndexSliderList"] = indexSliders

	fileName := htmlPath + "widgets/indexSlider.html"
	fileDir := filepath.Dir(fileName)

	//调用通用生成函数。
	GenFileByTemplate(fileName, fileDir, tmplPath+"widgets/indexSlider.html", data)
}
