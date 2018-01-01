package cache

import "github.com/astaxie/beego/orm"
import (
	"wx_article/models"
)

var allUrls = make(map[string]bool, 5000)

func init()  {
	var articles []models.Article
	o := orm.NewOrm()
	_, err := o.Raw("SELECT url FROM wx_article").QueryRows(&articles)
	if err != nil {
		panic(err)
	}
	for _, article := range articles {
		Add(article.Url)
	}
}

func Add(url string) (b bool) {
	if allUrls[url] {
		return false
	}
	allUrls[url] = true
	return true
}

func Size() (siez int) {
	return len(allUrls)
}