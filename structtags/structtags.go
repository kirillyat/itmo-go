//go:build !solution

package structtags

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// fieldMapCache stores a mapping from struct types to their processed field maps.
var fieldMapCache sync.Map

func Unpack(req *http.Request, ptr interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	t := reflect.TypeOf(ptr).Elem()
	fieldsMapInterface, ok := fieldMapCache.Load(t)
	if !ok {
		valueMap := make(map[string]reflect.StructField)
		vType := reflect.TypeOf(ptr).Elem()
		for i := 0; i < vType.NumField(); i++ {
			field := vType.Field(i)
			tag := field.Tag.Get("http")
			if tag == "" {
				tag = strings.ToLower(field.Name)
			}
			valueMap[tag] = field
		}
		fieldMapCache.Store(t, valueMap)
		fieldsMapInterface = valueMap
	}

	fieldsMap := fieldsMapInterface.(map[string]reflect.StructField)
	v := reflect.ValueOf(ptr).Elem()

	// Update struct field for each parameter in the request.
	for name, values := range req.Form {
		field, ok := fieldsMap[name]
		if !ok {
			continue // ignore unrecognized HTTP parameters
		}
		fieldValue := v.FieldByIndex(field.Index)
		for _, value := range values {
			if err := populate(fieldValue, value); err != nil {
				return fmt.Errorf("%s: %v", name, err)
			}
		}
	}
	return nil
}

func populate(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)
	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)
	case reflect.Slice:
		elem := reflect.New(v.Type().Elem()).Elem()
		if err := populate(elem, value); err != nil {
			return err
		}
		v.Set(reflect.Append(v, elem))
	default:
		return fmt.Errorf("unsupported kind %s", v.Type())
	}
	return nil
}
