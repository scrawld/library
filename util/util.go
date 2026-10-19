package util

import (
	"fmt"
	"reflect"
)

// Assign 结构体同名字段赋值,excludes 为排除的字段
func Assign(origin, target any, excludes ...string) error {
	var (
		originVal, targetVal   = reflect.ValueOf(origin), reflect.ValueOf(target)
		originKind, targetKind = originVal.Kind(), targetVal.Kind()
	)
	if originKind != reflect.Pointer || targetKind != reflect.Pointer {
		return fmt.Errorf("origin and target must be pointers, current %s %s", originKind, targetKind)
	}
	originVal, targetVal = originVal.Elem(), targetVal.Elem()
	originKind, targetKind = originVal.Kind(), targetVal.Kind()

	if originKind != reflect.Struct || targetKind != reflect.Struct {
		return fmt.Errorf("origin and target must be struct, current %s %s", originKind, targetKind)
	}

	excludeMap := map[string]struct{}{}
	for _, v := range excludes {
		excludeMap[v] = struct{}{}
	}
	for i := range originVal.NumField() {
		if _, isExclude := excludeMap[originVal.Type().Field(i).Name]; isExclude {
			continue
		}
		var (
			originFieldVal = originVal.Field(i)
			targetFieldVal = targetVal.FieldByName(originVal.Type().Field(i).Name)
		)
		if !originFieldVal.CanSet() || !targetFieldVal.CanSet() {
			continue
		}
		if originFieldVal.Type() != targetFieldVal.Type() {
			continue
		}
		targetFieldVal.Set(originFieldVal)
	}
	return nil
}
