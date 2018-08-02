package models

import (
	"github.com/jinzhu/gorm"
	"github.com/golangpkg/qor-cms/conf/db"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
	"path/filepath"
	"fmt"
	"math"
)

type Article struct {
	gorm.Model
	//Id         int64  `orm:"auto"`
	Title       string                    //标题
	ImgUrl      string                    //文章图片
	Content     string `gorm:"type:text"` //内容
	Description string `gorm:"type:text"` //内容描述
	Category    Category
	CategoryId  int64  `form:"category"` //分类
	Url         string                   //地址
	IsPublish   bool                     //是否发布。
	//publish2.Schedule
}

var (
	tmplPath      = beego.AppConfig.String("publishArticleTmplPath")
	htmlPath      = beego.AppConfig.String("publishArticleHtmlPath")
	webBasePath   = beego.AppConfig.String("webBasePath")
	siteName      = beego.AppConfig.String("siteName")
	imageBasePath = beego.AppConfig.String("uploadBaseUrl")
	page          = 10
)

func GenArticleAndCategoryList(maxPage int) {
	//生成文章
	var categorys []Category
	db.DB.Find(&categorys)
	logs.Info("GenArticleAndCategoryList#####")
	logs.Info(len(categorys))
	for _, category := range categorys {
		//GenArticleList(i)
		logs.Info(category)
		//公用生成模板。按照分类生成。
		GenArticleList(category.ID, maxPage)
	}
	//公用生成模板。全部文章生成。
	GenArticleList(0, maxPage)
}

func GenArticleList(id uint, maxPage int) {
	var count int
	if id == 0 { //普通文章
		db.DB.Model(&Article{}).Where(" is_publish = ? ", "1").Count(&count)
	} else { //分类
		db.DB.Model(&Article{}).Where(" is_publish = ? and category_id = ? ", "1", id).Count(&count)
	}

	pageAll := math.Ceil(float64(count) / float64(page))
	logs.Info("pageAll : ", count, pageAll)
	for i := 1; i <= int(pageAll); i ++ {
		GenArticlePage(i, count, page, id, maxPage)
	}
}

func GenArticlePage(pageNum, count, page int, id uint, maxPage int) {
	// Get all records
	var articles []Article
	limitPage := (pageNum - 1) * page //开始的数据num，page 大小 这个是从 1 开始的。主要是为了分页标签方便。

	if id == 0 { //普通文章
		db.DB.Where(" is_publish = ? ", "1").Order("id desc").
			Limit(page).Offset(limitPage).Find(&articles)
	} else { //分类
		db.DB.Where(" is_publish = ? and category_id = ? ", "1", id).
			Order("id desc").Limit(page).Offset(limitPage).Find(&articles)
	}

	logs.Info(" ################# page limit offset :", pageNum, limitPage, page)
	logs.Info(" ################# :", tmplPath, htmlPath, pageNum)
	data := make(map[string]interface{})
	data["ArticleList"] = articles
	data["WebBasePath"] = webBasePath
	data["SiteName"] = siteName
	data["ImageBasePath"] = imageBasePath
	data["CategoryId"] = id //增加分类Id。

	//将分页参数传入到页面中。
	pageInfo := Page{PageSize: page, TotalCount: count, CurrentPage: pageNum}
	strUrl := "/index%d.html"
	if id > 0 { //分类
		strUrl = fmt.Sprintf("/cat-%d-", id) + "index%d.html"
	}
	data["PageHtml"] = pageInfo.ToHtml(strUrl)

	indexPageName := "index.html"
	if pageNum > 1 { //第一页就是index.html
		indexPageName = fmt.Sprintf("index%d.html", pageNum)
	}

	fileName := htmlPath + indexPageName
	tmpName := tmplPath + "article/list.html"

	if id > 0 { //分类
		fileName = htmlPath + fmt.Sprintf("cat/%d/", id) + indexPageName
		tmpName = tmplPath + "article/categoryList.html"
	}
	fileDir := filepath.Dir(fileName)
	//调用通用生成函数。
	GenFileByTemplate(fileName, fileDir, tmpName, data)

	//支持增量更新，当maxPage = 0 的时候是全量，否则要要小于 maxPage.
	if id == 0 && (maxPage == 0 || pageNum <= maxPage) { //普通文章
		for _, article := range articles {
			GenArticleDetial(article)
		}
	}
	//进行debug，将数据打印到页面当中。
	//t.Execute(os.Stdout, data)

}

func GenArticleDetial(article Article) {
	data := make(map[string]interface{})
	data["Article"] = article
	data["SiteName"] = siteName

	fileName := htmlPath + article.Url
	fileDir := filepath.Dir(fileName)

	//调用通用生成函数。
	GenFileByTemplate(fileName, fileDir, tmplPath+"article/detail.html", data)

}
