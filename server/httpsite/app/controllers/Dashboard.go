package controllers

import (
	//"fmt"
	"time"
	"strconv"
	"sync"
	"devel.mephi.ru/dyokunev/dc-thermal-logger/server/httpsite/app/models"
//	"devel.mephi.ru/dyokunev/dc-thermal-logger/server/httpsite/app"
	"github.com/revel/revel"
	"github.com/revel/revel/cache"
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

		"ServerRack5":		groupInfo{DefaultSensorId: 66,	SensorIds: []int{ 64,  65,  66 }},
		"ServerRack4":		groupInfo{DefaultSensorId: 69,	SensorIds: []int{ 67,  68,  69 }},
		"ServerRack3":		groupInfo{DefaultSensorId: 74,	SensorIds: []int{ 72,  73,  74 }},
		"ServerRack2":		groupInfo{DefaultSensorId: 76,	SensorIds: []int{ 75,  76, -77 }},
		"ServerRack1":		groupInfo{DefaultSensorId: 85,	SensorIds: []int{ 84,  83,  85 }},
		"ServerRack0":		groupInfo{DefaultSensorId: 81,	SensorIds: []int{ 80,  82,  81 }},
		"ServerRack6":		groupInfo{DefaultSensorId: 89,	SensorIds: []int{ 89,  88,  90 }},
		"ServerRack7":		groupInfo{DefaultSensorId: 93,	SensorIds: []int{ 93,  91,  92 }},
		"ServerRack8":		groupInfo{DefaultSensorId: 99,	SensorIds: []int{ 100, 99,  101 }},
		"ServerRack9":		groupInfo{DefaultSensorId: 97,	SensorIds: []int{ 100, 97,  96  }},
		"ServerRack10":		groupInfo{DefaultSensorId: 106,	SensorIds: []int{ 106, 105, 104 }},
		"ServerRack11":		groupInfo{DefaultSensorId: 109,	SensorIds: []int{ 109, 108, 107 }},
		"ServerRack13":		groupInfo{DefaultSensorId: 112,	SensorIds: []int{ 112, 113, 114 }},
		"ServerRack14":		groupInfo{DefaultSensorId: 118,	SensorIds: []int{ 116, 118, 117 }},
	}

func (c Dashboard) get(key, defaultValue string) (result string) {
	c.Params.Bind(&result, key)
	if result == "" {
		return defaultValue
	}

	return
}

func (c Dashboard) considerGetParameters() {
	for k,v := range map[string]string{
			"mobileAutoupdate":	"toggle",
			"mobileShowNames":	"true",
	} {
		c.RenderArgs[k] = c.get(k, v);
	}
}

func (c Dashboard) page() {
	c.RenderArgs["groups" ] = groups
	sensors := map[int]sensorInfo{}

	if err := cache.Get("sensors", &sensors); err == nil {
		c.RenderArgs["sensors"] = sensors
		c.considerGetParameters()
		return;
	}

	mutex := sync.Mutex{}
	sem := make(chan bool, 16)
	for groupName,grpInfo := range groups {
		sem <- true
		go func(groupName string, groupInfo groupInfo) {
			defer func() { <-sem }()
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

				mutex.Lock()
				sensors[sensorId] = sensor
				mutex.Unlock()
			}
		}(groupName, grpInfo);
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
	go cache.Set("sensors", sensors, 10*time.Second)

	c.RenderArgs["sensors"] = sensors
	c.considerGetParameters()
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
	c.considerGetParameters()
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

func (c Dashboard) MapMainDCHtml() revel.Result {
	c.page()
	return c.Render()
}

func (c Dashboard) MapReserveDCHtml() revel.Result {
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

