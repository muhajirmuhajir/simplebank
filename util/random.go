package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generate random int betwen min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString return random name
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner return string
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney return int64
func RandomMoney() int64 {
	return RandomInt(0, 50000)
}

// RandomCurrency return string
func RandomCurrency() string {
	currencies := []string{"IDR", "USD", "EUR"}

	n := len(currencies)
	return currencies[rand.Intn(n)]
}
