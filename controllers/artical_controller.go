package controllers

import (
	"github.com/astaxie/beego"
	"wx_article/models"
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"time"
)

type ArticalController struct {
	beego.Controller
}

type RequestResult struct{
	Result bool `json:"result"`
}

func (this *ArticalController) Save()  {
	var msgs []models.AppMsg
	json.Unmarshal(this.Ctx.Input.RequestBody,&msgs)

	//TODO 批量操作

	for _,msg := range msgs {
		app := models.App{}
		app.Name = msg.AppName
		app.Publisher = msg.PublisherUsername
		o := orm.NewOrm()
		err := o.Read(&app,"Publisher")
		if err == orm.ErrNoRows {
			id,er := o.Insert(&app)
			if er != nil{
				panic("save app error")
			}
			app.Id = id
		}
		//TODO 支持更新名称
		artical := models.Article{}
		artical.AppId = app.Id
		artical.Digest = msg.Digest
		artical.HasRead = false
		artical.Title = msg.Title
		artical.Url = msg.Url
		if msg.PublishTime != "" {
			t,e := time.Parse("2006-01-02 15:04:05",msg.PublishTime)
			if e != nil {
				panic(e)
			}
			artical.PublishTime = t
		}
		_, err2 := o.Insert(&artical)
		if err2 != nil{
			panic("save article error")
		}
	}

	rr := RequestResult{Result:true}
	this.Data["json"] = &rr
	this.ServeJSON()
}

func (this *ArticalController) SetRead(){
	this.Ctx.WriteString("OK")
}

