package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"time"
	"wx_article/models"
	"wx_article/util"
	"wx_article/cache"
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
		if !cache.Add(msg.Url) {
			continue
		}
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
			loc, _ := time.LoadLocation("Local")
			//t, e := time.Parse("2006-01-02 15:04:05", msg.PublishTime)
			t, e := time.ParseInLocation("2006-01-02 15:04:05", msg.PublishTime, loc)
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
	appId, err := this.GetInt("appId", 0)
	if err != nil {
		panic(err)
	}

	var pageSize = page.PageSize
	var pageNum = page.PageNum
	o := orm.NewOrm()
	var countSql = "select count(1) as total_size from wx_article where hasRead = ? "
	if appId > 0 {
		countSql += " and id_app = ? "
	}

	if appId > 0 {
		err = o.Raw(countSql, false, appId).QueryRow(&page)
	} else {
		err = o.Raw(countSql, false).QueryRow(&page)
	}

	if err != nil {
		panic(err)
	}

	page.TotalPage = page.TotalSize / pageSize
	if page.TotalSize%pageSize > 0 {
		page.TotalPage = page.TotalPage + 1
	}
	qs := o.QueryTable("wx_article")
	qs = qs.Filter("hasRead", false)
	if appId > 0 {
		qs = qs.Filter("AppId", appId)
	}
	qs = qs.OrderBy("-id").Limit(page.PageSize, page.PageSize*(pageNum-1))

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

func (this *ArticleController) ListDelArticle() {
	o := orm.NewOrm()
	var delMsgs []*models.DeleteMessage
	if _, err := o.QueryTable("wx_delete_message").All(&delMsgs); err != nil {
		panic("find delete message error")
	}

	var apps []*models.App
	if _, err := o.QueryTable("wx_app").All(&apps); err != nil {
		panic("find app error")
	}
	appMap := make(map[string]models.App)
	for _, app := range apps {
		appMap[app.Publisher] = *app
	}

	var articleResult []ArticleResult = make([]ArticleResult, 0)
	var existArticleMap = make(map[int64]string)
	for _, delMsg := range delMsgs {
		var articles []*models.Article
		_, err := o.QueryTable("wx_article").Filter("AppId", appMap[delMsg.Publisher].Id).
			Filter("HasRead", false).Filter("PublishTime__lte", delMsg.CreateTime).All(&articles)
		if err != nil {
			panic("find article error")
		}
		for _, article := range articles {
			if _, ok := existArticleMap[article.Id]; ok {
				continue
			}
			ar := ArticleResult{Url: article.Url, Id: article.Id, Title: article.Title,
				PublishTime: util.Time(article.PublishTime),
				AppId:       article.AppId, AppName: appMap[delMsg.Publisher].Name}
			articleResult = append(articleResult, ar)
			existArticleMap[article.Id] = ""
		}
	}

	this.Data["json"] = &articleResult
	this.ServeJSON()
}

func (this *ArticleController) ReadDelArticle() {
	o := orm.NewOrm()
	var delMsgs []*models.DeleteMessage
	if _, err := o.QueryTable("wx_delete_message").All(&delMsgs); err != nil {
		panic("find delete message error")
	}

	var apps []*models.App
	if _, err := o.QueryTable("wx_app").All(&apps); err != nil {
		panic("find app error")
	}
	appMap := make(map[string]int64)
	for _, app := range apps {
		appMap[app.Publisher] = app.Id
	}

	for _, delMsg := range delMsgs {
		if _, err := o.QueryTable("wx_article").Filter("AppId", appMap[delMsg.Publisher]).
			Filter("HasRead", false).Filter("PublishTime__lte", delMsg.CreateTime).
			Update(orm.Params{"hasRead": true, "readTime": time.Now()}); err != nil {
			panic("update article error")
		}
	}
	if _, err := o.QueryTable("wx_delete_message").Filter("id__gt", 0).Delete(); err != nil {
		panic("remove delete message error")
	}

	this.serveOk()
}

func (this *ArticleController) Test()  {
	//this.Data["json"] = cache.Add("http://mp.weixin.qq.com/s?__biz=MzA4MTQ4NjQzMw==&mid=2652708426&idx=1&sn=6691dd192ebb91db141bba11b2358481&chksm=847d8944b30a0052a1f6682d92da4e1af58111b7a7fdda9e2e7f1e400fba0f6afeee87689281&scene=0#rd")
	this.Data["json"] = cache.Size()
	this.ServeJSON()
}