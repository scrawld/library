package util

import (
	"reflect"

	"github.com/shopspring/decimal"
)

// StructDiff
func StructDiff(s1, s2 interface{}, excludes ...string) (before, after map[string]interface{}) {
	before, after = map[string]interface{}{}, map[string]interface{}{}

	valS1 := reflect.ValueOf(s1)
	valS2 := reflect.ValueOf(s2)
	if valS1.Kind() == reflect.Ptr {
		valS1 = valS1.Elem()
	}
	if valS2.Kind() == reflect.Ptr {
		valS2 = valS2.Elem()
	}

	for i := 0; i < valS1.NumField(); i++ {
		is_exclude := false
		for _, col := range excludes {
			if valS1.Type().Field(i).Name == col {
				is_exclude = true
				break
			}
		}
		if is_exclude {
			continue
		}

		tmpS1 := valS1.Field(i)
		tmpS2 := valS2.FieldByName(valS1.Type().Field(i).Name)

		if tmpS1.Kind() != tmpS2.Kind() {
			continue
		}
		switch tmpS1.Kind() {
		case reflect.Slice:
			//d := SliceDiff(tmpS1.Interface()), tmpS2.Interface())
			//if len(d) == 0 {
			//	continue
			//}
			continue
		case reflect.Struct:
			if tmpS1.Type().Name() == "Decimal" {
				dec1, ok1 := tmpS1.Interface().(decimal.Decimal)
				dec2, ok2 := tmpS2.Interface().(decimal.Decimal)
				if ok1 && ok2 && dec1.Equal(dec2) {
					continue
				}
			} else {
				continue
			}
		case reflect.Map:
			continue
		default:
			if tmpS1.Interface() == tmpS2.Interface() {
				continue
			}
		}
		key := valS1.Type().Field(i).Tag.Get("json")
		if len(key) == 0 {
			key = valS1.Type().Field(i).Name
		}
		before[key] = tmpS1.Interface()
		after[key] = tmpS2.Interface()
	}
	return
}
