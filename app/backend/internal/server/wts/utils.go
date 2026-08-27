package wts

import "math"

func safeInt8(n int) int8 {
	if n > math.MaxInt8 {
		return math.MaxInt8
	}
	return int8(n)
}
