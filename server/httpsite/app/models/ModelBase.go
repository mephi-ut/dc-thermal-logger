package models

import (
	"fmt"
	"github.com/revel/revel"
	"reflect"
	"devel.mephi.ru/dc-thermal-logger/server/httpsite/app"
	"gopkg.in/reform.v1"
)

type ModelBase struct {
	this interface{}
}

var modelTableRegistry = map[string]reform.View{}

func (obj *ModelBase) Init(parent interface{}) error {
	obj.this = parent
	return nil
}

func (obj *ModelBase) AfterFind() error {
	err := fmt.Errorf("This shouldn't happened. Seems AfterFind() is not defined for some model (an example of required implementation is \"func (obj *…) AfterFind() error {return obj.Init(obj)}\").")
	revel.ERROR.Printf("%v", err.Error())
	return err
}

func (model ModelBase) Select(tail string) (interface{}, error) {
	modelType  := reflect.Indirect(reflect.ValueOf(model.this)).Type()
	modelName  := modelType.Name()
	sliceValue := reflect.MakeSlice(reflect.SliceOf(modelType), 0, 0)

	modelTable,ok := modelTableRegistry[modelName]
	if !ok {
		err := fmt.Errorf("Model \"%s\" is not registered", modelName)
		revel.ERROR.Printf("%s", err.Error())
		return sliceValue.Interface(), err
	}

	rows, err := app.DB.SelectRows(modelTable, tail)
	if err != nil {
		revel.ERROR.Printf("Cannot fetch log records: %s", err.Error())
		return sliceValue.Interface(), err
	}
	defer rows.Close()

	rowPtrValue := reflect.New(modelType)
	rowPtrValue.MethodByName("Init").Call([]reflect.Value{ rowPtrValue })
	for {
		err := app.DB.NextRow(rowPtrValue.Interface().(reform.Struct), rows)
		if err != nil {
			revel.TRACE.Printf("NextRow() ended -> %v", err.Error())
			break
		}
		sliceValue = reflect.Append(sliceValue, rowPtrValue.Elem())
	}

	return sliceValue.Interface(), nil
}

func modelRegister(model interface{}, table reform.View) {
	v := reflect.ValueOf(model)
	t := v.Type()

	modelTableRegistry[t.Name()] = table
}

