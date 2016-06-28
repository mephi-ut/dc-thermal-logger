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
	Name  string
	Value float32
}

type Dashboard struct {
	*revel.Controller
}

var groups = map[string][]int {
		"AirConditioning0": []int{ 1 },
		"AirConditioning2": []int{ 2 },
	}

func (c Dashboard) Page() revel.Result {
	sensors := map[int]sensorInfo{}

	for _,group := range groups {
		for _,sensorId := range group {
			sensor := sensorInfo{}
			sensor.Name  = models.SensorNameMap[sensorId]

			if sensor.Name == "" {
				continue
			}

			theLastHistoryRecord,err := models.HistoryRecord.Order("date", "DESC").First(app.DB, models.HistoryRecordFilter{SensorId: sensorId, AggregationType: models.AGGR_SECOND})
			if err != nil {
				if err != reform.ErrNoRows {
					revel.ERROR.Printf("%s", err.Error())
				}
				continue
			}
			sensor.Value = theLastHistoryRecord.ConvertedValue

			sensors[sensorId] = sensor
		}
	}

	c.RenderArgs["sensors"] = sensors
	c.RenderArgs["groups" ] = groups

	return c.Render()
}
