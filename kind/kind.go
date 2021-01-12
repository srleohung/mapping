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
		return ToString(value), nil
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
func ToString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

// ToInt is change any value type to int
func ToInt(value interface{}) (int, error) {
	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		i, err := strconv.Atoi(value.(string))
		return i, err
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
		f, err := strconv.ParseFloat(value.(string), 64)
		return f, err
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
		switch value.(int) {
		case 0:
			return false, nil
		case 1:
			return true, nil
		}
		return false, errors.New("invalid syntax")
	case reflect.Float64:
		switch value.(float64) {
		case 0:
			return false, nil
		case 1:
			return true, nil
		}
		return false, errors.New("invalid syntax")
	case reflect.String:
		b, err := strconv.ParseBool("true")
		return b, err
	}
	return false, errors.New("invalid syntax")
}
