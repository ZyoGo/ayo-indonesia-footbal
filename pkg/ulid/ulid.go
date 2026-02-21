package ulid

import (
	"crypto/rand"
	"log"
	"math"
	"math/big"
	mrand "math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

var entropyPool = sync.Pool{
	New: func() interface{} {
		source := mrand.New(mrand.NewSource(cryptoRandSeed()))
		// Use ulid.Monotonic to ensure monotonicity within the same millisecond
		return ulid.Monotonic(source, 0)
	},
}

// cryptoRandSeed generates a cryptographically secure random seed
func cryptoRandSeed() int64 {
	max := big.NewInt(math.MaxInt64)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.Fatalf("Failed to seed random number generator: %v", err)
	}
	return n.Int64()
}

// GenerateID creates a new ULID string
func GenerateID() string {
	t := time.Now().UTC()

	// Retrieve a monotonic entropy source from the pool
	entropy := entropyPool.Get().(*ulid.MonotonicEntropy)

	// Generate a new ULID
	id := ulid.MustNew(ulid.Timestamp(t), entropy)

	// Put the entropy source back in the pool
	entropyPool.Put(entropy)

	return id.String()
}
