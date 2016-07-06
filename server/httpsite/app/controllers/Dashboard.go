package controllers

import (
	//"fmt"
	"time"
	"strconv"
	"devel.mephi.ru/dyokunev/dc-thermal-logger/server/httpsite/app/models"
//	"devel.mephi.ru/dyokunev/dc-thermal-logger/server/httpsite/app"
	"github.com/revel/revel"
	"gopkg.in/reform.v1"
//	"golang.org/x/net/websocket";
)

const (
	SENSORS_COUNT = 256
)

type sensorInfo struct {
	Id		int
	FullName	string
	GroupName	string
	Name		string
	Value		float32
	LastTimestamp	time.Time
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

		"ServerRack5":		groupInfo{DefaultSensorId: 66,	SensorIds: []int{ 64, 65, 66 }},
		"ServerRack4":		groupInfo{DefaultSensorId: 69,	SensorIds: []int{ 67, 68, 69 }},
	}

func (c Dashboard) page() {
	sensors := map[int]sensorInfo{}

	for groupName,groupInfo := range groups {
		for _,sensorId := range groupInfo.SensorIds {
			sensor := sensorInfo{}
			sensor.Id        = sensorId
			sensor.GroupName = groupName
			sensor.Name      = models.SensorNameMap[sensorId]
			sensor.FullName  = models.SensorFullNameMap[sensorId]

			if sensor.Name == "" {
				continue
			}

			theLastHistoryRecord,err := models.HistoryRecord.Order("date", "DESC").Where("counter > 20").First(models.HistoryRecordFilter{SensorId: sensorId, AggregationType: models.AGGR_MINUTE})
			if err != nil {
				if err != reform.ErrNoRows {
					revel.ERROR.Printf("%s", err.Error())
				}
				continue
			}
			sensor.LastTimestamp = time.Time(theLastHistoryRecord.Date)
			sensor.Value         = float32(int((theLastHistoryRecord.ConvertedValue-273.15)*10))/10

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

func (c Dashboard) minimal() {
	sensors := map[int]sensorInfo{}

	for _,groupInfo := range groups {
		for _,sensorId := range groupInfo.SensorIds {
			if sensorId >= 16 {
				continue // Only Air Conditioning
			}

			sensor := sensorInfo{}
			sensor.FullName  = models.SensorFullNameMap[sensorId]

			theLastHistoryRecord,err := models.HistoryRecord.Order("date", "DESC").Where("counter > 20").First(models.HistoryRecordFilter{SensorId: sensorId, AggregationType: models.AGGR_MINUTE})
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
}

func (c Dashboard) Minimal() revel.Result {
	c.minimal()
	return c.Render()
}

func (c Dashboard) TSJson() revel.Result {
	return c.RenderJson(map[string]string{ "ts": time.Now().Format("2006-01-02 15:04:05") })
}

func (c Dashboard) MapHtml() revel.Result {
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

func (c Dashboard) MinimalJson() revel.Result {
	c.minimal()
	sensors := map[string]sensorInfo{}
	for k,v := range c.RenderArgs["sensors"].(map[int]sensorInfo) {
		sensors[strconv.Itoa(k)] = v
	}
	c.RenderArgs["sensors"] = sensors
	return c.RenderJson(c.RenderArgs)
}

/*
func (c Dashboard) Websocket(ws *websocket.Conn) revel.Result {
	revel.TRACE.Printf("New WS connection")
	var JSON = Codec{jsonMarshal, jsonUnmarshal}

	recvMessages := make(chan message)

	go func() {
		var msg message
		for {
			err := websocket.JSON.Receive(ws, &msg)
			if (err != nil) {	// disconnected
				revel.TRACE.Printf("WS is closed")
				close(recvMessages)
				return
			}
			recvMessages <- msg
		}
	}()

	for {
		select {
			case msg, ok := <-recvMessages:
				if (!ok) { // If the channel is closed
					revel.TRACE.Printf("recvMessages channel is closed")
					return nil
				}

				if (msg.Type == MSGTYPE_PING) {
					msg.Type = MSGTYPE_PONG
					websocket.JSON.Send(ws, &msg)
					break
				}

				switch (msgWords)

				revel.INFO.Printf("Got a WS message: %v\n", msg)
		}
	}

	return nil
}
*/

