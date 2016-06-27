package models

// generated with gopkg.in/reform.v1

import (
	"fmt"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)

type rawRecordTableType struct {
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *rawRecordTableType) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("raw_records").
func (v *rawRecordTableType) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *rawRecordTableType) Columns() []string {
	return []string{"id", "date", "sensor_id", "channel_id", "raw_value"}
}

// NewStruct makes a new struct for that view or table.
func (v *rawRecordTableType) NewStruct() reform.Struct {
	return new(rawRecord)
}

// NewRecord makes a new record for that table.
func (v *rawRecordTableType) NewRecord() reform.Record {
	return new(rawRecord)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *rawRecordTableType) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// rawRecordTable represents raw_records view or table in SQL database.
var rawRecordTable = &rawRecordTableType{
	s: parse.StructInfo{Type: "rawRecord", SQLSchema: "", SQLName: "raw_records", Fields: []parse.FieldInfo{{Name: "Id", Type: "int", Column: "id"}, {Name: "Date", Type: "time.Time", Column: "date"}, {Name: "SensorId", Type: "int", Column: "sensor_id"}, {Name: "ChannelId", Type: "int", Column: "channel_id"}, {Name: "RawValue", Type: "int", Column: "raw_value"}}, PKFieldIndex: 0},
	z: new(rawRecord).Values(),
}

// String returns a string representation of this struct or record.
func (s rawRecord) String() string {
	res := make([]string, 5)
	res[0] = "Id: " + reform.Inspect(s.Id, true)
	res[1] = "Date: " + reform.Inspect(s.Date, true)
	res[2] = "SensorId: " + reform.Inspect(s.SensorId, true)
	res[3] = "ChannelId: " + reform.Inspect(s.ChannelId, true)
	res[4] = "RawValue: " + reform.Inspect(s.RawValue, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *rawRecord) Values() []interface{} {
	return []interface{}{
		s.Id,
		s.Date,
		s.SensorId,
		s.ChannelId,
		s.RawValue,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *rawRecord) Pointers() []interface{} {
	return []interface{}{
		&s.Id,
		&s.Date,
		&s.SensorId,
		&s.ChannelId,
		&s.RawValue,
	}
}

// View returns View object for that struct.
func (s *rawRecord) View() reform.View {
	return rawRecordTable
}

// Select is a wrapper for SelectRows() and NextRow(): it makes a query and collects the result into a slice
func (s *rawRecord) Select(db *reform.DB, args ...interface{}) (result []rawRecord, err error) {
	var tail string

	if len(args) > 0 {
		switch arg := args[0].(type) {
		case string:
			tail = arg
			args = args[1:]
		case rawRecord:
			err = fmt.Errorf("This case is not implemented yet.")
			return
		default:
			err = fmt.Errorf("Invalid first element of \"args\". It should be a string or rawRecord.")
			return
		}
	}

	rows, err := db.SelectRows(rawRecordTable, tail, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	for {
		err := db.NextRow(s, rows)
		if err != nil {
			break
		}
		result = append(result, *s)
	}

	return
}

// Table returns Table object for that record.
func (s *rawRecord) Table() reform.Table {
	return rawRecordTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *rawRecord) PKValue() interface{} {
	return s.Id
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *rawRecord) PKPointer() interface{} {
	return &s.Id
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *rawRecord) HasPK() bool {
	return s.Id != rawRecordTable.z[rawRecordTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *rawRecord) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.Id = int(i64)
	} else {
		s.Id = pk.(int)
	}
}

var (
	// check interfaces
	_ reform.View   = rawRecordTable
	_ reform.Struct = new(rawRecord)
	_ reform.Table  = rawRecordTable
	_ reform.Record = new(rawRecord)
	_ fmt.Stringer  = new(rawRecord)

	// querier
	RawRecord = rawRecord{} // Should be read only
)

func init() {
	parse.AssertUpToDate(&rawRecordTable.s, new(rawRecord))
}
