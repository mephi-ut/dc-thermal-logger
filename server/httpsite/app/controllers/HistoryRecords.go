package controllers

import (
	//"fmt"
	"devel.mephi.ru/dyokunev/dc-thermal-logger/server/httpsite/app/models"
	//"devel.mephi.ru/dyokunev/dc-thermal-logger/server/httpsite/app"
	"github.com/revel/revel"
)

type HistoryRecords struct {
	*revel.Controller
}

func (c HistoryRecords) Find(historyRecord models.HistoryRecordFilter, order string, limit int) revel.Result {
	revel.INFO.Printf("%v", historyRecord)

	historyRecords,err := models.HistoryRecord.Order(order).Limit(limit).Select(historyRecord)
	if err != nil {
		revel.ERROR.Printf("%v", err.Error())
	}
	for i,_ := range historyRecords {
		historyRecords[i].ConvertedValue -= 273.15
	}
	c.RenderArgs["historyRecords"] = historyRecords

	return c.RenderJson(c.RenderArgs)
}
