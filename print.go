/**
 * Copyright 2020 Innodev LLC. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package print

import (
	"fmt"
	"reflect"
	"strings"
)

const IndentStr = "    "

func Printf(fmtStr string, vals ...interface{}) {
	Print(fmt.Sprintf(fmtStr, vals...))
}

func Print(v interface{}) {
	fmt.Print(StringS(0, v))
}

func handleMap(depth int, val reflect.Value) string {
	out := "\n"
	iter := val.MapRange()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		out += StringKV(depth, k.Interface(), v.Interface())
	}
	return out
}

func handleSlice(depth int, val reflect.Value) string {
	out := ""
	for i := 0; i < val.Len(); i++ {
		out += StringS(depth, val.Index(i).Interface())
	}
	return out
}

func stringS(depth int, v interface{}) (out string) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += IndentStr
	}
	if strger, ok := v.(fmt.Stringer); ok {
		return indent + strger.String()
	}
	rv := reflect.ValueOf(v)
	t := rv.Type()
	if t.Kind() == reflect.Ptr && !rv.IsNil() {
		return StringS(depth, rv.Elem().Interface())
	}
	if t.Kind() == reflect.Map {
		return handleMap(depth, rv)
	}
	if t.Kind() == reflect.Slice {
		return handleSlice(depth, rv)
	}
	if t.Kind() != reflect.Struct {
		return fmt.Sprintf(indent+"%+v", v)
	}

	for i := 0; i < t.NumField(); i++ {
		if !rv.Field(i).CanInterface() {
			continue
		}
		field := t.Field(i)
		name := field.Name
		if field.Tag.Get("cli") != "" {
			name = field.Tag.Get("cli")
			name = strings.Split(name, ",")[0]
		}
		res := StringKV(depth, name, rv.Field(i).Interface())
		if strings.Contains(field.Tag.Get("cli"), "omitempty") && stringS(0, rv.Field(i).Interface()) == "" {
			continue
		}
		out += res
	}
	return
}

func stringKV(depth int, k interface{}, v interface{}) (out string) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += IndentStr
	}
	out = indent + fmt.Sprint(k) + ": "
	if strger, ok := v.(fmt.Stringer); ok {
		out += indent + strger.String()
		return
	}
	rv := reflect.ValueOf(v)
	t := rv.Type()
	if t.Kind() == reflect.Ptr && !rv.IsNil() {
		out += StringS(depth, rv.Elem().Interface())
		return
	}
	if t.Kind() == reflect.Slice {
		out += "\n" + handleSlice(depth+1, rv)
		return
	}
	if t.Kind() == reflect.Map {
		out += handleMap(depth+1, rv)
		return
	}
	if t.Kind() != reflect.Struct {
		out += fmt.Sprintf("%+v", v)
		return
	}
	if t.NumField() > 0 {
		out += "\n"
	} else {
		out = ""
	}

	for i := 0; i < t.NumField(); i++ {
		if !rv.Field(i).CanInterface() {
			continue
		}
		field := t.Field(i)
		name := field.Name
		if field.Tag.Get("cli") != "" {
			name = field.Tag.Get("cli")
			name = strings.Split(name, ",")[0]
		}
		res := StringKV(depth+1, name, rv.Field(i).Interface())
		if strings.Contains(field.Tag.Get("cli"), "omitempty") && stringS(0, rv.Field(i).Interface()) == "" {
			continue
		}
		out += res
	}
	return
}

func StringS(depth int, v interface{}) string {
	return stringS(depth, v) + "\n"
}

func StringKV(depth int, k interface{}, v interface{}) string {
	return stringKV(depth, k, v) + "\n"
}

func PrintKV(depth int, k interface{}, v interface{}) {
	fmt.Print(StringKV(depth, k, v))
}
