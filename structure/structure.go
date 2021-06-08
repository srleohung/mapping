package structure

import (
	"errors"
	"reflect"
	"strings"

	"github.com/srleohung/mapping/kind"
)

// GetType is to get the type from the value
func GetType(i interface{}) reflect.Type {
	t := reflect.TypeOf(i)
	return getType(t)
}

func getType(t reflect.Type) reflect.Type {
	switch t.Kind() {
	case reflect.Ptr:
		return getType(t.Elem())
	default:
		return t
	}
}

// GetTypeName is to get the type name from the value
func GetTypeName(i interface{}) string {
	return reflect.TypeOf(i).Name()
}

// IsStruct is the check type is structure
func IsStruct(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Struct:
		return true
	default:
		return false
	}
}

// IsPublic is the check structure is public
func IsPublic(v reflect.Value) bool {
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).CanInterface() {
			return true
		}
	}
	return false
}

// GetValue is to get the value from the value
func GetValue(i interface{}) reflect.Value {
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr:
		return getValue(reflect.ValueOf(i).Elem())
	default:
		return reflect.ValueOf(i)
	}
}

func getValue(v reflect.Value) reflect.Value {
	switch v.Type().Kind() {
	case reflect.Ptr:
		return getValue(v.Elem())
	default:
		return v
	}
}

// GetFieldNames is to get all field names from the value
func GetFieldNames(i interface{}) (names []string) {
	if t := GetType(i); IsStruct(t) {
		for j := 0; j < t.NumField(); j++ {
			f := t.Field(j)
			names = append(names, f.Name)
		}
	}
	return names
}

// SearchFieldName is to search the field name from the value by key
func SearchFieldName(i interface{}, k, v string) string {
	if t := GetType(i); IsStruct(t) {
		for j := 0; j < t.NumField(); j++ {
			f := t.Field(j)
			a := strings.Split(f.Tag.Get(k), ",")
			for _, s := range a {
				if s == v {
					return f.Name
				}
			}
		}
	}
	return ""
}

// SearchFieldNames is to search all field names from the value by key
func SearchFieldNames(i interface{}, k, v string) (names []string) {
	if t := GetType(i); IsStruct(t) {
		for j := 0; j < t.NumField(); j++ {
			f := t.Field(j)
			a := strings.Split(f.Tag.Get(k), ",")
			for _, s := range a {
				if s == v {
					names = append(names, f.Name)
					break
				}
			}
		}
	}
	return names
}

// SetFieldValue is to set the field value on the structure
func SetFieldValue(i interface{}, f string, n interface{}) error {
	v := GetValue(i)
	return setFieldValue(v, f, n)
}

func setFieldValue(v reflect.Value, f string, n interface{}) error {
	if !IsStruct(v.Type()) {
		return errors.New("input structure error")
	}
	fv := v.FieldByName(f)
	if !fv.IsValid() {
		return errors.New("invalid field")
	}
	if !fv.CanSet() {
		return errors.New("cannot set structure field")
	}
	nv := reflect.ValueOf(n)
	if fv.Type() != nv.Type() {
		if nn, err := kind.ToType(n, fv.Type().Kind()); err == nil {
			nv = reflect.ValueOf(nn)
			if fv.Type() != nv.Type() {
				return errors.New("structure field type does not match value type")
			}
		} else {
			return errors.New("structure field type does not match value type")
		}
	}
	fv.Set(nv)
	return nil
}

// StructToMap is to convert the structure to a map
func StructToMap(s interface{}) map[string]interface{} {
	m, _ := structToMap(s)
	return m
}

func structToMap(s interface{}) (map[string]interface{}, error) {
	t := GetType(s)
	v := GetValue(s)
	m := make(map[string]interface{})
	if !IsStruct(t) {
		return m, errors.New("input structure error")
	}
	if !IsPublic(v) {
		return m, errors.New("cannot return value obtained from unexported field or method")
	}
	for i := 0; i < t.NumField(); i++ {
		fn := t.Field(i).Name
		if !v.Field(i).CanInterface() {
			continue
		}
		fv := getValue(v.Field(i)).Interface()
		switch GetType(fv).Kind() {
		case reflect.Slice:
			a := make([]interface{}, GetValue(fv).Len())
			for j := 0; j < GetValue(fv).Len(); j++ {
				if stm, err := structToMap(GetValue(fv).Index(j).Interface()); err == nil {
					a[j] = stm
				} else {
					a[j] = GetValue(fv).Index(j).Interface()
				}
			}
			m[fn] = a
		case reflect.Struct:
			if stm, err := structToMap(fv); err == nil {
				m[fn] = stm
			} else {
				m[fn] = fv
			}
		default:
			m[fn] = fv
		}
	}
	return m, nil
}

// StructToStruct is to transform a structure into another structure
func StructToStruct(s interface{}, d interface{}) error {
	sm := StructToMap(s)
	return structToStruct(sm, d)
}

func structToStruct(sm map[string]interface{}, d interface{}) error {
	const tag string = "struct"
	t := GetType(d)
	v := GetValue(d)
	if !IsStruct(t) {
		return errors.New("input structure error")
	}
	if !IsPublic(v) {
		return errors.New("cannot return value obtained from unexported field or method")
	}
	for i := 0; i < t.NumField(); i++ {
		fn := t.Field(i).Name
		ft := t.Field(i).Tag
		if !v.Field(i).CanInterface() {
			continue
		}
		fv := getValue(v.Field(i)).Interface()
		switch GetType(fv).Kind() {
		case reflect.Slice, reflect.Array:
			tvs := strings.Split(ft.Get(tag), ",")
			var sv interface{}
			var ok bool
			for _, tv := range tvs {
				ks := strings.Split(tv, ".")
				for i, k := range ks {
					if i == 0 {
						if sv, ok = sm[k]; !ok {
							break
						}
					} else {
						switch GetType(sv).Kind() {
						case reflect.Slice:
							break
						case reflect.Array:
							break
						default:
							if sv, ok = sv.(map[string]interface{})[k]; !ok {
								break
							}
						}
					}
				}
				if ok {
					break
				}
			}
			if ok {
				switch GetType(fv).Kind() {
				case reflect.Slice, reflect.Array:
					afv := reflect.MakeSlice(GetType(fv), GetValue(sv).Len(), GetValue(sv).Len())
					for i, ssv := range sv.([]interface{}) {
						a := afv.Index(i)
						ai := reflect.New(GetType(a.Interface())).Interface()
						if err := structToStruct(sm, ai); err == nil {
							a.Set(GetValue(ai))
						}
						for k, mv := range ssv.(map[string]interface{}) {
							if kf := SearchFieldNames(a.Interface(), tag, k); len(kf) != 0 {
								for _, kfv := range kf {
									setFieldValue(a, kfv, mv)
								}
							}
						}
					}
					if err := SetFieldValue(d, fn, afv.Interface()); err != nil {
						continue
					}
				default:
					continue
				}
			}
		case reflect.Struct:
			fv = reflect.New(GetType(fv)).Interface()
			if err := structToStruct(sm, fv); err == nil {
				if err := SetFieldValue(d, fn, GetValue(fv).Interface()); err != nil {
					continue
				}
			} else {
				tvs := strings.Split(ft.Get(tag), ",")
				var sv interface{}
				var ok bool
				for _, tv := range tvs {
					ks := strings.Split(tv, ".")
					for i, k := range ks {
						if i == 0 {
							if sv, ok = sm[k]; !ok {
								break
							}
						} else {
							switch GetType(sv).Kind() {
							case reflect.Slice:
								break
							case reflect.Array:
								break
							default:
								if sv, ok = sv.(map[string]interface{})[k]; !ok {
									break
								}
							}
						}
					}
					if ok {
						break
					}
				}
				if ok {
					SetFieldValue(d, fn, sv)
				}
			}
		default:
			tvs := strings.Split(ft.Get(tag), ",")
			var sv interface{}
			var ok bool
			for _, tv := range tvs {
				ks := strings.Split(tv, ".")
				for i, k := range ks {
					if i == 0 {
						if sv, ok = sm[k]; !ok {
							break
						}
					} else {
						switch GetType(sv).Kind() {
						case reflect.Slice:
							break
						case reflect.Array:
							break
						default:
							if sv, ok = sv.(map[string]interface{})[k]; !ok {
								break
							}
						}
					}
				}
				if ok {
					break
				}
			}
			if ok {
				SetFieldValue(d, fn, sv)
			}
		}
	}
	return nil
}
