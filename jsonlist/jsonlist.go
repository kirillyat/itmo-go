//go:build !solution

package jsonlist

import (
	"encoding/json"
	"io"
	"reflect"
)

func Marshal(w io.Writer, slice interface{}) error {
	sliceVal := reflect.ValueOf(slice)
	if sliceVal.Kind() != reflect.Slice {
		return &json.UnsupportedTypeError{Type: reflect.TypeOf(slice)}
	}

	length := sliceVal.Len()
	for i := 0; i < length; i++ {
		elem := sliceVal.Index(i).Interface()
		bytes, err := json.Marshal(elem)
		if err != nil {
			return err
		}
		if i > 0 {
			if _, err := w.Write([]byte(" ")); err != nil {
				return err
			}
		}
		if _, err := w.Write(bytes); err != nil {
			return err
		}
	}
	return nil
}

func Unmarshal(r io.Reader, slicePointer interface{}) error {
	slicePtrVal := reflect.ValueOf(slicePointer)
	if slicePtrVal.Kind() != reflect.Ptr || slicePtrVal.Elem().Kind() != reflect.Slice {
		return &json.UnsupportedTypeError{Type: reflect.TypeOf(slicePointer)}
	}

	sliceVal := slicePtrVal.Elem()
	elemType := sliceVal.Type().Elem()
	decoder := json.NewDecoder(r)
	for {
		elemPtr := reflect.New(elemType).Interface()
		if err := decoder.Decode(elemPtr); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		sliceVal.Set(reflect.Append(sliceVal, reflect.ValueOf(elemPtr).Elem()))
	}

	return nil
}
