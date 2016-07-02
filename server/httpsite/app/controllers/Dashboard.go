package controllers

import (
	//"fmt"
	"strconv"
	"devel.mephi.ru/dyokunev/dc-thermal-logger/server/httpsite/app/models"
	"devel.mephi.ru/dyokunev/dc-thermal-logger/server/httpsite/app"
	"github.com/revel/revel"
	"gopkg.in/reform.v1"
)

const (
	SENSORS_COUNT = 256
)

type sensorInfo struct {
	GroupName	string
	Name		string
	Value		float32
}
type groupInfo struct {
	DefaultSensorId	  int
	SensorIds	[]int
}

type Dashboard struct {
	*revel.Controller
}

var groups = map[string]groupInfo {
		"AirConditioning0":	groupInfo{DefaultSensorId: 1,	SensorIds: []int{ 1 }},
		"AirConditioning2":	groupInfo{DefaultSensorId: 2,	SensorIds: []int{ 2 }},
		"AirConditioning1":	groupInfo{DefaultSensorId: 3,	SensorIds: []int{ 3 }},

		"ServerRack5":		groupInfo{DefaultSensorId: 66,	SensorIds: []int{ 64, 66, 65 }},
		"ServerRack4":		groupInfo{DefaultSensorId: 69,	SensorIds: []int{ 67, 69, 68 }},
	}

func (c Dashboard) page() {
	sensors := map[int]sensorInfo{}

	for groupName,groupInfo := range groups {
		for _,sensorId := range groupInfo.SensorIds {
			sensor := sensorInfo{}
			sensor.GroupName = groupName
			sensor.Name      = models.SensorNameMap[sensorId]

			if sensor.Name == "" {
				continue
			}

			theLastHistoryRecord,err := models.HistoryRecord.Order("date", "DESC").First(app.DB, models.HistoryRecordFilter{SensorId: sensorId, AggregationType: models.AGGR_MINUTE})
			if err != nil {
				if err != reform.ErrNoRows {
					revel.ERROR.Printf("%s", err.Error())
				}
				continue
			}
			sensor.Value = float32(int((theLastHistoryRecord.ConvertedValue-273.15)*10))/10

			sensors[sensorId] = sensor
		}
	}

	c.RenderArgs["sensors"] = sensors
	c.RenderArgs["groups" ] = groups
}

func (c Dashboard) Page() revel.Result {
	c.page()
	return c.Render()
}

func (c Dashboard) PageJson() revel.Result {
	c.page()
	sensors := map[string]sensorInfo{}
	for k,v := range c.RenderArgs["sensors"].(map[int]sensorInfo) {
		sensors[strconv.Itoa(k)] = v
	}
	c.RenderArgs["sensors"] = sensors
	return c.RenderJson(c.RenderArgs)
}

