package util

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/shopspring/decimal"
)

// RandNum
func RandNum(min, max int) int {
	if max < min || max == min {
		return min
	}
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min+1)
}

// SplitRedPacket 将总红包金额拆分成指定份数,同时限制每个红包的最大和最小值
func SplitRedPacket(packetCount int64, totalAmount, minAmount, maxAmount decimal.Decimal) ([]decimal.Decimal, error) {
	precision := int32(2) // 保留小数位数
	totalAmount, minAmount, maxAmount = totalAmount.Truncate(precision), minAmount.Truncate(precision), maxAmount.Truncate(precision)

	if packetCount <= 0 || totalAmount.Sign() != 1 || maxAmount.Sign() != 1 {
		return nil, fmt.Errorf("Invalid parameters: packetCount=%d, totalAmount=%s, maxAmount=%s", packetCount, totalAmount, maxAmount)
	}
	if minAmount.GreaterThanOrEqual(maxAmount) {
		return nil, fmt.Errorf("Min amount (%s) must be less than max amount (%s)", minAmount, maxAmount)
	}
	minTotalAmount := decimal.NewFromInt(packetCount).Mul(minAmount)
	// 总红包金额 必须大于 红包份数*每个红包的最小值
	if totalAmount.LessThanOrEqual(minTotalAmount) {
		return nil, fmt.Errorf("Total amount must be greater than to %s", minTotalAmount)
	}
	rand.Seed(time.Now().UnixNano())

	// 用于生成随机金额的函数
	generateRandom := func(amount decimal.Decimal) decimal.Decimal {
		divisor := decimal.New(1, precision)
		i := rand.Intn(int(amount.Mul(divisor).IntPart()) + 1)
		return decimal.NewFromInt(int64(i)).Div(divisor)
	}

	remaining := totalAmount.Sub(minTotalAmount)    // 剩余的总金额
	packets := make([]decimal.Decimal, packetCount) // 初始化红包列表,每个红包初始为最小金额

	for i := range packets {
		extra := decimal.Zero

		if remaining.IsPositive() {
			randomAmount := maxAmount.Sub(minAmount) // 随机最大金额

			if remaining.LessThan(randomAmount) {
				randomAmount = remaining // 剩余总金额不够随机最大金额,直接使用剩余总金额随机
			}
			extra = generateRandom(randomAmount)
			remaining = remaining.Sub(extra)
		}

		packets[i] = minAmount.Add(extra)
	}
	return packets, nil
}
