package main

import (
	"fmt"
	"github.com/astaxie/beego"
	_ "wx_article/models"
	_ "wx_article/routers"
	_ "wx_article/cache"
)

func main() {
	fmt.Println("start")
	beego.Run()
}
