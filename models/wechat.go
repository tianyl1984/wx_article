package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	Id int64 `PK`
	Name string
	Publisher string
}

type Article struct {
	Id int64 `PK`
	AppId int64 `orm:"column(id_app)"`
	Title string
	Url string
	Digest string
	PublishTime time.Time `orm:"column(publishTime)"`
	Uuid string
	HasRead bool `orm:"column(hasRead)"`
	OfflineUrl string `orm:"column(offlineUrl)"`
}

func init()  {
	orm.Debug = true
	err1 := orm.RegisterDriver("mysql",orm.DRMySQL)
	if err1 != nil {
		panic(err1)
	}
	//err := orm.RegisterDataBase("default","mysql","root:tyl123@tcp(192.168.0.111:3306)/weixin?charset=utf8",3,3)
	err := orm.RegisterDataBase("default","mysql","root:tyl123@tcp(127.0.0.1:3306)/weixin?charset=utf8",3,3)
	if err != nil {
		panic(err)
	}
	orm.RegisterModelWithPrefix("wx_",new(App),new(Article))
}
