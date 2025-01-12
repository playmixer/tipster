package tipster

import (
	"crypto/sha256"
	"fmt"
	"time"

	"golang.org/x/exp/rand"
)

var (
	alphabetLower = "abcdefghijklmnopqrstuvwxyz"
	alphabetUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers       = "1234567890"
)

func init() {
	rand.Seed(uint64(time.Now().Unix()))
}

func hashBtoS(b []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(b))
}

// RandomString генерирует строку заданой длины.
func randomString(n uint, text string) string {
	var letterRunes = []byte(text)
	b := make([]byte, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
