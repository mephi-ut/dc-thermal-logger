package controllers

import (
	//"fmt"
	"strconv"
	"devel.mephi.ru/dyokunev/dc-thermal-logger/server/httpsite/app/models"
	//"devel.mephi.ru/dyokunev/dc-thermal-logger/server/httpsite/app"
	"github.com/revel/revel"
)

type HistoryRecords struct {
	*revel.Controller
}

func (c HistoryRecords) Find(historyRecord models.HistoryRecordFilter, order string, limit int) revel.Result {
	revel.INFO.Printf("%v", historyRecord)

	scope := models.HistoryRecord.Order(order)

	if (historyRecord.SensorId == 0) {
		scope  = scope.Where("sensor_id IN (1,2,3)")
		limit *= 3
	}

	historyRecords,err := scope.Limit(limit).Select(historyRecord)
	if err != nil {
		revel.ERROR.Printf("%v", err.Error())
	}
	for i,_ := range historyRecords {
		historyRecords[i].ConvertedValue -= 273.15
	}

	if (historyRecord.SensorId == 0) {
		historyRecordsBySensor := [4][]models.HistoryRecordFilter{}	// TODO: "Filter" shouldn't be used as just the original type
		for _,r := range historyRecords {
			historyRecordsBySensor[ r.SensorId ] = append(historyRecordsBySensor[ r.SensorId ], models.HistoryRecordFilter(r))
		}
		for i:=1; i<=3; i++ {
			c.RenderArgs["historyRecords"+strconv.Itoa(i)] = historyRecordsBySensor[i]
		}
	} else {
		c.RenderArgs["historyRecords"] = historyRecords
	}

	return c.RenderJson(c.RenderArgs)
}
