package models

// generated with gopkg.in/reform.v1

import (
	"fmt"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)

type logRecordView struct {
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *logRecordView) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("log_records").
func (v *logRecordView) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *logRecordView) Columns() []string {
	return []string{"date", "sensor_id", "channel_id", "value", "converted_value"}
}

// NewStruct makes a new struct for that view or table.
func (v *logRecordView) NewStruct() reform.Struct {
	return new(LogRecord)
}

// LogRecordView represents log_records view or table in SQL database.
var LogRecordView = &logRecordView{
	s: parse.StructInfo{Type: "LogRecord", SQLSchema: "", SQLName: "log_records", Fields: []parse.FieldInfo{{Name: "Date", Type: "time.Time", Column: "date"}, {Name: "SensorId", Type: "int", Column: "sensor_id"}, {Name: "ChannelId", Type: "int", Column: "channel_id"}, {Name: "Value", Type: "int", Column: "value"}, {Name: "ConvertedValue", Type: "int", Column: "converted_value"}}, PKFieldIndex: -1},
	z: new(LogRecord).Values(),
}

// String returns a string representation of this struct or record.
func (s LogRecord) String() string {
	res := make([]string, 5)
	res[0] = "Date: " + reform.Inspect(s.Date, true)
	res[1] = "SensorId: " + reform.Inspect(s.SensorId, true)
	res[2] = "ChannelId: " + reform.Inspect(s.ChannelId, true)
	res[3] = "Value: " + reform.Inspect(s.Value, true)
	res[4] = "ConvertedValue: " + reform.Inspect(s.ConvertedValue, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *LogRecord) Values() []interface{} {
	return []interface{}{
		s.Date,
		s.SensorId,
		s.ChannelId,
		s.Value,
		s.ConvertedValue,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *LogRecord) Pointers() []interface{} {
	return []interface{}{
		&s.Date,
		&s.SensorId,
		&s.ChannelId,
		&s.Value,
		&s.ConvertedValue,
	}
}

// View returns View object for that struct.
func (s *LogRecord) View() reform.View {
	return LogRecordView
}

// check interfaces
var (
	_ reform.View   = LogRecordView
	_ reform.Struct = new(LogRecord)
	_ fmt.Stringer  = new(LogRecord)
)

func init() {
	parse.AssertUpToDate(&LogRecordView.s, new(LogRecord))
}
