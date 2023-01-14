package utils

import (
	"math/rand"
	"strconv"
	"time"
)

// 生成随机验证码
func GetRand(l int) string {
	rand.Seed(time.Now().UnixNano())
	s := ""
	for i := 0; i < l; i++ {
		s += strconv.Itoa(rand.Intn(10))
	}
	return s
}
