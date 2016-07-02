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

	db    *reform.DB
	where [][]interface{}
	order []string
	limit int
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
	s: parse.StructInfo{Type: "historyRecord", SQLSchema: "", SQLName: "history_records", Fields: []parse.FieldInfo{{Name: "Id", Type: "int", Column: "id"}, {Name: "Date", Type: "MyTime", Column: "date"}, {Name: "AggregationType", Type: "AggregationType", Column: "aggregation_period"}, {Name: "SensorId", Type: "int", Column: "sensor_id"}, {Name: "RawValue", Type: "float32", Column: "raw_value"}, {Name: "ConvertedValue", Type: "float32", Column: "converted_value"}, {Name: "Counter", Type: "int", Column: "counter"}}, PKFieldIndex: 0},
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
	return &historyRecordScope{historyRecord: *s, db: defaultDB_historyRecord}
}

// Sets DB to do queries
func (s *historyRecord) DB(db *reform.DB) (scope *historyRecordScope) { return s.Scope().DB(db) }
func (s *historyRecordScope) DB(db *reform.DB) *historyRecordScope {
	s.db = db
	return s
}

// Sets default DB (to do not call the scope.DB() method every time)
func (s *historyRecord) SetDefaultDB(db *reform.DB) (err error) {
	defaultDB_historyRecord = db
	return nil
}

// Compiles SQL tail for defined limit scope
// TODO: should be compiled via dialects
func (s *historyRecordScope) getLimitTail() (tail string, args []interface{}, err error) {
	if s.limit <= 0 {
		return
	}

	tail = fmt.Sprintf("%v", s.limit)
	return
}

// Compiles SQL tail for defined order scope
// TODO: should be compiled via dialects
func (s *historyRecordScope) getOrderTail() (tail string, args []interface{}, err error) {
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
func (s *historyRecordScope) getWhereTailForFilter(filter HistoryRecordFilter) (tail string, whereTailArgs []interface{}, err error) {
	var whereTailStringParts []string

	sample := historyRecord(filter)

	v := reflect.ValueOf(sample)
	vT := v.Type()

	numField := v.NumField()

	placeholderCounter := 0
	for i := 0; i < numField; i++ {
		f := v.Field(i)
		fT := f.Type()

		if f.Interface() == reflect.Zero(fT).Interface() {
			continue
		}

		vs := vT.Field(i)
		rN := vs.Tag.Get("reform")

		placeholderCounter++
		whereTailStringParts = append(whereTailStringParts, rN+" = "+s.db.Dialect.Placeholder(placeholderCounter)) // TODO: escape field name
		whereTailArgs = append(whereTailArgs, f.Interface())
	}

	tail = strings.Join(whereTailStringParts, " AND ")

	return
}

// parseQuerierArgs considers different ways of defning the tail (using scope properties or/and in_args)
func (s *historyRecordScope) parseWhereTailComponent(in_args []interface{}) (tail string, args []interface{}, err error) {
	if len(in_args) > 0 {
		switch arg := in_args[0].(type) {
		case string:
			tail = arg
			args = in_args[1:]
			return
		case historyRecord:
			if len(in_args) > 1 {
				s = s.Where(in_args[1:]...)
			}
			tail, args, err = s.getWhereTailForFilter(HistoryRecordFilter(arg))
		case HistoryRecordFilter:
			if len(in_args) > 1 {
				s = s.Where(in_args[1:]...)
			}
			tail, args, err = s.getWhereTailForFilter(arg)
		default:
			err = fmt.Errorf("Invalid first element of \"in_args\" (%v). It should be a string or HistoryRecordFilter.", reflect.ValueOf(arg).Type().Name())
			return
		}
	}

	return
}

// Compiles SQL tail for defined filter
// TODO: should be compiled via dialects
func (s *historyRecordScope) getWhereTail() (tail string, whereTailArgs []interface{}, err error) {
	var whereTailStringParts []string

	if len(s.where) == 0 {
		return
	}

	for _, whereComponent := range s.where {
		var whereTailStringPart string
		var whereTailArgsPart []interface{}

		whereTailStringPart, whereTailArgsPart, err = s.parseWhereTailComponent(whereComponent)
		if err != nil {
			return
		}

		if len(whereTailStringPart) > 0 {
			whereTailStringParts = append(whereTailStringParts, whereTailStringPart)
		}
		whereTailArgs = append(whereTailArgs, whereTailArgsPart...)
	}

	if len(whereTailStringParts) == 0 {
		return
	}

	tail = "(" + strings.Join(whereTailStringParts, ") AND (") + ")"

	return
}

func (s *historyRecordScope) Where(in_args ...interface{}) *historyRecordScope {
	s.where = append(s.where, in_args)
	return s
}

// Compiles SQL tail for defined db/where/order/limit scope
// TODO: should be compiled via dialects
func (s *historyRecordScope) getTail() (tail string, args []interface{}, err error) {
	whereTailString, whereTailArgs, err := s.getWhereTail()

	if err != nil {
		return
	}
	orderTailString, orderTailArgs, err := s.getOrderTail()
	if err != nil {
		return
	}
	limitTailString, _, err := s.getLimitTail()
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

	if len(limitTailString) > 0 {
		limitTailString = " LIMIT " + limitTailString + " "
	}

	tail = whereTailString + orderTailString + limitTailString
	return

}

// Select is a wrapper for SelectRows() and NextRow(): it makes a query and collects the result into a slice
func (s *historyRecord) Select(args ...interface{}) (result []historyRecord, err error) {
	return s.Scope().Select(args...)
}
func (s *historyRecordScope) Select(args ...interface{}) (result []historyRecord, err error) {
	tail, args, err := s.Where(args...).getTail()
	if err != nil {
		return
	}

	rows, err := s.db.SelectRows(historyRecordTable, tail, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	for {
		err := s.db.NextRow(s, rows)
		if err != nil {
			break
		}
		result = append(result, (*s).historyRecord)
	}

	return
}

// "First" a method to select and return only one record.
func (s *historyRecord) First(args ...interface{}) (result historyRecord, err error) {
	return s.Scope().First(args...)
}
func (s *historyRecordScope) First(args ...interface{}) (result historyRecord, err error) {
	tail, args, err := s.Where(args...).getTail()
	if err != nil {
		return
	}

	err = s.db.SelectOneTo(&result, tail, args...)

	return
}

// Sets order. Arguments should be passed by pairs column-{ASC,DESC}. For example Order("id", "ASC", "value" "DESC")
func (s *historyRecord) Order(args ...interface{}) (scope *historyRecordScope) {
	return s.Scope().Order(args...)
}
func (s *historyRecordScope) Order(argsI ...interface{}) *historyRecordScope {
	switch len(argsI) {
	case 0:
	case 1:
		arg := argsI[0].(string)
		args0 := strings.Split(arg, ",")
		var args []string
		for _, arg0 := range args0 {
			args = append(args, strings.Split(arg0, ":")...)
		}
		s.order = args
	default:
		var args []string
		for _, argI := range argsI {
			args = append(args, argI.(string))
		}
		s.order = args
	}

	return s
}

// Sets limit.
func (s *historyRecord) Limit(limit int) (scope *historyRecordScope) { return s.Scope().Limit(limit) }
func (s *historyRecordScope) Limit(limit int) *historyRecordScope {
	s.limit = limit
	return s
}

// "Reload" reloads record using Primary Key
func (s *HistoryRecordFilter) Reload(db *reform.DB) error { return (*historyRecord)(s).Reload(db) }
func (s *historyRecord) Reload(db *reform.DB) (err error) {
	return db.FindByPrimaryKeyTo(s, s.PKValue())
}

// Create and Insert inserts new record to DB
func (s *historyRecord) Create() (err error) { return s.Scope().Create() }
func (s *historyRecordScope) Create() (err error) {
	return s.db.Insert(s)
}
func (s *historyRecord) Insert() (err error) { return s.Scope().Insert() }
func (s *historyRecordScope) Insert() (err error) {
	return s.db.Insert(s)
}

// Save inserts new record to DB is PK is zero and updates existing record if PK is not zero
func (s *historyRecord) Save() (err error) { return s.Scope().Save() }
func (s *historyRecordScope) Save() (err error) {
	return s.db.Save(s)
}

// Update updates existing record in DB
func (s *historyRecord) Update() (err error) { return s.Scope().Update() }
func (s *historyRecordScope) Update() (err error) {
	return s.db.Update(s)
}

// Delete deletes existing record in DB
func (s *historyRecord) Delete() (err error) { return s.Scope().Delete() }
func (s *historyRecordScope) Delete() (err error) {
	return s.db.Delete(s)
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
func (s *HistoryRecordFilter) SetPK(pk interface{}) { (*historyRecord)(s).SetPK(pk) }
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
	HistoryRecord           = historyRecord{} // Should be read only
	defaultDB_historyRecord *reform.DB
)

func init() {
	parse.AssertUpToDate(&historyRecordTable.s, new(historyRecord))
}
