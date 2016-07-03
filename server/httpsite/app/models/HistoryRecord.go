package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"math"
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
	return []byte("\"" + t.Format("2006-01-02 15:04:05") + "\""), nil
}

var aggregationTypes = []AggregationType{AGGR_SECOND, AGGR_MINUTE, AGGR_5MINUTES, AGGR_HOUR, AGGR_DAY, AGGR_WEEK}
var sensorIdMap = map[int]map[int]int{
	1:  {0: 64, 1: 65, 2: 66, 3: 67, 4: 68, 5: 69, 6: 70, 7: 71},
	2:  {0: 72, 1: 73, 2: 74, 3: 75, 4: 76, 5: 77, 6: 78, 7: 79},
	32: {0: 1, 1: 2, 2: 3, 3: 4},
}

type convertionFunction func(rawValue float32) (convertedValue float32)

const (
	THERMO0_R0      float64 = 4.7
	THERMO0_B       float64 = 6119
	THERMO0_T0      float64 = 298.15
	THERMO0_R_corr  float64 = 0.452
	THERMO0_R_const float64 = 1.8
)

func calculateThermo0_fitted(rawValue float32) float32 {
	v := float64(rawValue) / 4096

	r := THERMO0_R_const * (1/v - 1)

	t := THERMO0_B / (math.Log((r-THERMO0_R_corr)/THERMO0_R0) + THERMO0_B/THERMO0_T0)

	return float32(t)
}

/*func calculateThermo0_graduated(rawValue float32) (float32) {
	return 0
}*/
func calculateThermo0(rawValue float32) float32 {
	return /*(*/ calculateThermo0_fitted(rawValue) /* + calculateThermo0_graduated(rawValue)) / 2*/
}

func calculateThermo1(rawValue float32) (convertedValue float32) {
	return 0
}

var valueConvertionMap = map[int]convertionFunction{
	1: func(rawValue float32) (convertedValue float32) { return float32(rawValue)/10 + 273.15 },
	2: func(rawValue float32) (convertedValue float32) { return float32(rawValue)/10 + 273.15 },
	3: func(rawValue float32) (convertedValue float32) { return float32(rawValue)/10 + 273.15 },
	4: func(rawValue float32) (convertedValue float32) { return float32(rawValue)/10 + 273.15 },

	64: func(rawValue float32) (convertedValue float32) { return calculateThermo0(rawValue) },
	65: func(rawValue float32) (convertedValue float32) { return calculateThermo0(rawValue) },
	66: func(rawValue float32) (convertedValue float32) { return calculateThermo0(rawValue) },
	67: func(rawValue float32) (convertedValue float32) { return calculateThermo0(rawValue) },
	68: func(rawValue float32) (convertedValue float32) { return calculateThermo0(rawValue) },
	69: func(rawValue float32) (convertedValue float32) { return calculateThermo0(rawValue) },
	70: func(rawValue float32) (convertedValue float32) { return calculateThermo1(rawValue) },

	71: func(rawValue float32) (convertedValue float32) { return calculateThermo0(rawValue) },
	72: func(rawValue float32) (convertedValue float32) { return calculateThermo0(rawValue) },
	73: func(rawValue float32) (convertedValue float32) { return calculateThermo0(rawValue) },
	74: func(rawValue float32) (convertedValue float32) { return calculateThermo0(rawValue) },
	75: func(rawValue float32) (convertedValue float32) { return calculateThermo0(rawValue) },
	76: func(rawValue float32) (convertedValue float32) { return calculateThermo0(rawValue) },
	77: func(rawValue float32) (convertedValue float32) { return calculateThermo0(rawValue) },
	78: func(rawValue float32) (convertedValue float32) { return calculateThermo0(rawValue) },
	79: func(rawValue float32) (convertedValue float32) { return calculateThermo1(rawValue) },
}
var SensorNameMap = map[int]string{
	1: "Исходящий",
	2: "Исходящий",
	3: "Исходящий",
	4: "Исходящий",

	64: "Зад",
	65: "Перед",
	66: "Пол",
	67: "Зад",
	68: "Перед",
	69: "Пол",
	70: "Не подключен",
	71: "Потолок",

	72: "Зад",
	73: "Перед",
	74: "Пол",
	75: "Зад",
	76: "Перед",
	77: "Пол",
	78: "Не подключен",
	79: "Потолок",
}
var SensorFullNameMap = map[int]string{
	1: "KD 1.1",
	2: "KD 3.1",
	3: "KD 2.1",

	64: "Стойка A05 - зад",
	65: "Стойка A05 - перед",
	66: "Стойка A05 - пол",
	67: "Стойка B05 - зад",
	68: "Стойка B05 - перед",
	69: "Стойка B05 - пол",
	70: "Стойка B05 - не подключен",
	71: "Стойка B05 - потолок",

	72: "Стойка - зад",
	73: "Стойка - перед",
	74: "Стойка - пол",
	75: "Стойка - зад",
	76: "Стойка - перед",
	77: "Стойка - пол",
	78: "Стойка - не подключен",
	79: "Стойка - потолок",
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
