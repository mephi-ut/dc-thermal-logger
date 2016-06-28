package models

import (
	"time"
)

//go:generate reform

//reform:raw_records
type rawRecord struct {
	ModelBase

	Id           int       `reform:"id,pk"`
	Date         time.Time `reform:"date"`
	RawSensorId  int       `reform:"raw_sensor_id"`
	RawChannelId int       `reform:"raw_channel_id"`
	RawValue     int       `reform:"raw_value"`
}

func (r *rawRecord) ToHistoryRecords() (result map[AggregationType]historyRecord, err error) {
	result = make(map[AggregationType]historyRecord)

	historyRecordAsIs := historyRecord{
		Date:     MyTime(r.Date),
		SensorId: sensorIdMap[r.RawSensorId][r.RawChannelId],
		RawValue: float32(r.RawValue),
	}

	if historyRecordAsIs.SensorId == 0 {
		err = ErrInvalidSensorId
		return
	}

	err = historyRecordAsIs.ConvertValue()
	if err != nil {
		return
	}

	for _, aggregationType := range aggregationTypes {
		historyRecord := historyRecordAsIs
		historyRecord.AggregationType = aggregationType
		historyRecord.FixDate()

		result[aggregationType] = historyRecord
	}

	return
}
