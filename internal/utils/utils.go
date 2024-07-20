package utils

import (
	"math/rand"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateString(n int) string {
	buf := make([]rune, n)
	for i := range n {
		buf[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(buf)
}
