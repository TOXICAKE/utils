package utils

import (
	"log"
	"reflect"
)

var DebugMode = false

func DebugEcho(v ...interface{}) {
	if DebugMode {
		log.Println(v...)
	}
}

type N struct {
	Arg1 string
	Arg2 int
}

// CopyStructByName 复制结构体中字段名相同且类型相同的成员，第二个参数需要提供地址
func CopyStructByName(src interface{}, dest interface{}) {

	type Members map[string]reflect.Type

	srcT := reflect.TypeOf(src)
	srcV := reflect.ValueOf(src)
	destT := reflect.TypeOf(dest).Elem()
	destV := reflect.ValueOf(dest).Elem()

	srcMembers := Members{}
	for i := 0; i < srcT.NumField(); i++ {
		srcMembers[srcT.Field(i).Name] = srcV.Type()
	}

	for i := 0; i < destT.NumField(); i++ {
		name := destT.Field(i).Name
		if _, ok := srcMembers[name]; ok {
			if srcV.FieldByName(name).Type() == destV.Field(i).Type() {
				destV.Field(i).Set(srcV.FieldByName(name))
			}
		}
	}
}
