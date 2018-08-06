package models

import (
	"github.com/jinzhu/gorm"
	"github.com/golangpkg/qor-cms/conf/db"
	"github.com/astaxie/beego/logs"
	"path/filepath"
	"fmt"
	"math"
	"github.com/qor/admin"
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

func InitJobUi(adminConf *admin.Admin) {
	//################################ 招聘模块 ################################
	jobCompany := adminConf.AddResource(&JobCompany{}, &admin.Config{Name: "招聘公司管理", Menu: []string{"招聘管理"}})
	jobCompany.Meta(&admin.Meta{Name: "Name", Label: "公司名称"})
	jobCompany.Meta(&admin.Meta{Name: "IndustryType", Label: "行业分类", Config: &admin.SelectOneConfig{
		Collection: []string{"计算机/互联网/通信/电子", "会计/金融/银行/保险", "贸易/消费/制造/营运", "制药/医疗",
			"广告/媒体", "房地产/建筑", "专业服务/教育/培训", "服务业", "物流/运输", "能源/原材料", "政府/非营利组织/其他"}}})
	jobCompany.Meta(&admin.Meta{Name: "CompanyInfo", Label: "公司描述", Type: "kindeditor"})

	job := adminConf.AddResource(&Job{}, &admin.Config{Name: "招聘职位管理", Menu: []string{"招聘管理"}})
	job.Meta(&admin.Meta{Name: "IsPublish", Label: "是否发布", Type: "checkbox"})
	job.Meta(&admin.Meta{Name: "Title", Label: "标题", Type: "text"})
	job.Meta(&admin.Meta{Name: "Salary", Label: "薪水（/月）", Type: "text"})
	job.Meta(&admin.Meta{Name: "JobCompany", Label: "公司名称"})
	job.Meta(&admin.Meta{Name: "Locale", Label: "工作地点", Config: &admin.SelectOneConfig{
		Collection: []string{"北京", "上海", "广州", "深圳", "天津", "杭州", "成都"}}})
	job.Meta(&admin.Meta{Name: "Education", Label: "学历", Config: &admin.SelectOneConfig{
		Collection: []string{"大专", "本科", "硕士", "博士"}}})
	job.Meta(&admin.Meta{Name: "Age", Label: "年龄（岁）", Type: "text"})
	job.Meta(&admin.Meta{Name: "WorkExperience", Label: "工作年限", Config: &admin.SelectOneConfig{
		Collection: []string{"1-3年", "3-5年", "5-10年", "10年以上"}}})
	job.Meta(&admin.Meta{Name: "Department", Label: "所属部门", Type: "text"})
	job.Meta(&admin.Meta{Name: "ReportTo", Label: "汇报对象", Type: "text"})
	job.Meta(&admin.Meta{Name: "PublishDate", Label: "发部日期", Type: "date"})
	job.Meta(&admin.Meta{Name: "JobInfo", Label: "职位描述", Type: "kindeditor"})
	job.IndexAttrs("Title", "Salary", "JobCompany", "Locale", "Education", "Age", "WorkExperience")
	//新增
	job.NewAttrs("IsPublish", "Title", "Salary", "JobCompany", "Locale", "Education", "Age", "WorkExperience",
		"Department", "ReportTo", "PublishDate", "JobInfo")
	//编辑
	job.EditAttrs("IsPublish", "Title", "Salary", "JobCompany", "Locale", "Education", "Age", "WorkExperience",
		"Department", "ReportTo", "PublishDate", "JobInfo")
	//增加发布功能：
	// 发布按钮，显示到右侧上面。
	job.Action(&admin.Action{
		Name:  "publishJobList",
		Label: "发布",
		Handler: func(actionArgument *admin.ActionArgument) error {
			logs.Info("############### publishJobList ###############")
			//生成html代码。
			GenJobList()
			return nil
		},
		Modes: []string{"collection"},
	})
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
