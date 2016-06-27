package models

// generated with gopkg.in/reform.v1

import (
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)

type historyRecordTableType struct {
	s parse.StructInfo
	z []interface{}
}

type historyRecordScope struct {
	historyRecord
	order []string
}

type HistoryRecordFilter historyRecord

// Schema returns a schema name in SQL database ("").
func (v *historyRecordTableType) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("history_records").
func (v *historyRecordTableType) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *historyRecordTableType) Columns() []string {
	return []string{"id", "date", "aggregation_period", "sensor_id", "raw_value", "converted_value", "counter"}
}

// NewStruct makes a new struct for that view or table.
func (v *historyRecordTableType) NewStruct() reform.Struct {
	return new(historyRecord)
}

// NewRecord makes a new record for that table.
func (v *historyRecordTableType) NewRecord() reform.Record {
	return new(historyRecord)
}

func (v *historyRecordTableType) NewScope() *historyRecordScope {
	return &historyRecordScope{}
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *historyRecordTableType) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// historyRecordTable represents history_records view or table in SQL database.
var historyRecordTable = &historyRecordTableType{
	s: parse.StructInfo{Type: "historyRecord", SQLSchema: "", SQLName: "history_records", Fields: []parse.FieldInfo{{Name: "Id", Type: "int", Column: "id"}, {Name: "Date", Type: "time.Time", Column: "date"}, {Name: "AggregationType", Type: "AggregationType", Column: "aggregation_period"}, {Name: "SensorId", Type: "int", Column: "sensor_id"}, {Name: "RawValue", Type: "float32", Column: "raw_value"}, {Name: "ConvertedValue", Type: "float32", Column: "converted_value"}, {Name: "Counter", Type: "int", Column: "counter"}}, PKFieldIndex: 0},
	z: new(historyRecord).Values(),
}

// String returns a string representation of this struct or record.
func (s historyRecord) String() string {
	res := make([]string, 7)
	res[0] = "Id: " + reform.Inspect(s.Id, true)
	res[1] = "Date: " + reform.Inspect(s.Date, true)
	res[2] = "AggregationType: " + reform.Inspect(s.AggregationType, true)
	res[3] = "SensorId: " + reform.Inspect(s.SensorId, true)
	res[4] = "RawValue: " + reform.Inspect(s.RawValue, true)
	res[5] = "ConvertedValue: " + reform.Inspect(s.ConvertedValue, true)
	res[6] = "Counter: " + reform.Inspect(s.Counter, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *historyRecord) Values() []interface{} {
	return []interface{}{
		s.Id,
		s.Date,
		s.AggregationType,
		s.SensorId,
		s.RawValue,
		s.ConvertedValue,
		s.Counter,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *historyRecord) Pointers() []interface{} {
	return []interface{}{
		&s.Id,
		&s.Date,
		&s.AggregationType,
		&s.SensorId,
		&s.RawValue,
		&s.ConvertedValue,
		&s.Counter,
	}
}

// View returns View object for that struct.
func (s *historyRecord) View() reform.View {
	return historyRecordTable
}

// Generate a scope for object
func (s *historyRecord) Scope() *historyRecordScope {
	return &historyRecordScope{historyRecord: *s}
}

// Compiles SQL tail for defined order scope
// TODO: should be compiled via dialects
func (s *historyRecordScope) getOrderTail(db *reform.DB) (tail string, args []interface{}, err error) {
	var fieldName string
	var orderStringParts []string

	for idx, orderStr := range s.order {
		switch idx % 2 {
		case 0:
			fieldName = orderStr
		case 1:
			orderDirection := orderStr

			orderStringParts = append(orderStringParts, fieldName+" "+orderDirection) // TODO: escape field name
		}
	}

	tail = strings.Join(orderStringParts, ", ")

	return
}

// Compiles SQL tail for defined filter
// TODO: should be compiled via dialects
func (s *historyRecordScope) getWhereTail(db *reform.DB, filter HistoryRecordFilter) (tail string, whereTailArgs []interface{}, err error) {
	var whereTailStringParts []string

	sample := historyRecord(filter)

	v := reflect.ValueOf(sample)
	vT := v.Type()

	numField := v.NumField()

	counter := 0
	for i := 0; i < numField; i++ {
		f := v.Field(i)
		fT := f.Type()

		if f.Interface() == reflect.Zero(fT).Interface() {
			continue
		}

		s := vT.Field(i)
		rN := s.Tag.Get("reform")

		counter++
		whereTailStringParts = append(whereTailStringParts, rN+" = "+db.Dialect.Placeholder(counter)) // TODO: escape field name
		whereTailArgs = append(whereTailArgs, f.Interface())
	}

	tail = strings.Join(whereTailStringParts, " AND ")

	return
}

// Compiles SQL tail for defined order scope and filter
// TODO: should be compiled via dialects
func (s *historyRecordScope) compileTailUsingFilter(db *reform.DB, filter HistoryRecordFilter) (tail string, args []interface{}, err error) {
	whereTailString, whereTailArgs, err := s.getWhereTail(db, filter)
	if err != nil {
		return
	}
	orderTailString, orderTailArgs, err := s.getOrderTail(db)
	if err != nil {
		return
	}

	args = append(whereTailArgs, orderTailArgs...)

	if len(whereTailString) > 0 {
		whereTailString = " WHERE " + whereTailString + " "
	}

	if len(orderTailString) > 0 {
		orderTailString = " ORDER BY " + orderTailString + " "
	}

	tail = whereTailString + orderTailString
	return

}

// parseQuerierArgs considers different ways of defning the tail (using scope properties or/and in_args)
func (s *historyRecordScope) parseQuerierArgs(db *reform.DB, in_args []interface{}) (tail string, args []interface{}, err error) {
	if len(in_args) > 0 {
		switch arg := in_args[0].(type) {
		case string:
			if len(s.order) > 0 {
				err = fmt.Errorf("This case is not implemented yet. You cannot use Order() and string tail argument in one request.")
				return
			}
			tail = arg
			args = in_args[1:]
		case historyRecord:
			if len(args) > 1 {
				err = fmt.Errorf("Too many arguments.")
				return
			}
			tail, args, err = s.compileTailUsingFilter(db, HistoryRecordFilter(arg))
		case HistoryRecordFilter:
			if len(args) > 1 {
				err = fmt.Errorf("Too many arguments.")
				return
			}
			tail, args, err = s.compileTailUsingFilter(db, arg)
		default:
			err = fmt.Errorf("Invalid first element of \"args\". It should be a string or HistoryRecordFilter.")
		}
	}

	return
}

// Select is a wrapper for SelectRows() and NextRow(): it makes a query and collects the result into a slice
func (s *historyRecord) Select(db *reform.DB, args ...interface{}) (result []historyRecord, err error) {
	return s.Scope().Select(db, args...)
}
func (s *historyRecordScope) Select(db *reform.DB, args ...interface{}) (result []historyRecord, err error) {
	tail, args, err := s.parseQuerierArgs(db, args)
	if err != nil {
		return
	}

	rows, err := db.SelectRows(historyRecordTable, tail, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	for {
		err := db.NextRow(s, rows)
		if err != nil {
			break
		}
		result = append(result, (*s).historyRecord)
	}

	return
}

// "First" a method to select and return only one record.
func (s *historyRecord) First(db *reform.DB, args ...interface{}) (result historyRecord, err error) {
	return s.Scope().First(db, args...)
}
func (s *historyRecordScope) First(db *reform.DB, args ...interface{}) (result historyRecord, err error) {
	tail, args, err := s.parseQuerierArgs(db, args)
	if err != nil {
		return
	}

	err = db.SelectOneTo(&result, tail, args...)

	return
}

// Create and Insert inserts new record to DB
func (s *historyRecord) Create(db *reform.DB) (err error) { return s.Scope().Create(db) }
func (s *historyRecordScope) Create(db *reform.DB) (err error) {
	return db.Insert(s)
}
func (s *historyRecord) Insert(db *reform.DB) (err error) { return s.Scope().Insert(db) }
func (s *historyRecordScope) Insert(db *reform.DB) (err error) {
	return db.Insert(s)
}

// Save inserts new record to DB is PK is zero and updates existing record if PK is not zero
func (s *historyRecord) Save(db *reform.DB) (err error) { return s.Scope().Save(db) }
func (s *historyRecordScope) Save(db *reform.DB) (err error) {
	return db.Save(s)
}

// Update updates existing record in DB
func (s *historyRecord) Update(db *reform.DB) (err error) { return s.Scope().Update(db) }
func (s *historyRecordScope) Update(db *reform.DB) (err error) {
	return db.Update(s)
}

// Delete deletes existing record in DB
func (s *historyRecord) Delete(db *reform.DB) (err error) { return s.Scope().Delete(db) }
func (s *historyRecordScope) Delete(db *reform.DB) (err error) {
	return db.Delete(s)
}

// Sets order. Arguments should be passed by pairs column-{ASC,DESC}. For example Order("id", "ASC", "value" "DESC")
func (s *historyRecord) Order(args ...string) (scope *historyRecordScope) {
	return s.Scope().Order(args...)
}
func (s *historyRecordScope) Order(args ...string) *historyRecordScope {
	s.order = args
	return s
}

// Table returns Table object for that record.
func (s *historyRecord) Table() reform.Table {
	return historyRecordTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *historyRecord) PKValue() interface{} {
	return s.Id
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *historyRecord) PKPointer() interface{} {
	return &s.Id
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *historyRecord) HasPK() bool {
	return s.Id != historyRecordTable.z[historyRecordTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *historyRecord) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.Id = int(i64)
	} else {
		s.Id = pk.(int)
	}
}

var (
	// check interfaces
	_ reform.View   = historyRecordTable
	_ reform.Struct = new(historyRecord)
	_ reform.Table  = historyRecordTable
	_ reform.Record = new(historyRecord)
	_ fmt.Stringer  = new(historyRecord)

	// querier
	HistoryRecord = historyRecord{} // Should be read only
)

func init() {
	parse.AssertUpToDate(&historyRecordTable.s, new(historyRecord))
}
