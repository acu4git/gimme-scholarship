package testutil

import "math/rand"

func RandLetters(n int) string {
	rs1Letters := []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = rs1Letters[rand.Intn(len(rs1Letters))] //nolint:gosec
	}
	return string(b)
}
