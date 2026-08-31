package cachetime

import (
	"math/rand/v2"
	"time"
)

// WithJitter 回傳 base 加上 [-spread, +spread] 範圍內隨機抖動後的時間長度。
// 例如 WithJitter(30*time.Minute, 5*time.Minute) 會回傳 25~35 分鐘之間的隨機值。
func WithJitter(base time.Duration, spread time.Duration) time.Duration {
	if spread <= 0 {
		return base
	}

	// 產生 [-spread, +spread] 之間的隨機偏移量
	offsetRange := int64(spread) * 2
	offset := rand.Int64N(offsetRange) - int64(spread)

	result := base + time.Duration(offset)
	if result < 0 {
		result = 0
	}
	return result
}

// WithJitterRatio 依比例抖動,例如 ratio=0.2 代表在 base 的 ±20% 範圍內隨機。
// 比起固定數值的 spread,這個寫法在 base 值調整時不用跟著重新計算 spread,較不容易忘記同步改動。
func WithJitterRatio(base time.Duration, ratio float64) time.Duration {
	if ratio <= 0 {
		return base
	}
	spread := time.Duration(float64(base) * ratio)
	return WithJitter(base, spread)
}
