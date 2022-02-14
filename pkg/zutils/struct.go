package zutils

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// MapFillStruct use map fill struct, returns error when field unkown or type error
func MapFillStruct(data map[string]interface{}, obj interface{}) error {
	for k, v := range data {
		if v == nil {
			continue
		}
		err := setField(obj, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// MapFillStructMust use map fill struct, none error will returns
func MapFillStructMust(data map[string]interface{}, obj interface{}) {
	for k, v := range data {
		if v == nil {
			continue
		}
		setField(obj, k, v)
	}
}

// StructToMap convert struct to map
func StructToMap(v interface{}, data map[string]interface{}) {
	if v == nil || data == nil {
		return
	}
	structToMap(reflect.ValueOf(v), data)
}

func structToMap(v reflect.Value, data map[string]interface{}) {
	t := v.Type()
	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	// Only struct are supported
	if t.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < t.NumField(); i++ {
		if !t.Field(i).Anonymous {
			name := t.Field(i).Tag.Get("map")
			if name == "-" {
				continue
			}
			if name == "" {
				name = t.Field(i).Name
			}
			data[name] = v.Field(i).Interface()
		} else {
			structToMap(v.Field(i).Addr(), data)
		}
	}
}

// StructToStruct deep copy struct
func StructToStruct(src, dst interface{}) {
	if src == nil || dst == nil {
		return
	}
	structToStruct(reflect.ValueOf(src), reflect.ValueOf(dst))
}

func structToStruct(src, dst reflect.Value) {
	st := src.Type()
	dt := dst.Type()
	if st.Kind() == reflect.Ptr {
		src = src.Elem()
		st = st.Elem()
	}
	if dt.Kind() == reflect.Ptr {
		dst = dst.Elem()
		dt = dt.Elem()
	}
	// Only struct are supported
	if st.Kind() != reflect.Struct || dt.Kind() != reflect.Struct {
		return
	}
	var field reflect.Value
	for i := 0; i < st.NumField(); i++ {
		if !st.Field(i).Anonymous {
			field = dst.FieldByName(st.Field(i).Name)
			if field.IsValid() && field.CanSet() {
				field.Set(src.Field(i))
			}
		} else {
			structToStruct(src.Field(i).Addr(), dst)
		}
	}
}

func setField(obj interface{}, name string, value interface{}) error {
	rv := reflect.ValueOf(obj)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return fmt.Errorf("cannot set %s on nil", name)
	}

	field := findField(rv, name)
	if !field.IsValid() {
		return fmt.Errorf("no such field: %s in obj", name)
	}

	if !field.CanSet() {
		return fmt.Errorf("cannot set %s field value", name)
	}

	fieldType := field.Type() // struct file type
	val := reflect.ValueOf(value)

	var err error
	if fieldType != val.Type() {
		val, err = typeConversion(fmt.Sprintf("%v", value), field.Type().Name()) // convert type
		if err != nil {
			return err
		}
	}

	field.Set(val)
	return nil
}

func findField(rv reflect.Value, name string) reflect.Value {
	rp := rv.Type()
	for i := 0; i < rp.NumField(); i++ {
		if rp.Field(i).Anonymous {
			return findField(rv.Field(i), name)
		}
		// use struct tag to named a field
		n := rp.Field(i).Tag.Get("map")
		if n == "-" {
			continue
		}
		if n == name {
			return rv.Field(i)
		}
		// directly named
		n = rp.Field(i).Name
		if n == name {
			return rv.Field(i)
		}
		// Camel named
		if CamelCase(name) == n {
			return rv.Field(i)
		}
	}

	return rv.FieldByName(name)
}

func typeConversion(value string, ntype string) (reflect.Value, error) {
	switch ntype {
	case "string":
		return reflect.ValueOf(value), nil
	case "time.Time", "Time", "time":
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	case "bool":
		b, err := strconv.ParseBool(value)
		return reflect.ValueOf(b), err
	case "int":
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	case "int8":
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	case "int16":
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int16(i)), err
	case "int32":
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int32(i)), err
	case "int64":
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	case "uint":
		i, err := strconv.ParseUint(value, 10, 64)
		return reflect.ValueOf(uint(i)), err
	case "uint8":
		i, err := strconv.ParseUint(value, 10, 64)
		return reflect.ValueOf(uint8(i)), err
	case "uint16":
		i, err := strconv.ParseUint(value, 10, 64)
		return reflect.ValueOf(uint16(i)), err
	case "uint32":
		i, err := strconv.ParseUint(value, 10, 64)
		return reflect.ValueOf(uint32(i)), err
	case "uint64":
		i, err := strconv.ParseUint(value, 10, 64)
		return reflect.ValueOf(i), err
	case "float32":
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	case "float64":
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}

	return reflect.ValueOf(value), errors.New("unsupport type: " + ntype + " " + value)
}
