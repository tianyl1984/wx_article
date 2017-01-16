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

type ArticleResult struct {
	Url string `json:"url"`
	Title string `json:"title"`
	Id int64 `json:"id"`
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

	this.serveOk()
}

func (this *ArticalController) SetRead(){
	articleId,err := this.GetInt64("articleId")
	if err != nil {
		panic(err)
	}
	o := orm.NewOrm()
	if _,ok := o.Raw("update wx_article set hasRead = ? where id = ? ",true, articleId).Exec();ok != nil {
		panic(ok)
	}
	this.serveOk()
}

func (this *ArticalController) List()  {
	page := models.Page{}
	if pageSize,err := this.GetInt("pageSize",15); err != nil {
		panic(err)
	}else{
		page.PageSize = pageSize
	}
	if pageNum,err := this.GetInt("pageNum",1); err != nil {
		panic(err)
	}else{
		page.PageNum = pageNum
	}

	var pageSize = page.PageSize
	var pageNum = page.PageNum
	o := orm.NewOrm()
	if err := o.Raw("select count(1) as total_size from wx_article where hasRead = ? ",false).QueryRow(&page);err != nil{
		panic(err)
	}
	page.TotalPage = page.TotalSize/pageSize
	if page.TotalSize%pageSize > 0 {
		page.TotalPage = page.TotalPage + 1
	}
	qs := o.QueryTable("wx_article")
	qs = qs.Filter("hasRead",false).OrderBy("-id").Limit(page.PageSize,page.PageSize*(pageNum - 1))

	var maps []orm.Params
	if _,err := qs.Values(&maps);err != nil{
		panic(err)
	}

	var articles []ArticleResult
	for _,m := range maps {
		var temp = ArticleResult{}
		if id,ok := m["Id"].(int64); ok {
			temp.Id = id
		}
		if str,ok := m["Title"].(string); ok {
			temp.Title = str
		}
		if str,ok := m["Url"].(string); ok {
			temp.Url = str
		}
		articles = append(articles,temp)
	}

	page.Data = articles

	this.Data["json"] = &page
	this.ServeJSON()
}

func (this *ArticalController) serveOk() {
	rr := RequestResult{Result:true}
	this.Data["json"] = &rr
	this.ServeJSON()
}
