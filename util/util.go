package util

import (
	"fmt"
	"reflect"
)

/**
 * Assign 结构体赋值
 *
 * Example:
 *
 * type User1 struct {
 * 	Name string
 * 	Age  int
 * }
 *
 * type User2 struct {
 * 	Name string
 * 	Age  int
 * }
 *
 * u1 := &User1{Name: "zhangsan", Age: 20}
 * u2 := &User2{}
 *
 * Assign(u1, u2)
 * fmt.Println(u1, u2)
 */
func Assign(origin, target interface{}, excludes ...string) error {
	var (
		originVal, targetVal   = reflect.ValueOf(origin), reflect.ValueOf(target)
		originKind, targetKind = originVal.Kind(), targetVal.Kind()
	)
	if originKind != reflect.Ptr || targetKind != reflect.Ptr {
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
	for i := 0; i < originVal.NumField(); i++ {
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
