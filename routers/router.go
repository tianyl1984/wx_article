package routers

import (
	"github.com/astaxie/beego"
	"wx_article/controllers"
)

func init() {
	beego.Router("/article/save", &controllers.ArticleController{}, "*:Save")
	beego.Router("/article/setRead", &controllers.ArticleController{}, "*:SetRead")
	beego.Router("/article/addFavorite", &controllers.ArticleController{}, "*:AddFavorite")
	beego.Router("/article/list", &controllers.ArticleController{}, "*:List")
	beego.Router("/message/addDeleteMessage", &controllers.MessageController{}, "*:AddDeleteMessage")
}
