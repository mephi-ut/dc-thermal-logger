package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

const (
	AGGR_UNKNOWN AggregationType = iota
	AGGR_SECOND
	AGGR_MINUTE
	AGGR_5MINUTES
	AGGR_HOUR
	AGGR_DAY
	AGGR_WEEK
)

var (
	ErrInvalidSensorId = errors.New("Invalid sensor/channel ID")
)

//go:generate reform

type AggregationType int
type MyTime time.Time

//reform:history_records
type historyRecord struct {
	ModelBase

	Id              int             `reform:"id,pk"`
	Date            MyTime          `reform:"date"`
	AggregationType AggregationType `reform:"aggregation_period"`
	SensorId        int             `reform:"sensor_id"`
	RawValue        float32         `reform:"raw_value"`
	ConvertedValue  float32         `reform:"converted_value"`
	Counter         int             `reform:"counter"`
}

func (t MyTime) Format(fmt string) string {
	return (time.Time)(t).Format(fmt)
}
func (t MyTime) Unix() int64 {
	return (time.Time)(t).Unix()
}
func (t *MyTime) Scan(src interface{}) error {
	*t = MyTime(src.(time.Time))
	return nil
}
func (t MyTime) Value() (driver.Value, error) {
	return []byte(t.Format("2006-01-02 15:04:05")), nil
}

func (t *MyTime) MarshalJSON() ([]byte, error) {
	return []byte("\""+t.Format("2006-01-02 15:04:05")+"\""), nil
}

var aggregationTypes = []AggregationType{AGGR_SECOND, AGGR_MINUTE, AGGR_5MINUTES, AGGR_HOUR, AGGR_DAY, AGGR_WEEK}
var sensorIdMap = map[int]map[int]int{
	1: {0: 1, 1: 2, 2: 3, 3: 4},
}

type convertionFunction func(rawValue float32) (convertedValue float32)

var valueConvertionMap = map[int]convertionFunction{
	1: func(rawValue float32) (convertedValue float32) { return float32(rawValue) / 10 },
	2: func(rawValue float32) (convertedValue float32) { return float32(rawValue) / 10 },
	3: func(rawValue float32) (convertedValue float32) { return float32(rawValue) / 10 },
	4: func(rawValue float32) (convertedValue float32) { return float32(rawValue) / 10 },
}
var SensorNameMap = map[int]string{
	1: "Исходящий",
	2: "Исходящий",
	3: "Исходящий",
	4: "Исходящий",
}

func (t AggregationType) ToString() string {
	switch t {
	case AGGR_SECOND:
		return "second"
	case AGGR_MINUTE:
		return "minute"
	case AGGR_5MINUTES:
		return "5minutes"
	case AGGR_HOUR:
		return "hour"
	case AGGR_DAY:
		return "day"
	case AGGR_WEEK:
		return "week"
	}
	panic(fmt.Errorf("This shouldn't happened"))
	return ""
}

func (t *AggregationType) FromString(s string) {
	switch s {
	case "second":
		*t = AGGR_SECOND
	case "minute":
		*t = AGGR_MINUTE
	case "5minutes":
		*t = AGGR_5MINUTES
	case "hour":
		*t = AGGR_HOUR
	case "day":
		*t = AGGR_DAY
	case "week":
		*t = AGGR_WEEK
	default:
		panic(fmt.Errorf("This shouldn't happened"))
	}

	return
}

func (t *AggregationType) Scan(value interface{}) error {
	t.FromString(string(value.([]uint8)))
	return nil
}
func (t AggregationType) Value() (driver.Value, error) { return t.ToString(), nil }

func (r *historyRecord) ConvertValue() error {
	f, ok := valueConvertionMap[r.SensorId]
	if !ok {
		return ErrInvalidSensorId
	}

	r.ConvertedValue = f(r.RawValue)

	return nil
}

func (r *historyRecord) FixDate() {
	var divider int64
	switch r.AggregationType {
	case AGGR_SECOND:
		divider = 1
	case AGGR_MINUTE:
		divider = 60
	case AGGR_5MINUTES:
		divider = 60 * 5
	case AGGR_HOUR:
		divider = 3600
	case AGGR_DAY:
		divider = 3600 * 24
	case AGGR_WEEK:
		divider = 3600 * 24 * 7
	default:
		panic(fmt.Errorf("This shouldn't happened"))
	}

	unixTS := r.Date.Unix()
	unixTS /= divider
	unixTS *= divider
	//unixTS += (divider / 2)

	r.Date = MyTime(time.Unix(unixTS, 0))

	return
}

func (r *historyRecord) Merge(another historyRecord) {
	r.ConvertedValue = (r.ConvertedValue*float32(r.Counter) + another.ConvertedValue) / (float32(r.Counter) + 1)
	r.RawValue = (r.RawValue*float32(r.Counter) + another.RawValue) / (float32(r.Counter) + 1)
	r.Counter++
}
