package util

import (
	"math/rand"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.NewSource(time.Now().UnixNano())
}

func randomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandomUser() string {
	return randomString(6)
}

func RandomId() int64 {
	return randomInt(100000000000000000, 999999999999999999)
}
