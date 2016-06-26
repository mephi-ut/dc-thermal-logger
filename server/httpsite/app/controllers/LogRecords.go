package controllers

import (
	"devel.mephi.ru/dc-thermal-logger/server/httpsite/app/helpers"
	"github.com/revel/revel"
)

type LogRecords struct {
	*revel.Controller
}

func (c LogRecords) Dashboard() revel.Result {
	c.RenderArgs["logRecords"],_ = helpers.LogRecord.Select("")

	return c.Render()
}
