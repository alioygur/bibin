package mysql

import (
	"errors"
	"reflect"
)

type (
	sqlStruct struct {
		tagName string
		columns []string
		values  []reflect.Value
	}
)

func newSQLStruct(vv interface{}) (*sqlStruct, error) {
	var s sqlStruct
	s.tagName = "db"

	v := reflect.ValueOf(vv)

	if v.Kind() != reflect.Ptr {
		return nil, errors.New("must pass a pointer")
	}

	v = v.Elem()

	if v.Kind() != reflect.Struct {
		return nil, errors.New("must pass pointer of struct")
	}

	s.setMap(v)

	if len(s.columns) == 0 {
		return nil, errors.New("did you forgot struct field tag? there are no available fields")
	}

	return &s, nil
}

func (s *sqlStruct) Columns(expects ...string) []string {
	var columns []string
	for _, c := range s.columns {
		if inStringSlice(c, expects) {
			continue
		}
		columns = append(columns, "`"+c+"`")
	}
	return columns
}

func (s *sqlStruct) Values(expects ...string) []interface{} {
	var values []interface{}
	for i, c := range s.columns {
		if inStringSlice(c, expects) {
			continue
		}
		values = append(values, s.values[i].Interface())
	}
	return values
}

func (s *sqlStruct) Ptrs() []interface{} {
	var ptrs []interface{}
	for _, v := range s.values {
		ptrs = append(ptrs, v.Addr().Interface())
	}
	return ptrs
}

func (s *sqlStruct) setMap(v reflect.Value) {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// is it embeded field?
		if t.Field(i).Anonymous && v.Field(i).Kind() == reflect.Struct {
			s.setMap(v.Field(i))
		}
		// if tag "db" not exists, then skip this field
		if t.Field(i).Tag.Get(s.tagName) == "" {
			continue
		}

		s.columns = append(s.columns, t.Field(i).Tag.Get(s.tagName))
		s.values = append(s.values, v.Field(i))
	}
}

func inStringSlice(needle string, stack []string) bool {
	for _, s := range stack {
		if needle == s {
			return true
		}
	}
	return false
}
