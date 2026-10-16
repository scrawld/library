package namematch

import (
	"sort"
	"strings"
)

// 分隔符换成空格,撇号直接删除
var punctReplacer = strings.NewReplacer(
	",", " ",
	".", " ",
	"-", " ",
	"'", "",
)

// normalize 转小写、去标点、拆词排序后用空格拼接
func normalize(s string) string {
	s = strings.ToLower(s)
	s = punctReplacer.Replace(s)
	f := strings.Fields(s)
	sort.Strings(f)
	return strings.Join(f, " ")
}
