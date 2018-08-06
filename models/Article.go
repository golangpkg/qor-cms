package models

import (
	"github.com/jinzhu/gorm"
	"github.com/golangpkg/qor-cms/conf/db"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
	"path/filepath"
	"fmt"
	"math"
	"github.com/qor/admin"
	"regexp"
	"strings"
	"github.com/qor/qor/resource"
	"github.com/qor/qor"
	"time"
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

type Category struct {
	gorm.Model
	//Id         int64  `orm:"auto"`
	Name string //用户名
	//Description    string  //描述
}

var (
	tmplPath      = beego.AppConfig.String("publishArticleTmplPath")
	htmlPath      = beego.AppConfig.String("publishArticleHtmlPath")
	webBasePath   = beego.AppConfig.String("webBasePath")
	siteName      = beego.AppConfig.String("siteName")
	imageBasePath = beego.AppConfig.String("uploadBaseUrl")
	page          = 10
)

func InitArticleUi(adminConf *admin.Admin) {
	// Create resources from GORM-backend model
	//文章分类
	category := adminConf.AddResource(&Category{}, &admin.Config{Name: "分类管理", Menu: []string{"资源管理"}})
	category.Meta(&admin.Meta{Name: "Name", Label: "名称"})
	//PageCount: 5,
	article := adminConf.AddResource(&Article{}, &admin.Config{Name: "文章管理", Menu: []string{"资源管理"}})
	article.Meta(&admin.Meta{Name: "Title", Label: "标题", Type: "text"})
	article.Meta(&admin.Meta{Name: "ImgUrl", Label: "图片", Type: "kindimage"})
	article.Meta(&admin.Meta{Name: "Content", Label: "内容", Type: "kindeditor"})
	article.Meta(&admin.Meta{Name: "Category", Label: "分类"})
	article.Meta(&admin.Meta{Name: "CreatedAt", Label: "创建时间"})
	article.Meta(&admin.Meta{Name: "UpdatedAt", Label: "更新时间"})
	article.Meta(&admin.Meta{Name: "Url", Label: "地址", Type: "readonly"})
	article.Meta(&admin.Meta{Name: "IsPublish", Label: "是否发布", Type: "checkbox"})
	article.IndexAttrs("Title", "Category", "IsPublish", "Url", "ImgUrl", "CreatedAt", "UpdatedAt")
	//新增
	article.NewAttrs("Title", "Url", "IsPublish", "Category", "ImgUrl", "Content")
	//编辑
	article.EditAttrs("Title", "Url", "IsPublish", "Category", "ImgUrl", "Content")
	//增加发布功能：
	// 发布按钮，显示到右侧上面。
	article.Action(&admin.Action{
		Name:  "publishAll",
		Label: "全部发布",
		Handler: func(actionArgument *admin.ActionArgument) error {
			logs.Info("############### publishAll ###############")
			//生成html代码。
			GenArticleAndCategoryList(0)
			return nil
		},
		Modes: []string{"collection"},
	})
	// 发布按钮，显示到右侧上面。
	article.Action(&admin.Action{
		Name:  "publish5page",
		Label: "增量发布5页",
		Handler: func(actionArgument *admin.ActionArgument) error {
			logs.Info("############### publish5page ###############")
			//生成html代码。
			GenArticleAndCategoryList(5)
			return nil
		},
		Modes: []string{"collection"},
	})

	article.AddProcessor(&resource.Processor{
		Name: "process_store_data", // register another processor with
		Handler: func(value interface{}, metaValues *resource.MetaValues, context *qor.Context) error {
			if article, ok := value.(*Article); ok {
				// do something...
				logs.Info("################ article ##################")
				if article.Url == "" {
					t := article.CreatedAt //time.Now()
					if t.IsZero() { //如果创建事件为空。
						t = time.Now()
					}
					url := fmt.Sprintf("%d-%02d/%d.html", t.Year(), t.Month(), t.Unix())
					logs.Info(t, url)
					article.Url = url
				}
				//更新摘要。新建，修改都更新。
				if article.Content != "" {
					//如果摘要为空，且内容不为空。
					//去除所有尖括号内的HTML代码，并换成换行符
					re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
					newContent := re.ReplaceAllString(article.Content, "\n")
					//去除连续的换行符
					re, _ = regexp.Compile("\\s{2,}")
					newContent = re.ReplaceAllString(newContent, "\n")
					newContent = strings.TrimSpace(newContent)
					newContentRune := []rune(newContent)
					if len(newContentRune) > 75 {
						article.Description = string(newContentRune[0:75])
					} else {
						article.Description = newContent
					}
					logs.Info("description: ", article.Description)
				}

			}
			return nil
		},
	})
}

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

	var allArticles []Article
	for i := 1; i <= int(pageAll); i ++ {
		tmpList := GenArticlePage(i, count, page, id, maxPage)
		for _, article := range tmpList {
			allArticles = append(allArticles, article)
		}
	}

	logs.Info(" ################# allArticles len :", len(allArticles))

	//生成全部数据。
	for i, article := range allArticles {
		logs.Info(" ################# number article:", tmplPath, htmlPath, i)
		//logs.Info("number article :", i )
		data := make(map[string]interface{})
		data["Article"] = article
		data["SiteName"] = siteName
		data["CategoryId"] = article.CategoryId
		//#设置每一页的前一页，下一页数据。
		if i != 0 {
			tmpArticle := allArticles[i-1]
			data["PrevName"] = tmpArticle.Title
			data["PrevUrl"] = tmpArticle.Url
		}
		if i != (len(allArticles) - 1) {
			tmpArticle := allArticles[i+1]
			data["NextName"] = tmpArticle.Title
			data["NextUrl"] = tmpArticle.Url
		}

		fileName := htmlPath + article.Url
		GenArticleDetial(data, fileName)
	}
}

func GenArticlePage(pageNum, count, page int, id uint, maxPage int) (articles []Article) {
	// Get all records
	limitPage := (pageNum - 1) * page //开始的数据num，page 大小 这个是从 1 开始的。主要是为了分页标签方便。

	if id == 0 { //普通文章
		db.DB.Where(" is_publish = ? ", "1").Order("id desc").
			Limit(page).Offset(limitPage).Find(&articles)
	} else { //分类
		db.DB.Where(" is_publish = ? and category_id = ? ", "1", id).
			Order("id desc").Limit(page).Offset(limitPage).Find(&articles)
	}

	logs.Info(" ################# page limit offset :", pageNum, limitPage, page)
	//logs.Info(" ################# :", tmplPath, htmlPath, pageNum)
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
		return articles
	} else {
		return []Article{}
	}
	//进行debug，将数据打印到页面当中。
	//t.Execute(os.Stdout, data)
}

//生成文章详细。
func GenArticleDetial(data map[string]interface{}, fileName string) {

	fileDir := filepath.Dir(fileName)
	//调用通用生成函数。
	GenFileByTemplate(fileName, fileDir, tmplPath+"article/detail.html", data)
}
