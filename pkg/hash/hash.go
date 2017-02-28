package hash

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// CreateHash creates a new random SHA-1 hash.
func CreateHash() string {
	nowUnix := time.Now().UnixNano()
	number := rand.Int()
	pid := os.Getpid()

	data := fmt.Sprintf("%d%d%d", nowUnix, number, pid)
	hash := sha1.Sum([]byte(data))

	return fmt.Sprintf("%x", hash)
}
