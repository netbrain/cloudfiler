package web

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type DataInjector struct {
	source map[string][]string
	target reflect.Value
}

func NewDataInjector(source map[string][]string, target interface{}) *DataInjector {

	if source == nil || target == nil {
		return nil
	}

	val := reflect.ValueOf(target)

	if val.Kind() != reflect.Ptr {
		panic("target must be ptr")
	}

	if val.IsNil() || !val.IsValid() {
		return nil
	}

	return &DataInjector{
		source: source,
		target: val,
	}
}

func (f *DataInjector) Inject() {
	target := f.target.Elem()
	for i := target.NumField() - 1; i >= 0; i-- {
		targetFieldTyp := target.Type().Field(i)
		targetFieldVal := target.FieldByName(targetFieldTyp.Name)
		if !targetFieldVal.CanSet() {
			//skip unsettable field as it
			//is probably unexported
			continue
		}

		var sfName string
		if tagName := targetFieldTyp.Tag.Get("name"); tagName != "" {
			sfName = tagName
		} else {
			sfName = f.nameAsCamelBack(targetFieldTyp.Name)
		}
		sourceVal := f.source[sfName]

		if val := f.convertValues(targetFieldVal, sourceVal); val != nil {
			targetFieldVal.Set(reflect.ValueOf(val))
		}
	}
}

func (f *DataInjector) convertValues(t reflect.Value, s []string) (out interface{}) {
	if s == nil {
		return nil
	}

	if len(s) == 1 {
		in := s[0]
		switch t.Kind() {
		case reflect.String:
			out = in
		case reflect.Int:
			out, _ = strconv.Atoi(in)
		case reflect.Bool:
			out, _ = strconv.ParseBool(in)
		case reflect.Slice:
			switch t.Type().Elem().Kind() {
			case reflect.Uint8: //byte
				out = []byte(in)
			case reflect.Int: //int
				i, _ := strconv.Atoi(in)
				out = []int{i}
			case reflect.String:
				out = s
			default:
				panic("Unknown slice element kind")
			}
		default:
			panic("Unknown kind")
		}
	} else {
		in := s
		if t.Kind() == reflect.Slice {
			switch t.Type().Elem().Kind() {
			case reflect.String:
				out = in
			case reflect.Int:
				ints := make([]int, 0)
				for _, str := range in {
					i, _ := strconv.Atoi(str)
					ints = append(ints, i)
				}
				out = ints
			case reflect.Bool:
				bools := make([]bool, 0)
				for _, str := range in {
					b, _ := strconv.ParseBool(str)
					bools = append(bools, b)
				}
				out = bools
			default:
				panic("Unknown kind")
			}
		} else {
			panic(fmt.Sprintf("Expected slice target for input %#v", s))
		}
	}
	return
}

func (f *DataInjector) nameAsCamelBack(name string) string {
	nameParts := strings.SplitN(name, "", 2)
	nameParts[0] = strings.ToLower(nameParts[0])
	return strings.Join(nameParts, "")
}
