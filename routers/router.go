package routers

import (
	"github.com/astaxie/beego"
	"wx_article/controllers"
)

func init()  {
	beego.Router("/article/save",&controllers.ArticalController{},"*:Save")
	beego.Router("/article/setRead",&controllers.ArticalController{},"*:SetRead")
	beego.Router("/article/list",&controllers.ArticalController{},"*:List")
}