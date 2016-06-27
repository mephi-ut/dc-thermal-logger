package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const (
	AGGR_SECOND AggregationType = iota
	AGGR_MINUTE
	AGGR_HOUR
	AGGR_DAY
	AGGR_WEEK
)

//go:generate reform

type AggregationType int

//reform:history_records
type historyRecord struct {
	ModelBase

	Id              int             `reform:"id,pk"`
	Date            time.Time       `reform:"date"`
	AggregationType AggregationType `reform:"aggregation_period"`
	SensorId        int             `reform:"sensor_id"`
	RawValue        float32         `reform:"raw_value"`
	ConvertedValue  float32         `reform:"converted_value"`
	Counter         int             `reform:"counter"`
}

var aggregationTypes = []AggregationType{AGGR_SECOND, AGGR_MINUTE, AGGR_HOUR, AGGR_DAY, AGGR_WEEK}
var sensorIdMap = map[int]map[int]int{
	1: {0: 1},
}

func (t AggregationType) ToString() string {
	switch t {
	case AGGR_SECOND:
		return "second"
	case AGGR_MINUTE:
		return "minute"
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

func (t *AggregationType) Scan(value interface{}) error { t.FromString(string(value.([]uint8))); return nil }
func (t AggregationType) Value() (driver.Value, error)  { return t.ToString(), nil }

func (r *historyRecord) ConvertValue() {
	r.ConvertedValue = float32(r.RawValue)
}

func (r *historyRecord) FixDate() {
	var divider int64
	switch r.AggregationType {
	case AGGR_SECOND:
		divider = 1
	case AGGR_MINUTE:
		divider = 60
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

	r.Date = time.Unix(unixTS, 0)

	return
}

func (r *historyRecord) Merge(another historyRecord) {
	r.ConvertedValue = (r.ConvertedValue*float32(r.Counter) + another.ConvertedValue) / (float32(r.Counter) + 1)
	r.RawValue = (r.RawValue*float32(r.Counter) + another.RawValue) / (float32(r.Counter) + 1)
	r.Counter++
}
