package controllers

import (
	"github.com/astaxie/beego/orm"
	"wx_article/models"
)

type AppController struct {
	BaseController
}

func (this *AppController) List() {
	o := orm.NewOrm()
	var apps []models.App
	o.QueryTable("wx_app").All(&apps)

	var result []map[string]interface{}

	for _, app := range apps {
		result = append(result, map[string]interface{}{"name": app.Name, "id": app.Id})
	}

	this.Data["json"] = &result
	this.ServeJSON()
}
