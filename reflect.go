package main

import (
	"reflect"
)

type VariableType int

const (
	Null VariableType = iota
	Bool
	Number
	String
	Unknown
)

func variableTypeOf(v interface{}) VariableType {
	switch {
	case isNull(v):
		return Null
	case isBool(v):
		return Bool
	case isNumber(v):
		return Number
	case isString(v):
		return String
	default:
		return Unknown
	}
}

func isSlice(v interface{}) bool {
	switch reflect.ValueOf(v).Kind() {
	case reflect.Slice:
		return true
	}
	return false
}

func isMap(v interface{}) bool {
	switch reflect.ValueOf(v).Kind() {
	case reflect.Map:
		return true
	}
	return false
}

func isNull(v interface{}) bool {
	return v == nil
}

func isBool(v interface{}) bool {
	return reflect.ValueOf(v).Kind() == reflect.Bool
}

func isNumber(v interface{}) bool {
	switch reflect.ValueOf(v).Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64:
		return true
	}
	return false
}

func isString(v interface{}) bool {
	return reflect.ValueOf(v).Kind() == reflect.String
}

func isScalar(v interface{}) bool {
	return isNull(v) || isBool(v) || isNumber(v) || isString(v)
}
