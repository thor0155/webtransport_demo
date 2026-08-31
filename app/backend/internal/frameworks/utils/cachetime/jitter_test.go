package cachetime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithJitter_WithinRange(t *testing.T) {
	base := 30 * time.Minute
	spread := 5 * time.Minute

	for i := 0; i < 1000; i++ {
		result := WithJitter(base, spread)
		assert.GreaterOrEqual(t, result, base-spread)
		assert.LessOrEqual(t, result, base+spread)
	}
}

func TestWithJitter_ZeroSpread(t *testing.T) {
	base := 30 * time.Minute
	assert.Equal(t, base, WithJitter(base, 0))
}

func TestWithJitter_NeverNegative(t *testing.T) {
	// 極端情況:spread 大於 base,結果不該變成負數
	for i := 0; i < 1000; i++ {
		result := WithJitter(1*time.Second, 10*time.Second)
		assert.GreaterOrEqual(t, result, time.Duration(0))
	}
}

func TestWithJitterRatio(t *testing.T) {
	base := 30 * time.Minute
	for i := 0; i < 1000; i++ {
		result := WithJitterRatio(base, 0.2) // ±20%
		assert.GreaterOrEqual(t, result, base-6*time.Minute)
		assert.LessOrEqual(t, result, base+6*time.Minute)
	}
}
