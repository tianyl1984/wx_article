package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"time"
	"wx_article/models"
	"wx_article/util"
)

type ArticleController struct {
	BaseController
}

type ArticleResult struct {
	Url         string    `json:"url"`
	Title       string    `json:"title"`
	Id          int64     `json:"id"`
	PublishTime util.Time `json:"publishTime"`
	AppName     string    `json:"appName"`
	AppId       int64     `json:"appId"`
}

func (this *ArticleController) Save() {
	var msgs []models.AppMsg
	json.Unmarshal(this.Ctx.Input.RequestBody, &msgs)

	//TODO 批量操作

	for _, msg := range msgs {
		app := models.App{}
		app.Name = msg.AppName
		app.Publisher = msg.PublisherUsername
		o := orm.NewOrm()
		err := o.Read(&app, "Publisher")
		if err == orm.ErrNoRows {
			id, er := o.Insert(&app)
			if er != nil {
				panic("save app error")
			}
			app.Id = id
		}
		//TODO 支持更新名称
		artical := models.Article{}
		artical.AppId = app.Id
		artical.Digest = msg.Digest
		artical.HasRead = false
		artical.Favorite = false
		artical.Title = msg.Title
		artical.Url = msg.Url
		if msg.PublishTime != "" {
			t, e := time.Parse("2006-01-02 15:04:05", msg.PublishTime)
			if e != nil {
				panic(e)
			}
			artical.PublishTime = t
		}
		_, err2 := o.Insert(&artical)
		if err2 != nil {
			panic("save article error")
		}
	}

	this.serveOk()
}

func (this *ArticleController) SetRead() {
	articleId, err := this.GetInt64("articleId")
	if err != nil {
		panic(err)
	}
	o := orm.NewOrm()
	if _, ok := o.Raw("update wx_article set hasRead = ?,readTime = ? where id = ? ", true, time.Now(), articleId).Exec(); ok != nil {
		panic(ok)
	}
	this.serveOk()
}

func (this *ArticleController) AddFavorite() {
	articleId, err := this.GetInt64("articleId")
	if err != nil {
		panic(err)
	}
	o := orm.NewOrm()
	if _, ok := o.Raw("update wx_article set favorite = ?,favTime = ? where id = ? ", true, time.Now(), articleId).Exec(); ok != nil {
		panic(ok)
	}
	this.serveOk()
}

func (this *ArticleController) List() {
	page := models.Page{}
	if pageSize, err := this.GetInt("pageSize", 15); err != nil {
		panic(err)
	} else {
		page.PageSize = pageSize
	}
	if pageNum, err := this.GetInt("pageNum", 1); err != nil {
		panic(err)
	} else {
		page.PageNum = pageNum
	}

	var pageSize = page.PageSize
	var pageNum = page.PageNum
	o := orm.NewOrm()
	if err := o.Raw("select count(1) as total_size from wx_article where hasRead = ? ", false).QueryRow(&page); err != nil {
		panic(err)
	}
	page.TotalPage = page.TotalSize / pageSize
	if page.TotalSize%pageSize > 0 {
		page.TotalPage = page.TotalPage + 1
	}
	qs := o.QueryTable("wx_article")
	qs = qs.Filter("hasRead", false).OrderBy("-id").Limit(page.PageSize, page.PageSize*(pageNum-1))

	var maps []orm.Params
	if _, err := qs.Values(&maps); err != nil {
		panic(err)
	}

	var appIds []int64
	var appIdMap map[int64]interface{} = make(map[int64]interface{}, len(maps))
	var articles []ArticleResult
	for _, m := range maps {
		var temp = ArticleResult{}
		if id, ok := m["Id"].(int64); ok {
			temp.Id = id
		}
		if str, ok := m["Title"].(string); ok {
			temp.Title = str
		}
		if str, ok := m["Url"].(string); ok {
			temp.Url = str
		}
		if t, ok := m["PublishTime"].(time.Time); ok {
			temp.PublishTime = util.Time(t)
		}
		if id, ok := m["AppId"].(int64); ok {
			if _, ok := appIdMap[id]; !ok {
				appIdMap[id] = ""
				appIds = append(appIds, id)
			}
			temp.AppId = id
		}
		articles = append(articles, temp)
	}
	if len(appIds) > 0 {
		appQs := o.QueryTable("wx_app")
		appMap := make(map[int64]string)
		appList := []orm.Params{}
		appQs.Filter("Id__in", appIds).Values(&appList)
		for _, app := range appList {
			if id, ok := app["Id"].(int64); ok {
				if str, ok := app["Name"].(string); ok {
					appMap[id] = str
				}
			}
		}

		for i, _ := range articles {
			articles[i].AppName = appMap[articles[i].AppId]
		}
	}

	page.Data = articles

	this.Data["json"] = &page
	this.ServeJSON()
}
