package util

import (
	"math/rand"
	"slices"
	"time"
)

// InSlice 函数搜索数组中是否存在指定的值
func InSlice[T comparable](val T, slice []T) (bool, int) {
	for k, v := range slice {
		if val == v {
			return true, k
		}
	}
	return false, -1
}

// SliceDiff 函数用于比较两个数组的值，并返回差集
func SliceDiff[T comparable](slice1, slice2 []T) (r []T) {
	for _, v := range slice1 {
		if slices.Contains(slice2, v) {
			continue
		}
		r = append(r, v)
	}
	for _, v := range slice2 {
		if slices.Contains(slice1, v) {
			continue
		}
		r = append(r, v)
	}
	return
}

// SliceUnique 函数用于移除数组中重复的值
func SliceUnique[T comparable](slice []T) (r []T) {
	m := map[T]struct{}{}
	for _, v := range slice {
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = struct{}{}
		r = append(r, v)
	}
	return
}

// SliceShuffle 函数把数组中的元素按随机顺序重新排列
func SliceShuffle[T any](slice []T) {
	rand.Seed(time.Now().UnixNano())
	for i := len(slice) - 1; i > 0; i-- {
		k := rand.Intn(i + 1)
		slice[k], slice[i] = slice[i], slice[k]
	}
	return
}

// SliceChunk 函数把一个数组分割为新的数组块
func SliceChunk[T comparable](slice []T, size int) (r [][]T) {
	if size <= 0 || len(slice) == 0 {
		return nil
	}
	for size < len(slice) {
		r, slice = append(r, slice[0:size:size]), slice[size:]
	}
	r = append(r, slice)
	return
}
