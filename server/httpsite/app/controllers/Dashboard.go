package controllers

import (
	"devel.mephi.ru/dyokunev/dc-thermal-logger/server/httpsite/app/models"
	"devel.mephi.ru/dyokunev/dc-thermal-logger/server/httpsite/app"
	"github.com/revel/revel"
)

type Dashboard struct {
	*revel.Controller
}

func (c Dashboard) Page() revel.Result {
	var err error

	c.RenderArgs["rawRecords"],err = models.RawRecord.Select(app.DB)
	if err != nil {
		revel.ERROR.Printf("%s", err.Error())
	}

	return c.Render()
}
