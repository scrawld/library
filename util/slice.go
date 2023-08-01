package util

import (
	"math/rand"
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
	for i := 0; i < len(slice1); i++ {
		if exists, _ := InSlice(slice1[i], slice2); exists {
			continue
		}
		r = append(r, slice1[i])
	}
	for i := 0; i < len(slice2); i++ {
		if exists, _ := InSlice(slice2[i], slice1); exists {
			continue
		}
		r = append(r, slice2[i])
	}
	return
}

// SliceUnique 函数用于移除数组中重复的值
func SliceUnique[T comparable](slice []T) (r []T) {
	m := map[T]struct{}{}
	for i := 0; i < len(slice); i++ {
		t := slice[i]
		if _, ok := m[t]; ok {
			continue
		}
		m[t] = struct{}{}
		r = append(r, t)
	}
	return
}

// SliceShuffle 函数把数组中的元素按随机顺序重新排列
func SliceShuffle[T any](slice []T) {
	var ran = rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := len(slice) - 1; i > 0; i-- {
		k := ran.Intn(i + 1)
		slice[k], slice[i] = slice[i], slice[k]
	}
	return
}

// SliceChunk 函数把一个数组分割为新的数组块
func SliceChunk[T comparable](slice []T, size int) (r [][]T) {
	for size < len(slice) {
		r, slice = append(r, slice[0:size:size]), slice[size:]
	}
	r = append(r, slice)
	return
}
