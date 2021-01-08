package structure

import (
	"errors"
	"reflect"
	"strings"
)

func GetTypeName(structure interface{}) string {
	return reflect.TypeOf(structure).Name()
}

func GetFieldNames(structure interface{}) (names []string) {
	var t reflect.Type
	switch reflect.TypeOf(structure).Kind() {
	case reflect.Struct:
		t = reflect.TypeOf(structure)
	case reflect.Ptr:
		t = reflect.TypeOf(structure).Elem()
	default:
		return names
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		names = append(names, f.Name)
	}
	return names
}

func SetFieldValue(structure interface{}, field string, value interface{}) error {
	var i reflect.Value
	switch reflect.TypeOf(structure).Kind() {
	case reflect.Struct:
		i = reflect.ValueOf(structure)
	case reflect.Ptr:
		i = reflect.ValueOf(structure).Elem()
	default:
		return errors.New("input structure error")
	}
	f := i.FieldByName(field)
	if !f.IsValid() {
		return errors.New("invalid field")
	}
	if !f.CanSet() {
		return errors.New("cannot set structure field")
	}
	t := f.Type()
	v := reflect.ValueOf(value)
	if t != v.Type() {
		return errors.New("structure field type does not match value type")
	}
	f.Set(v)
	return nil
}

func StructToMap(structure interface{}) map[string]interface{} {
	var t reflect.Type
	var v reflect.Value
	m := make(map[string]interface{})
	switch reflect.TypeOf(structure).Kind() {
	case reflect.Struct:
		t = reflect.TypeOf(structure)
		v = reflect.ValueOf(structure)
	case reflect.Ptr:
		t = reflect.TypeOf(structure).Elem()
		v = reflect.ValueOf(structure).Elem()
	default:
		return m
	}
	for i := 0; i < t.NumField(); i++ {
		fn := t.Field(i).Name
		fv := v.Field(i).Interface()
		if reflect.TypeOf(fv).Kind() == reflect.Slice {
			a := make([]interface{}, reflect.ValueOf(fv).Len())
			for j := 0; j < reflect.ValueOf(fv).Len(); j++ {
				a[j] = StructToMap(reflect.ValueOf(fv).Index(j).Interface())
			}
			m[fn] = a
		} else if reflect.TypeOf(fv).Kind() == reflect.Struct || reflect.TypeOf(fv).Kind() == reflect.Ptr {
			m[fn] = StructToMap(fv)
		} else {
			m[fn] = fv
		}
	}
	return m
}

func StructToStruct(source interface{}, destination interface{}) error {
	sm := StructToMap(source)
	var t reflect.Type
	var v reflect.Value
	switch reflect.TypeOf(destination).Kind() {
	case reflect.Struct:
		t = reflect.TypeOf(destination)
		v = reflect.ValueOf(destination)
	case reflect.Ptr:
		t = reflect.TypeOf(destination).Elem()
		v = reflect.ValueOf(destination).Elem()
	default:
		return errors.New("input structure error")
	}
	for i := 0; i < t.NumField(); i++ {
		fn := t.Field(i).Name
		ft := t.Field(i).Tag
		fv := v.Field(i).Interface()
		if reflect.TypeOf(fv).Kind() == reflect.Slice || reflect.TypeOf(fv).Kind() == reflect.Array {
			tv := strings.Split(ft.Get("struct"), ",")[0]
			ks := strings.Split(tv, ".")
			var sv interface{}
			var ok bool
			for i, k := range ks {
				if i == 0 {
					if sv, ok = sm[k]; !ok {
						break
					}
				} else {
					if sv, ok = sv.(map[string]interface{})[k]; !ok {
						break
					}
				}
			}
			if ok {
				afv := reflect.MakeSlice(reflect.TypeOf(fv), reflect.ValueOf(sv).Len(), reflect.ValueOf(sv).Len())
				err := SetFieldValue(destination, fn, afv.Interface())
				if err != nil {
					return err
				}
			}
		} else if reflect.TypeOf(fv).Kind() == reflect.Struct || reflect.TypeOf(fv).Kind() == reflect.Ptr {
			fv = reflect.New(reflect.TypeOf(fv)).Interface()
			err := structToStruct(sm, fv)
			if err != nil {
				return err
			}
			err = SetFieldValue(destination, fn, reflect.ValueOf(fv).Elem().Interface())
			if err != nil {
				return err
			}
		} else {
			tv := strings.Split(ft.Get("struct"), ",")[0]
			ks := strings.Split(tv, ".")
			var sv interface{}
			var ok bool
			for i, k := range ks {
				if i == 0 {
					if sv, ok = sm[k]; !ok {
						break
					}
				} else {
					if sv, ok = sv.(map[string]interface{})[k]; !ok {
						break
					}
				}
			}
			if ok {
				SetFieldValue(destination, fn, sv)
			}
		}
	}
	return nil
}

func structToStruct(sm map[string]interface{}, d interface{}) error {
	var t reflect.Type
	var v reflect.Value
	switch reflect.TypeOf(d).Kind() {
	case reflect.Struct:
		t = reflect.TypeOf(d)
		v = reflect.ValueOf(d)
	case reflect.Ptr:
		t = reflect.TypeOf(d).Elem()
		v = reflect.ValueOf(d).Elem()
	default:
		return errors.New("input structure error")
	}
	for i := 0; i < t.NumField(); i++ {
		fn := t.Field(i).Name
		ft := t.Field(i).Tag
		fv := v.Field(i).Interface()
		if reflect.TypeOf(fv).Kind() == reflect.Struct || reflect.TypeOf(fv).Kind() == reflect.Ptr {
			fv = reflect.New(reflect.TypeOf(fv)).Interface()
			err := structToStruct(sm, fv)
			if err != nil {
				return err
			}
			err = SetFieldValue(d, fn, fv)
			if err != nil {
				return err
			}
		} else {
			tv := strings.Split(ft.Get("struct"), ",")[0]
			ks := strings.Split(tv, ".")
			var sv interface{}
			var ok bool
			for i, k := range ks {
				if i == 0 {
					if sv, ok = sm[k]; !ok {
						break
					}
				} else {
					if sv, ok = sv.(map[string]interface{})[k]; !ok {
						break
					}
				}
			}
			if ok {
				SetFieldValue(d, fn, sv)
			}
		}
	}
	return nil
}
