package main

import (
	"fmt"
	"github.com/astaxie/beego"
	_ "wx_article/models"
	_ "wx_article/routers"
)

func main() {
	fmt.Println("start")
	beego.Run()
}
