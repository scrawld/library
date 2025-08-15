package util

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/shopspring/decimal"
)

// RandNum 生成 int 类型的随机数
func RandNum(min, max int) int {
	if max < min || max == min {
		return min
	}
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min+1)
}

// RandNumInt64 生成 int64 类型的随机数
func RandNumInt64(min, max int64) int64 {
	if max <= min {
		return min
	}
	return min + rand.Int63n(max-min+1)
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

// GenerateRandomLetters 生成指定长度的随机字母的随机码(不包含o/O/I/l)
func GenerateRandomLetters(length uint) string {
	rand.Seed(time.Now().UnixNano())
	charset := "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"

	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

// GenerateRandomCode 生成指定长度的随机字母和数字的随机码(不包含o/O/I/l和0)
func GenerateRandomCode(length uint) string {
	rand.Seed(time.Now().UnixNano())
	charset := "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ123456789"

	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

// GenRandomAmount 根据剩余额度生成随机金额（永不到达最大值）
// quota: 剩余额度
// minRatio: 最小比例
// maxRatio: 最大比例
func GenRandomAmount(quota decimal.Decimal, minRatio, maxRatio float64) decimal.Decimal {
	precision := int32(2) // 保留小数位数

	quota = quota.Truncate(precision)
	if quota.Sign() <= 0 {
		return decimal.Zero
	}
	// 根据剩余额度计算零值概率 
	if rand.Float64() < CalcZeroProbability(quota.InexactFloat64()) {
		return decimal.Zero
	}
	// 将 quota 转为 float64 用于计算
	quotaVal := quota.InexactFloat64()

	// 计算上下限
	minAmt := math.Max(0.01, quotaVal*minRatio)
	maxAmt := quotaVal * maxRatio

	// 如果最小值大于等于最大值，直接返回零
	if minAmt >= maxAmt {
		return decimal.Zero
	}

	// 生成随机金额
	rangeSize := maxAmt - minAmt
	zoneRand := rand.Float64()
	var amount float64

	switch {
	case zoneRand < 0.25: // 偏小
		amount = minAmt + rand.Float64()*rangeSize*0.3
	case zoneRand < 0.75: // 中等
		amount = minAmt + rangeSize*0.3 + rand.Float64()*rangeSize*0.4
	default: // 偏大
		amount = minAmt + rangeSize*0.7 + rand.Float64()*rangeSize*0.3
	}
	if amount >= quotaVal {
		return decimal.Zero
	}
	return decimal.NewFromFloat(amount).Truncate(precision)
}

// CalcZeroProbability 零值概率计算（剩余越少，零概率越高）
func CalcZeroProbability(quota float64) float64 {
	if quota <= 0 {
		return 1.0
	}
	zeroProbBase := 0.05  // 基础零值概率
	zeroProbGrowth := 0.8 // 零值概率增长率

	quota = math.Max(quota, 0.01)
	prob := zeroProbBase + (1 - math.Exp(-zeroProbGrowth/quota))
	return math.Min(0.95, math.Max(zeroProbBase, prob))
}
