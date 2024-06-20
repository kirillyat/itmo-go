//go:build !solution

package reversemap

import "reflect"

func ReverseMap(forward interface{}) interface{} {
	forwardValue := reflect.ValueOf(forward)

	if forwardValue.Kind() != reflect.Map {
		panic("input argument is not a map")
	}

	forwardType := forwardValue.Type()
	if forwardType.Key().Kind() != reflect.String && forwardType.Key().Kind() != reflect.Int {
		panic("keys in map must be of type string or int")
	}
	if forwardType.Elem().Kind() != reflect.String && forwardType.Elem().Kind() != reflect.Int {
		panic("values in map must be of type string or int")
	}

	reverseType := reflect.MapOf(forwardType.Elem(), forwardType.Key())
	reverse := reflect.MakeMap(reverseType)

	for _, key := range forwardValue.MapKeys() {
		reverse.SetMapIndex(forwardValue.MapIndex(key), key)
	}

	return reverse.Interface()
}
