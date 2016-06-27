package controllers

import (
	//"fmt"
	"devel.mephi.ru/dyokunev/dc-thermal-logger/server/httpsite/app/models"
	"devel.mephi.ru/dyokunev/dc-thermal-logger/server/httpsite/app"
	"github.com/revel/revel"
	"gopkg.in/reform.v1"
)

const (
	SENSORS_COUNT = 256
)

type sensorInfo struct {
	Temp float32
}

type Dashboard struct {
	*revel.Controller
}

func (c Dashboard) Page() revel.Result {
	var sensors [SENSORS_COUNT+1]sensorInfo

	for sensorId := 1; sensorId <= SENSORS_COUNT; sensorId++ {
		sensor := &sensors[sensorId]

		theLastHistoryRecord,err := models.HistoryRecord.Order("date", "DESC").First(app.DB, models.HistoryRecordFilter{SensorId: sensorId})
		if err != nil {
			if err != reform.ErrNoRows {
				revel.ERROR.Printf("%s", err.Error())
			}
			continue
		}
		sensor.Temp = theLastHistoryRecord.ConvertedValue
	}

	c.RenderArgs["sensors"] = sensors

	return c.Render()
}
