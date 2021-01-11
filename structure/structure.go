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

func SearchFieldName(structure interface{}, key string) (name string) {
	var t reflect.Type
	switch reflect.TypeOf(structure).Kind() {
	case reflect.Struct:
		t = reflect.TypeOf(structure)
	case reflect.Ptr:
		t = reflect.TypeOf(structure).Elem()
	default:
		return ""
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		v := strings.Split(f.Tag.Get("struct"), ",")[0]
		if v == key {
			return f.Name
		}
	}
	return ""
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
		if strings.ToUpper(string(fn[0])) != string(fn[0]) {
			continue
		}
		fv := v.Field(i).Interface()
		if reflect.TypeOf(fv).Kind() == reflect.Slice {
			a := make([]interface{}, reflect.ValueOf(fv).Len())
			for j := 0; j < reflect.ValueOf(fv).Len(); j++ {
				if stm, err := structToMap(reflect.ValueOf(fv).Index(j).Interface()); err == nil {
					a[j] = stm
				} else {
					a[j] = reflect.ValueOf(fv).Index(j).Interface()
				}
			}
			m[fn] = a
		} else if reflect.TypeOf(fv).Kind() == reflect.Struct || reflect.TypeOf(fv).Kind() == reflect.Ptr {
			if stm, err := structToMap(fv); err == nil {
				m[fn] = stm
			} else {
				m[fn] = fv
			}
		} else {
			m[fn] = fv
		}
	}
	return m
}

func structToMap(s interface{}) (map[string]interface{}, error) {
	var t reflect.Type
	var v reflect.Value
	m := make(map[string]interface{})
	switch reflect.TypeOf(s).Kind() {
	case reflect.Struct:
		t = reflect.TypeOf(s)
		v = reflect.ValueOf(s)
	case reflect.Ptr:
		t = reflect.TypeOf(s).Elem()
		v = reflect.ValueOf(s).Elem()
	default:
		return m, errors.New("input structure error")
	}
	for i := 0; i < t.NumField(); i++ {
		fn := t.Field(i).Name
		if strings.ToUpper(string(fn[0])) != string(fn[0]) {
			return m, errors.New("cannot return value obtained from unexported field or method")
		}
		fv := v.Field(i).Interface()
		if reflect.TypeOf(fv).Kind() == reflect.Slice {
			a := make([]interface{}, reflect.ValueOf(fv).Len())
			for j := 0; j < reflect.ValueOf(fv).Len(); j++ {
				if stm, err := structToMap(reflect.ValueOf(fv).Index(j).Interface()); err == nil {
					a[j] = stm
				} else {
					a[j] = reflect.ValueOf(fv).Index(j).Interface()
				}
			}
			m[fn] = a
		} else if reflect.TypeOf(fv).Kind() == reflect.Struct || reflect.TypeOf(fv).Kind() == reflect.Ptr {
			if stm, err := structToMap(fv); err == nil {
				m[fn] = stm
			} else {
				m[fn] = fv
			}
		} else {
			m[fn] = fv
		}
	}
	return m, nil
}

func StructToStruct(source interface{}, destination interface{}) error {
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
	sm := StructToMap(source)
	for i := 0; i < t.NumField(); i++ {
		fn := t.Field(i).Name
		ft := t.Field(i).Tag
		if strings.ToUpper(string(fn[0])) != string(fn[0]) {
			continue
		}
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
					if reflect.TypeOf(sv).Kind() == reflect.Slice || reflect.TypeOf(sv).Kind() == reflect.Array {
						break
					}
					if sv, ok = sv.(map[string]interface{})[k]; !ok {
						break
					}
				}
			}
			if ok {
				if reflect.TypeOf(sv).Kind() != reflect.Slice && reflect.TypeOf(sv).Kind() != reflect.Array {
					continue
				}
				afv := reflect.MakeSlice(reflect.TypeOf(fv), reflect.ValueOf(sv).Len(), reflect.ValueOf(sv).Len())
				for i, ssv := range sv.([]interface{}) {
					a := afv.Index(i)
					ai := reflect.New(reflect.TypeOf(a.Interface())).Interface()
					if err := structToStruct(sm, ai); err == nil {
						a.Set(reflect.ValueOf(ai).Elem())
					}
					for k, mv := range ssv.(map[string]interface{}) {
						if kf := SearchFieldName(a.Interface(), fn+"."+k); kf != "" {
							ak := a.FieldByName(kf)
							akt := ak.Type()
							akv := reflect.ValueOf(mv)
							if !ak.IsValid() || !ak.CanSet() || akt != akv.Type() {
								continue
							}
							ak.Set(akv)
						}
					}
				}
				if err := SetFieldValue(destination, fn, afv.Interface()); err != nil {
					return err
				}
			}
		} else if reflect.TypeOf(fv).Kind() == reflect.Struct || reflect.TypeOf(fv).Kind() == reflect.Ptr {
			fv = reflect.New(reflect.TypeOf(fv)).Interface()
			if err := structToStruct(sm, fv); err == nil {
				if err := SetFieldValue(destination, fn, reflect.ValueOf(fv).Elem().Interface()); err != nil {
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
						if reflect.TypeOf(sv).Kind() == reflect.Slice || reflect.TypeOf(sv).Kind() == reflect.Array {
							break
						}
						if sv, ok = sv.(map[string]interface{})[k]; !ok {
							break
						}
					}
				}
				if ok {
					SetFieldValue(destination, fn, sv)
				}
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
					if reflect.TypeOf(sv).Kind() == reflect.Slice || reflect.TypeOf(sv).Kind() == reflect.Array {
						break
					}
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
		if strings.ToUpper(string(fn[0])) != string(fn[0]) {
			return errors.New("cannot return value obtained from unexported field or method")
		}
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
					if reflect.TypeOf(sv).Kind() == reflect.Slice || reflect.TypeOf(sv).Kind() == reflect.Array {
						break
					}
					if sv, ok = sv.(map[string]interface{})[k]; !ok {
						break
					}
				}
			}
			if ok {
				if reflect.TypeOf(sv).Kind() != reflect.Slice && reflect.TypeOf(sv).Kind() != reflect.Array {
					continue
				}
				afv := reflect.MakeSlice(reflect.TypeOf(fv), reflect.ValueOf(sv).Len(), reflect.ValueOf(sv).Len())
				for i, ssv := range sv.([]interface{}) {
					a := afv.Index(i)
					for k, mv := range ssv.(map[string]interface{}) {
						if kf := SearchFieldName(a.Interface(), fn+"."+k); kf != "" {
							ak := a.FieldByName(kf)
							akt := ak.Type()
							akv := reflect.ValueOf(mv)
							if !ak.IsValid() || !ak.CanSet() || akt != akv.Type() {
								continue
							}
							ak.Set(akv)
						}
					}
				}
				if err := SetFieldValue(d, fn, afv.Interface()); err != nil {
					return err
				}
			}
		} else if reflect.TypeOf(fv).Kind() == reflect.Struct || reflect.TypeOf(fv).Kind() == reflect.Ptr {
			fv = reflect.New(reflect.TypeOf(fv)).Interface()
			if err := structToStruct(sm, fv); err != nil {
				return err
			}
			if err := SetFieldValue(d, fn, fv); err != nil {
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
					if reflect.TypeOf(sv).Kind() == reflect.Slice || reflect.TypeOf(sv).Kind() == reflect.Array {
						break
					}
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
