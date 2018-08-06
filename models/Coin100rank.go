package models

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"path/filepath"
	"github.com/golangpkg/qor-cms/conf/db"
)

//加密货币 前 100 排名。
type Coin100rank struct {
	gorm.Model
	OrderId    int64  //排序号
	Date       string //日期
	Currency   string //币种名称
	Name       string //币种名称
	Describe   string //币种描述
	IconUrl    string //币种图标
	Price      string //价格
	Vol        string //交易额
	Change     string //价格
	MaxSupply  string //发行量
	Supply     string //流通量
	MarketCap  string //总市值
	UpdateTime string //更新时间
}

//初始化ui界面显示内容。
func InitCoin100rankUi(adminConf *admin.Admin) {
	//文章分类
	coin100rank := adminConf.AddResource(&Coin100rank{}, &admin.Config{Name: "加密货币100", Menu: []string{"数据管理"}})
	coin100rank.Meta(&admin.Meta{Name: "OrderId", Label: "排序号"})
	coin100rank.Meta(&admin.Meta{Name: "Date", Label: "日期"})
	coin100rank.Meta(&admin.Meta{Name: "Currency", Label: "币种名称"})
	coin100rank.Meta(&admin.Meta{Name: "Name", Label: "币种名称"})
	coin100rank.Meta(&admin.Meta{Name: "Describe", Label: "币种描述"})
	coin100rank.Meta(&admin.Meta{Name: "IconUrl", Label: "币种图标"})
	coin100rank.Meta(&admin.Meta{Name: "Price", Label: "价格"})
	coin100rank.Meta(&admin.Meta{Name: "Vol", Label: "交易额"})
	coin100rank.Meta(&admin.Meta{Name: "Change", Label: "价格"})
	coin100rank.Meta(&admin.Meta{Name: "MaxSupply", Label: "发行量"})
	coin100rank.Meta(&admin.Meta{Name: "Supply", Label: "流通量"})
	coin100rank.Meta(&admin.Meta{Name: "MarketCap", Label: "总市值"})
	coin100rank.Meta(&admin.Meta{Name: "UpdateTime", Label: "更新时间"})
	//展示列表数据。
	coin100rank.IndexAttrs("Currency", "Name", "Price", "Vol", "Change", "MaxSupply", "Supply", "MarketCap")
	//新建&编辑
	coin100rank.NewAttrs("")
	coin100rank.EditAttrs("")
}

func GenCoin100rankList() {
	// Get all records

	var coin100ranks []Coin100rank
	//#使用 Preload 进行预加载数据。其实就是join Company表。然后把属性填充。
	db.DB.Where("  `date` = DATE_FORMAT(now(), '%Y%m%d') ", ).Find(&coin100ranks)

	data := make(map[string]interface{})
	data["DataList"] = coin100ranks
	data["WebBasePath"] = webBasePath
	data["SiteName"] = siteName
	data["ImageBasePath"] = imageBasePath

	fileName := htmlPath + "coin100/list.html"
	tmpName := tmplPath + "coin100/list.html"

	fileDir := filepath.Dir(fileName)
	//调用通用生成函数。
	GenFileByTemplate(fileName, fileDir, tmpName, data)
}
