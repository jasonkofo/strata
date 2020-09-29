package sqlgen

import (
	"math/rand"
	"time"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()),
)

// randomString definition
func randomString(length int) string {
	b := make([]byte, length)

	for i := 0; i < length; i++ {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

const charset = "abcdefghijklmnopqrstuvwxyz"
