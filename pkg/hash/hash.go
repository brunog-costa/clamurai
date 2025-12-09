package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
)

type hasher struct {
	hasher hash.Hash
}

func New(algo string) *hasher {
	var h hash.Hash

	switch algo {
	case "md5":
		h = md5.New()
	case "sha1":
		h = sha1.New()
	case "sha256":
		h = sha256.New()
	case "sha512":
		h = sha512.New()
	default:
		// Default to sha256 if unknown algorithm
		h = sha256.New()
		algo = "sha256"
	}

	return &hasher{
		hasher: h,
	}
}

func (h *hasher) HashSum(body []byte) string {
	// Reset the hasher in case it's being reused
	h.hasher.Reset()

	// Write the data to the hasher
	h.hasher.Write(body)

	// Calculate the final hash sum
	sum := h.hasher.Sum(nil)

	// Convert to hexadecimal string
	return hex.EncodeToString(sum)
}

// Convenience function for one-time hashing
func HashSum(algo string, body []byte) string {
	hasher := New(algo)
	return hasher.HashSum(body)
}
