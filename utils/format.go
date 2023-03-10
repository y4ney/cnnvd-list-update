package utils

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
)

func FormatURL(domain, path string) string {
	return "https://" + domain + "/" + path
}

func FormatMonth(month int) string {
	if month < 10 {
		return fmt.Sprintf("0%v", month)
	}
	return fmt.Sprintf("%v", month)
}

// RandInt 返回一个[0,MaxInt64)的一个随机整数
func RandInt() int {
	seed, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	return int(seed.Int64())
}
