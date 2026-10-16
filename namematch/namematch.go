package namematch

import (
	"github.com/xrash/smetrics"
)

// Similarity 返回两个姓名的相似度(0~100),忽略大小写/标点和词序
func Similarity(a, b string) float64 {
	a, b = normalize(a), normalize(b)
	if a == "" || b == "" {
		return 0
	}
	return smetrics.JaroWinkler(a, b, 0.7, 4) * 100
}
