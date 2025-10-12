package utils

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/vegidio/go-sak/crypto"
)

const base62Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func FakeTimestamp(s string) time.Time {
	hash, err := crypto.Sha256Hash(s)
	if err != nil {
		return time.Now()
	}

	// Parse the 256-bit value.
	val := new(big.Int)
	if _, ok := val.SetString(hash, 16); !ok {
		return time.Now()
	}

	// Constants: [start, end] window in UTC.
	start := time.Date(1980, 10, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2035, 10, 1, 0, 0, 0, 0, time.UTC)

	// The total span in nanoseconds fits in int64 (~55 years).
	totalDur := end.Sub(start)
	totalNanos := totalDur.Nanoseconds()

	// max = 2^256 - 1
	maxValue := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

	// Scale: floor(val * totalNanos / max).
	numer := new(big.Int).Mul(val, big.NewInt(totalNanos))
	scaled := new(big.Int).Quo(numer, maxValue) // integer division

	// Convert to duration (safe: scaled <= totalNanos, which fits int64).
	offset := time.Duration(scaled.Int64())
	return start.Add(offset)
}

// region - Private functions

func base62ToBigInt(s string) (*big.Int, error) {
	n := big.NewInt(0)
	base := big.NewInt(62)

	for i := 0; i < len(s); i++ {
		c := s[i]
		idx := int64(strings.IndexByte(base62Alphabet, c))

		if idx < 0 {
			return nil, fmt.Errorf("invalid base62 char %q at pos %d", c, i)
		}
		n.Mul(n, base)
		n.Add(n, big.NewInt(idx))
	}
	return n, nil
}

// endregion
