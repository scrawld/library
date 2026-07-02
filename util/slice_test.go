package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestSliceDiff(t *testing.T) {
	// 有交集也有差集
	assert.ElementsMatch(t, []int{1, 4}, SliceDiff([]int{1, 2, 3}, []int{2, 3, 4}))

	// 完全相同，无差集
	assert.Nil(t, SliceDiff([]int{1, 2}, []int{1, 2}))

	// 一方为空
	assert.ElementsMatch(t, []int{1, 2}, SliceDiff([]int{1, 2}, []int{}))

	// 两个都为空
	assert.Nil(t, SliceDiff([]int{}, []int{}))
}

func TestSliceUnique(t *testing.T) {
	// 包含重复元素
	assert.Equal(t, []int{1, 2, 3}, SliceUnique([]int{1, 2, 2, 3, 3, 3}))

	// 无重复元素
	assert.Equal(t, []int{1, 2, 3}, SliceUnique([]int{1, 2, 3}))

	// 空切片 / nil
	assert.Nil(t, SliceUnique([]int{}))
	assert.Nil(t, SliceUnique[int](nil))

	// 字符串类型
	assert.Equal(t, []string{"a", "b", "c"}, SliceUnique([]string{"a", "b", "a", "c", "b"}))
}

func TestSliceShuffle(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7}
	SliceShuffle(s)
	t.Logf("s: %v", s)
}

func TestSliceChunk(t *testing.T) {
	assert.Equal(t, [][]int{{1, 2}, {3, 4}, {5}}, SliceChunk([]int{1, 2, 3, 4, 5}, 2))

	// size大于len
	assert.Equal(t, [][]int{{1, 2, 3}}, SliceChunk([]int{1, 2, 3}, 5))

	// size为0
	assert.Nil(t, SliceChunk([]int{1, 2, 3}, 0))

	// size为负
	assert.Nil(t, SliceChunk([]int{1, 2, 3}, -1))

	// 空切片
	assert.Nil(t, SliceChunk([]int{}, 3))
}
