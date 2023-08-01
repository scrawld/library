package util

import (
	"reflect"
	"strings"
)

// ParseTagSetting get model's field tags
func ParseTagSetting(tags reflect.StructTag, key ...string) map[string]string {
	setting := map[string]string{}
	for _, v := range key {
		str := tags.Get(v)
		if len(str) == 0 {
			continue
		}
		tags := strings.Split(str, ";")
		for _, value := range tags {
			if len(value) == 0 {
				continue
			}
			v := strings.Split(value, ":")
			k := strings.TrimSpace(strings.ToUpper(v[0]))
			if len(v) >= 2 {
				setting[k] = strings.Join(v[1:], ":")
			} else {
				setting[k] = k
			}
		}
	}
	return setting
}
