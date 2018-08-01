package models

import (
	"fmt"
	"strconv"
	"strings"
)

const DEFAULT_PAGE_SIZE = 20

//定义 page对象。
type Page struct {
	PageSize    int
	TotalCount  int
	CurrentPage int
}

//活动开始索引
func (page *Page) GetStartIndex() int {
	return (page.GetCurrentPage() - 1) * page.PageSize
}

//获取结束索引
func (page *Page) GetEndIndex() int {
	return page.GetCurrentPage() * page.PageSize
}

//是否第一页
func (page *Page) IsFirstPage() bool {
	return page.GetCurrentPage() <= 1
}

//是否末页
func (page *Page) IsLastPage() bool {
	return page.GetCurrentPage() >= page.GetPageCount()
}

//获取下一页页码

func (page *Page) GetNextPage() int {
	if (page.IsLastPage()) {
		return page.GetCurrentPage()
	}
	return page.GetCurrentPage() + 1
}

//获取上一页页码
func (page *Page) GetPreviousPage() int {
	if (page.IsFirstPage()) {
		return 1
	}
	return page.GetCurrentPage() - 1
}

//获取当前页页码
func (page *Page) GetCurrentPage() int {
	if (page.CurrentPage == 0) {
		page.CurrentPage = 1
	}
	return page.CurrentPage
}

//取得总页数
func (page *Page) GetPageCount() int {
	if (page.TotalCount%page.PageSize == 0) {
		return page.TotalCount / page.PageSize
	} else {
		return page.TotalCount/page.PageSize + 1
	}
}

////取总记录数
//func (page *Page) GetTotalCount() int {
//	return page.TotalCount
//}

func (page *Page) HasNextPage() bool {
	return page.GetCurrentPage() < page.GetPageCount()
}

//该页是否有上一页.
func (page *Page) HasPreviousPage() bool {
	return page.GetCurrentPage() > 1
}

//故意要留的回车。方便页面展示。
const page_left = `<a class="mdl-button" style="min-width: 20px;" href="%s">&lt;</a>
`
const page_right = `<a class="mdl-button" style="min-width: 20px;" href="%s">&gt; </a>
`
const page_li = `<a class="mdl-button" style="min-width: 20px;" href="%s">%s</a>
`
const page_li_active = `<a class="mdl-button mdl-button--raised mdl-button--colored" style="min-width: 20px;" href="">%s</a>
`

//将int 转换成 string 格式。
//func ItoA(i int) string {
//	return strconv.FormatInt(i, 10)
//}

func urlFormat(url string, pageNo int) string {
	url_new := fmt.Sprintf(url, pageNo)
	if pageNo == 1 {
		url_new = strings.Replace(url_new, "index1.html", "index.html", -1)
	}
	return url_new
}

func (page *Page) ToHtml(url string) string {
	tmp_html := ""
	if !page.IsFirstPage() {
		tmp_html += fmt.Sprintf(page_left, urlFormat(url, page.GetPreviousPage()))
		tmp_html += fmt.Sprintf(page_li, urlFormat(url, 1), "1")
	}
	//当前页的尾巴长度
	var pagePos int = 5
	//显示前面的省略号
	if page.CurrentPage > (pagePos + 2) {
		tmp_html += fmt.Sprintf(page_li, "#", "...")
	}
	//增加前面的尾巴
	for i := pagePos; i >= 1; i-- {
		if (page.CurrentPage - i) > 1 {
			pageIndex := strconv.Itoa(page.CurrentPage - i)
			tmp_html += fmt.Sprintf(page_li, urlFormat(url, page.CurrentPage-i), pageIndex)
		}
	}
	//显示当前页号
	tmp_html += fmt.Sprintf(page_li_active, strconv.Itoa(page.CurrentPage))

	//增加后面的尾巴
	for i := 1; i <= pagePos; i++ {
		if (page.GetPageCount() - page.CurrentPage - i) > 0 {
			pageIndex := page.CurrentPage + i
			tmp_html += fmt.Sprintf(page_li, urlFormat(url, pageIndex), strconv.Itoa(pageIndex))
		}
	}
	//显示后面的省略号
	if (page.GetPageCount() - page.CurrentPage) > (pagePos + 1) {
		tmp_html += fmt.Sprintf(page_li, "", "...")
	}
	if page.HasNextPage() {
		pageIndex := strconv.Itoa(page.GetPageCount())
		tmp_html += fmt.Sprintf(page_li, urlFormat(url, page.GetPageCount()), pageIndex)
		tmp_html += fmt.Sprintf(page_right, urlFormat(url, page.GetNextPage()))
	}
	return tmp_html
}

//
//直接写死代码，返回html片段。
func (page *Page) String() string {
	return fmt.Sprintf("<p>page:[ PageSize: %d, TotalCount: %d, CurrentPage: %d ] </p> &nbsp;", page.PageSize, page.TotalCount, page.CurrentPage)
}
