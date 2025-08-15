package util_test

import (
	"testing"

	"github.com/scrawld/library/util"
	"github.com/shopspring/decimal"
)

// TestGenRandomAmount_DecreasingQuota 测试多次调用 GenRandomAmount 时 quota 逐渐变小
func TestGenRandomAmount_DecreasingQuota(t *testing.T) {
	quota := decimal.NewFromFloat(100.0) // 初始额度
	minRatio := 0.1
	maxRatio := 0.3

	for i := range 30 {
		amount := util.GenRandomAmount(quota, minRatio, maxRatio)
		t.Logf("Iteration %d: generated amount=%v, remaining quota=%v", i+1, amount, quota)

		// 检查金额不超过剩余额度
		if amount.GreaterThan(quota) {
			t.Errorf("Generated amount %v > remaining quota %v", amount, quota)
		}

		// 更新剩余额度
		quota = quota.Sub(amount)
	}

	t.Logf("Final remaining quota: %v", quota)
}

// TestCalcZeroProbability 测试 CalcZeroProbability 方法
func TestCalcZeroProbability(t *testing.T) {
	tests := []struct {
		quota float64
		min   float64
		max   float64
	}{
		{-1, 1.0, 1.0},     // quota <= 0 应返回 1
		{0, 1.0, 1.0},      // quota = 0 应返回 1
		{0.01, 0.05, 0.95}, // 边界小值
		{0.05, 0.05, 0.95}, // 边界小值
		{1, 0.05, 0.95},    // 常规值
		{10, 0.05, 0.95},   // 较大值
		{100, 0.05, 0.95},  // 很大值
		{1000, 0.05, 0.95}, // 超大值
	}

	for _, tt := range tests {
		prob := util.CalcZeroProbability(tt.quota)
		t.Logf("quota=%f => prob=%f", tt.quota, prob)

		if prob < tt.min || prob > tt.max {
			t.Errorf("CalcZeroProbability(%f) = %f, expected in range [%f, %f]",
				tt.quota, prob, tt.min, tt.max)
		}
	}
}
