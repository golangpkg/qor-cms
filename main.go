package main

import (
	_ "github.com/golangpkg/qor-cms/routers"
	//   "github.com/qor/qor"
	"github.com/qor/admin"

	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/golangpkg/qor-cms/conf/auth"
	"github.com/golangpkg/qor-cms/conf/db"
	"github.com/golangpkg/qor-cms/models"
	"github.com/qor/auth/auth_identity"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
)

func main() {
	//开启session。配置文件 配置下sessionon = true即可。
	beego.BConfig.WebConfig.Session.SessionOn = true

	//DB, _ := gorm.Open("sqlite3", "demo.db") //for sqlite3
	DB := db.DB
	DB.AutoMigrate(&models.Category{}, &models.Article{}, &auth_identity.AuthIdentity{}, &models.IndexSlider{})

	// Initalize
	Admin := admin.New(&admin.AdminConfig{SiteName: "qor-cms", DB: DB, Auth: auth.AdminAuth{}})

	// Create resources from GORM-backend model
	//文章分类
	category := Admin.AddResource(&models.Category{}, &admin.Config{Name: "分类管理", Menu: []string{"资源管理"}})
	category.Meta(&admin.Meta{Name: "Name", Label: "名称"})
	//PageCount: 5,
	article := Admin.AddResource(&models.Article{}, &admin.Config{Name: "文章管理", Menu: []string{"资源管理"}})
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
		Name:  "publish",
		Label: "发布",
		Handler: func(actionArgument *admin.ActionArgument) error {
			logs.Info("############### publish ###############")
			//生成html代码。
			models.GenArticleList()
			return nil
		},
		Modes: []string{"collection"},
	})

	article.AddProcessor(&resource.Processor{
		Name: "process_store_data", // register another processor with
		Handler: func(value interface{}, metaValues *resource.MetaValues, context *qor.Context) error {
			if article, ok := value.(*models.Article); ok {
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

	//增加首页轮播图。
	indexSlider := Admin.AddResource(&models.IndexSlider{}, &admin.Config{Name: "首页轮播图", Menu: []string{"资源管理"}})
	indexSlider.Meta(&admin.Meta{Name: "Image", Label: "图片（860X300）地址", Type: "kindimage"})
	indexSlider.Meta(&admin.Meta{Name: "Url", Label: "链接地址"})
	indexSlider.Action(&admin.Action{
		Name:  "publish",
		Label: "发布",
		Handler: func(actionArgument *admin.ActionArgument) error {
			logs.Info("############### publish ###############")
			//生成html代码。
			models.GenIndexSlider()
			return nil
		},
		Modes: []string{"collection"},
	})

	// 启动服务
	mux := http.NewServeMux()
	Admin.MountTo("/admin", mux)
	beego.Handler("/admin/*", mux)
	beego.Run()
}
