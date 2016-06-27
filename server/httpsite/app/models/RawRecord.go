package models

import (
	"time"
)

//go:generate reform

//reform:raw_records
type rawRecord struct {
	ModelBase

	Id        int       `reform:"id,pk"`
	Date      time.Time `reform:"date"`
	SensorId  int       `reform:"sensor_id"`
	ChannelId int       `reform:"channel_id"`
	RawValue  int       `reform:"raw_value"`
}
