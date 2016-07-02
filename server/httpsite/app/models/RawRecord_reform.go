package models

// generated with gopkg.in/reform.v1

import (
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)

type rawRecordTableType struct {
	s parse.StructInfo
	z []interface{}
}

type rawRecordScope struct {
	rawRecord

	db    *reform.DB
	where [][]interface{}
	order []string
	limit int
}

type RawRecordFilter rawRecord

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
	return []string{"id", "date", "raw_sensor_id", "raw_channel_id", "raw_value"}
}

// NewStruct makes a new struct for that view or table.
func (v *rawRecordTableType) NewStruct() reform.Struct {
	return new(rawRecord)
}

// NewRecord makes a new record for that table.
func (v *rawRecordTableType) NewRecord() reform.Record {
	return new(rawRecord)
}

func (v *rawRecordTableType) NewScope() *rawRecordScope {
	return &rawRecordScope{}
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *rawRecordTableType) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// rawRecordTable represents raw_records view or table in SQL database.
var rawRecordTable = &rawRecordTableType{
	s: parse.StructInfo{Type: "rawRecord", SQLSchema: "", SQLName: "raw_records", Fields: []parse.FieldInfo{{Name: "Id", Type: "int", Column: "id"}, {Name: "Date", Type: "time.Time", Column: "date"}, {Name: "RawSensorId", Type: "int", Column: "raw_sensor_id"}, {Name: "RawChannelId", Type: "int", Column: "raw_channel_id"}, {Name: "RawValue", Type: "int", Column: "raw_value"}}, PKFieldIndex: 0},
	z: new(rawRecord).Values(),
}

// String returns a string representation of this struct or record.
func (s rawRecord) String() string {
	res := make([]string, 5)
	res[0] = "Id: " + reform.Inspect(s.Id, true)
	res[1] = "Date: " + reform.Inspect(s.Date, true)
	res[2] = "RawSensorId: " + reform.Inspect(s.RawSensorId, true)
	res[3] = "RawChannelId: " + reform.Inspect(s.RawChannelId, true)
	res[4] = "RawValue: " + reform.Inspect(s.RawValue, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *rawRecord) Values() []interface{} {
	return []interface{}{
		s.Id,
		s.Date,
		s.RawSensorId,
		s.RawChannelId,
		s.RawValue,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *rawRecord) Pointers() []interface{} {
	return []interface{}{
		&s.Id,
		&s.Date,
		&s.RawSensorId,
		&s.RawChannelId,
		&s.RawValue,
	}
}

// View returns View object for that struct.
func (s *rawRecord) View() reform.View {
	return rawRecordTable
}

// Generate a scope for object
func (s *rawRecord) Scope() *rawRecordScope {
	return &rawRecordScope{rawRecord: *s, db: defaultDB_rawRecord}
}

// Sets DB to do queries
func (s *rawRecord) DB(db *reform.DB) (scope *rawRecordScope) { return s.Scope().DB(db) }
func (s *rawRecordScope) DB(db *reform.DB) *rawRecordScope {
	s.db = db
	return s
}

// Sets default DB (to do not call the scope.DB() method every time)
func (s *rawRecord) SetDefaultDB(db *reform.DB) (err error) {
	defaultDB_rawRecord = db
	return nil
}

// Compiles SQL tail for defined limit scope
// TODO: should be compiled via dialects
func (s *rawRecordScope) getLimitTail() (tail string, args []interface{}, err error) {
	if s.limit <= 0 {
		return
	}

	tail = fmt.Sprintf("%v", s.limit)
	return
}

// Compiles SQL tail for defined order scope
// TODO: should be compiled via dialects
func (s *rawRecordScope) getOrderTail() (tail string, args []interface{}, err error) {
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
func (s *rawRecordScope) getWhereTailForFilter(filter RawRecordFilter) (tail string, whereTailArgs []interface{}, err error) {
	var whereTailStringParts []string

	sample := rawRecord(filter)

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
func (s *rawRecordScope) parseWhereTailComponent(in_args []interface{}) (tail string, args []interface{}, err error) {
	if len(in_args) > 0 {
		switch arg := in_args[0].(type) {
		case string:
			tail = arg
			args = in_args[1:]
			return
		case rawRecord:
			if len(in_args) > 1 {
				s = s.Where(in_args[1:]...)
			}
			tail, args, err = s.getWhereTailForFilter(RawRecordFilter(arg))
		case RawRecordFilter:
			if len(in_args) > 1 {
				s = s.Where(in_args[1:]...)
			}
			tail, args, err = s.getWhereTailForFilter(arg)
		default:
			err = fmt.Errorf("Invalid first element of \"in_args\" (%v). It should be a string or RawRecordFilter.", reflect.ValueOf(arg).Type().Name())
			return
		}
	}

	return
}

// Compiles SQL tail for defined filter
// TODO: should be compiled via dialects
func (s *rawRecordScope) getWhereTail() (tail string, whereTailArgs []interface{}, err error) {
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

func (s *rawRecordScope) Where(in_args ...interface{}) *rawRecordScope {
	s.where = append(s.where, in_args)
	return s
}

// Compiles SQL tail for defined db/where/order/limit scope
// TODO: should be compiled via dialects
func (s *rawRecordScope) getTail() (tail string, args []interface{}, err error) {
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
func (s *rawRecord) Select(args ...interface{}) (result []rawRecord, err error) {
	return s.Scope().Select(args...)
}
func (s *rawRecordScope) Select(args ...interface{}) (result []rawRecord, err error) {
	tail, args, err := s.Where(args...).getTail()
	if err != nil {
		return
	}

	rows, err := s.db.SelectRows(rawRecordTable, tail, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	for {
		err := s.db.NextRow(s, rows)
		if err != nil {
			break
		}
		result = append(result, (*s).rawRecord)
	}

	return
}

// "First" a method to select and return only one record.
func (s *rawRecord) First(args ...interface{}) (result rawRecord, err error) {
	return s.Scope().First(args...)
}
func (s *rawRecordScope) First(args ...interface{}) (result rawRecord, err error) {
	tail, args, err := s.Where(args...).getTail()
	if err != nil {
		return
	}

	err = s.db.SelectOneTo(&result, tail, args...)

	return
}

// Sets order. Arguments should be passed by pairs column-{ASC,DESC}. For example Order("id", "ASC", "value" "DESC")
func (s *rawRecord) Order(args ...interface{}) (scope *rawRecordScope) {
	return s.Scope().Order(args...)
}
func (s *rawRecordScope) Order(argsI ...interface{}) *rawRecordScope {
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
func (s *rawRecord) Limit(limit int) (scope *rawRecordScope) { return s.Scope().Limit(limit) }
func (s *rawRecordScope) Limit(limit int) *rawRecordScope {
	s.limit = limit
	return s
}

// "Reload" reloads record using Primary Key
func (s *RawRecordFilter) Reload(db *reform.DB) error { return (*rawRecord)(s).Reload(db) }
func (s *rawRecord) Reload(db *reform.DB) (err error) {
	return db.FindByPrimaryKeyTo(s, s.PKValue())
}

// Create and Insert inserts new record to DB
func (s *rawRecord) Create() (err error) { return s.Scope().Create() }
func (s *rawRecordScope) Create() (err error) {
	return s.db.Insert(s)
}
func (s *rawRecord) Insert() (err error) { return s.Scope().Insert() }
func (s *rawRecordScope) Insert() (err error) {
	return s.db.Insert(s)
}

// Save inserts new record to DB is PK is zero and updates existing record if PK is not zero
func (s *rawRecord) Save() (err error) { return s.Scope().Save() }
func (s *rawRecordScope) Save() (err error) {
	return s.db.Save(s)
}

// Update updates existing record in DB
func (s *rawRecord) Update() (err error) { return s.Scope().Update() }
func (s *rawRecordScope) Update() (err error) {
	return s.db.Update(s)
}

// Delete deletes existing record in DB
func (s *rawRecord) Delete() (err error) { return s.Scope().Delete() }
func (s *rawRecordScope) Delete() (err error) {
	return s.db.Delete(s)
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
func (s *RawRecordFilter) SetPK(pk interface{}) { (*rawRecord)(s).SetPK(pk) }
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
	RawRecord           = rawRecord{} // Should be read only
	defaultDB_rawRecord *reform.DB
)

func init() {
	parse.AssertUpToDate(&rawRecordTable.s, new(rawRecord))
}
