package main

import (
	"fmt"
	"github.com/astaxie/beego"
	_ "wx_article/routers"
	_ "wx_article/models"
)

func main()  {
	fmt.Println("start")
	beego.Run()
}
