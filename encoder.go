package encoder

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

// Marshal will encode all the fields of a struct into url.Values.
func Marshal(s interface{}) url.Values {
	values := url.Values{}

	v := reflect.ValueOf(s)
	addValue(values, "", v)

	return values
}

func addValue(values url.Values, key string, value reflect.Value) {
	switch value.Kind() {
	default:
		values.Add(key, value.String())
	case reflect.Struct:
		t := reflect.TypeOf(value.Interface())
		for i := 0; i < value.NumField(); i++ {
			fv := value.Field(i)
			sf := t.Field(i)

			urlKey := sf.Tag.Get("url")
			if urlKey == "" {
				urlKey = sf.Name
			}
			if key != "" {
				urlKey = fmt.Sprintf("%s.%s", key, urlKey)
			}

			addValue(values, urlKey, fv)
		}
	case reflect.Ptr:
		if value.IsNil() {
			return
		}
		addValue(values, key, value.Elem())
	case reflect.Int:
		values.Add(key, strconv.Itoa(int(value.Int())))
	case reflect.Map:
		for _, k := range value.MapKeys() {
			var urlKey string
			if key == "" {
				urlKey = k.String()
			} else {
				urlKey = fmt.Sprintf("%s.%s", key, k)
			}
			addValue(values, urlKey, value.MapIndex(k))
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			addValue(values, key, value.Index(i))
		}
	}
}
