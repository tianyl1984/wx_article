package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
	"wx_article/models"
)

type MessageController struct {
	BaseController
}

func (this *MessageController) AddDeleteMessage() {
	var delMsgs []models.DelMsg
	json.Unmarshal(this.Ctx.Input.RequestBody, &delMsgs)
	for _, msg := range delMsgs {
		saveMsg := models.DeleteMessage{}
		saveMsg.Publisher = msg.Talker
		tmInt64, _ := strconv.ParseInt(msg.CreateTime, 10, 64)
		tm := time.Unix(tmInt64/1000, 0)
		saveMsg.CreateTime = tm
		o := orm.NewOrm()
		_, err := o.Insert(&saveMsg)
		if err != nil {
			panic("save deleteMessage error")
		}
	}
	this.serveOk()
}

func (this *MessageController) List() {
	o := orm.NewOrm()

	timeStr := "1486425770000"
	timeInt64, _ := strconv.ParseInt(timeStr, 10, 64)
	tm := time.Unix(timeInt64/1000, 0)

	toSave := models.DeleteMessage{}
	toSave.Publisher = "1234"
	toSave.CreateTime = tm
	o.Insert(&toSave)

	var msgs []*models.DeleteMessage
	o.QueryTable("wx_delete_message").All(&msgs)
	for _, msg := range msgs {
		fmt.Println(msg.CreateTime)
	}

	this.serveOk()
}
