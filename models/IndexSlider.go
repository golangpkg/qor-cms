package models

import (
	"github.com/jinzhu/gorm"
	"github.com/golangpkg/qor-cms/conf/db"
	"path/filepath"
)

type IndexSlider struct {
	gorm.Model
	//Id         int64  `orm:"auto"`
	Image string //图片地址
	Url   string //链接地址
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
