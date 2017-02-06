package controllers

import "github.com/astaxie/beego"

type RequestResult struct {
	Result bool `json:"result"`
}

type BaseController struct {
	beego.Controller
}

func (this *BaseController) serveOk() {
	rr := RequestResult{Result: true}
	this.Data["json"] = &rr
	this.ServeJSON()
}
