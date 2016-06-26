package models

import (
	"time"
)

//go:generate reform

//reform:log_records
type LogRecord struct {
	ModelBase

	Date           time.Time `reform:"date"`
	SensorId       int       `reform:"sensor_id"`
	ChannelId      int       `reform:"channel_id"`
	Value          int       `reform:"value"`
	ConvertedValue int       `reform:"converted_value"`
}

func (obj *LogRecord) AfterFind() error {
	return obj.Init(obj)
}

func init() {
	modelRegister(LogRecord{}, LogRecordView)
}

