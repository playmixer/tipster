package tipster

import (
	"crypto/sha256"
	"fmt"
)

func hashBtoS(b []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(b))
}
