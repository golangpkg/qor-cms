package models

import (
	"github.com/jinzhu/gorm"
	"github.com/golangpkg/qor-cms/conf/db"
	"github.com/astaxie/beego/logs"
	"path/filepath"
	"fmt"
	"math"
)

type Job struct {
	gorm.Model
	Title          string                     //标题
	Salary         string                     //薪水
	Locale         string                     //工作地点
	Education      string                     //学历
	Age            string                     //年龄
	WorkExperience string                     //工作年限
	Department     string                     //所属部门
	ReportTo       string                     //汇报对象
	PublishDate    string                     //发部日期
	IsPublish      bool                       //是否发布。
	JobInfo        string `gorm:"type:text"`  //职位描述
	JobCompany     JobCompany                 //职位所属公司
	JobCompanyName string `form:"JobCompany"` //职位所属公司名称
	JobCompanyId   int64  `form:"JobCompany"` //职位所属公司Id
}

type JobCompany struct {
	gorm.Model
	Name         string                    //标题
	IndustryType string                    //行业分类
	CompanyInfo  string `gorm:"type:text"` //企业描述
}

func GenJobList() {
	var count int

	db.DB.Model(&Job{}).Where(" is_publish = ? ", "1").Count(&count)

	pageAll := math.Ceil(float64(count) / float64(page))
	logs.Info("pageAll : ", count, pageAll)

	for i := 1; i <= int(pageAll); i ++ {
		GenJobPage(i, count, page, )
	}
}

func GenJobPage(pageNum, count, page int) {
	// Get all records
	limitPage := (pageNum - 1) * page //开始的数据num，page 大小 这个是从 1 开始的。主要是为了分页标签方便。

	var jobs []Job
	//#使用 Preload 进行预加载数据。其实就是join Company表。然后把属性填充。
	db.DB.Preload("JobCompany").Where(" is_publish = ? ", "1").Limit(page).Offset(limitPage).Find(&jobs)

	logs.Info(" ################# page limit offset :", pageNum, limitPage, page)
	data := make(map[string]interface{})
	data["DataList"] = jobs
	data["WebBasePath"] = webBasePath
	data["SiteName"] = siteName
	data["ImageBasePath"] = imageBasePath

	//将分页参数传入到页面中。
	pageInfo := Page{PageSize: page, TotalCount: count, CurrentPage: pageNum}
	strUrl := "/index%d.html"
	data["PageHtml"] = pageInfo.ToHtml(strUrl)

	indexPageName := "index.html"
	if pageNum > 1 { //第一页就是index.html
		indexPageName = fmt.Sprintf("index%d.html", pageNum)
	}

	fileName := htmlPath + "job/" + indexPageName
	tmpName := tmplPath + "job/list.html"

	fileDir := filepath.Dir(fileName)
	//调用通用生成函数。
	GenFileByTemplate(fileName, fileDir, tmpName, data)
	//生成全部数据。
	for i, job := range jobs {
		logs.Info(" ################# number job:", job, i)
		//logs.Info("number article :", i )
		data := make(map[string]interface{})
		data["Data"] = job
		data["SiteName"] = siteName
		fileName := fmt.Sprintf("%sjob/%d.html", htmlPath, job.ID)
		GenJobeDetial(data, fileName)
	}
}

//生成文章详细。
func GenJobeDetial(data map[string]interface{}, fileName string) {
	fileDir := filepath.Dir(fileName)
	//调用通用生成函数。
	GenFileByTemplate(fileName, fileDir, tmplPath+"job/detail.html", data)
}
