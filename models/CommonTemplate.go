package models

import (
	"os"
	"path/filepath"
	"github.com/astaxie/beego/logs"
	"log"
	"io/ioutil"
	"html/template"
	"github.com/astaxie/beego"
)

//通用按照模板生成文件方法。
func GenFileByTemplate(fileName, fileDir, templateFile string, data map[string]interface{}) {

	//https://stackoverflow.com/questions/32551811/read-file-as-template-execute-it-and-write-it-back
	//读取一个模板文件，然后执行在写入到一个文件。

	//增加自定义函数。
	funcs := make(map[string]interface{})
	funcs["str2html"] = beego.Str2html //使用到这个函数
	funcs["date"] = beego.Date         //使用日期格式化。

	//如果没有文件dir，则创建文件夹。
	_, err1 := os.Stat(fileDir)
	if os.IsNotExist(err1) {
		os.MkdirAll(fileDir, os.ModePerm)
	}

	fileBase := filepath.Base(fileName)
	f, err := os.Create(fileName)
	logs.Info("################# fileDir #################", fileDir)
	logs.Info("################# fileBase #################", fileBase)
	defer f.Close() //最后关闭文件。
	logs.Info("write to file.", f.Name())
	if err != nil {
		log.Println("create file: ", err)
		return
	}
	//每次使用参数传递的模板文件创建。
	tmplContent, err2 := ioutil.ReadFile(templateFile)
	if err2 != nil {
		logs.Error("ReadFile ERROR :", err2.Error())
	}

	tmp1 := template.New("letter")
	tmp1.Funcs(funcs)
	t := template.Must(tmp1.Parse(string(tmplContent)))
	err = t.Execute(f, data)
	if err != nil {
		log.Print("execute: ", err)
		return
	}
}
