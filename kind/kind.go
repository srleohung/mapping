package kind

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// ToType is change any value type to target type
func ToType(value interface{}, typ reflect.Kind) (interface{}, error) {
	switch typ {
	case reflect.String:
		return ToString(value)
	case reflect.Int:
		return ToInt(value)
	case reflect.Float64:
		return ToFloat64(value)
	case reflect.Bool:
		return ToBool(value)
	}
	return nil, errors.New("invalid syntax")
}

// ToString is change any value type to string
func ToString(value interface{}) (string, error) {
	return fmt.Sprintf("%v", value), nil
}

// ToInt is change any value type to int
func ToInt(value interface{}) (int, error) {
	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		return strconv.Atoi(value.(string))
	case reflect.Float64:
		return int(value.(float64)), nil
	case reflect.Bool:
		if value.(bool) {
			return 1, nil
		}
		return 0, nil
	}
	return 0, errors.New("invalid syntax")
}

// ToFloat64 is change any value type to float64
func ToFloat64(value interface{}) (float64, error) {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Int:
		return float64(value.(int)), nil
	case reflect.String:
		return strconv.ParseFloat(value.(string), 64)
	case reflect.Bool:
		if value.(bool) {
			return 1, nil
		}
		return 0, nil
	}
	return 0, errors.New("invalid syntax")
}

// ToBool is change any value type to bool
func ToBool(value interface{}) (bool, error) {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Int:
		if value.(int) > 0 {
			return true, nil
		} else {
			return false, nil
		}
	case reflect.Float64:
		if value.(float64) > 0 {
			return true, nil
		} else {
			return false, nil
		}
	case reflect.String:
		b, err := strconv.ParseBool(value.(string))
		return b, err
	}
	return false, errors.New("invalid syntax")
}
