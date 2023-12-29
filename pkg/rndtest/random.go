// Package rndtest вспомогательный модуль для тестирования.
package rndtest

import (
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// Int генерирует случайное значение типа int в диапазоне [min, max].
func Int(min, max int) int {
	return min + int(rand.Int63n(int64(max-min+1)))
}

// String генерирует случайную последовательность символов из алфавита (alphabet) длиной n.
func String(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}
