package booking

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"
)

var counter uint64

func GenerateBookingReference() string {

	// date prefix
	date := time.Now().Format("20060102")

	// atomic counter prevents collisions inside the same process
	c := atomic.AddUint64(&counter, 1)

	// random bytes add extra entropy
	b := make([]byte, 2)
	rand.Read(b)

	randomPart := hex.EncodeToString(b)

	return fmt.Sprintf(
		"HTL-%s-%s%X",
		date,
		randomPart,
		c,
	)
}