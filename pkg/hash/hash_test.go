package hash

import (
	"testing"
)

// Test data
var (
	testData       = []byte("hello world")
	expectedMD5    = "5eb63bbbe01eeed093cb22bb8f5acdc3"
	expectedSHA1   = "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed"
	expectedSHA256 = "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
	expectedSHA512 = "309ecc489c12d6eb4cc40f50c902f2b4d0ed77ee511a7c7a9bcd3ca86d4cd86f989dd35bc5ff499670da34255b45b0cfd830e81f605dcf7dc5542e93ae9cd76f"
)

// Testing the exposed hashing function
func TestHashSum(t *testing.T) {
	tests := []struct {
		name     string
		algo     string
		data     []byte
		expected string
	}{
		{"Convenience MD5", "md5", testData, expectedMD5},
		{"Convenience SHA1", "sha1", testData, expectedSHA1},
		{"Convenience SHA256", "sha256", testData, expectedSHA256},
		{"Convenience SHA512", "sha512", testData, expectedSHA512},
		{"Convenience unknown algo", "unknown", testData, expectedSHA256}, // Should default to SHA256
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HashSum(tt.algo, tt.data)
			if result != tt.expected {
				t.Errorf("HashSum() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Tests that the same input always produces the same output
func TestHashSumConsistency(t *testing.T) {
	// Hash the same data multiple times
	for i := 0; i < 10; i++ {
		result := HashSum("sha256", testData)
		if result != expectedSHA256 {
			t.Errorf("Consistency test failed at iteration %d: got %v, want %v", i, result, expectedSHA256)
		}
	}
}

// Benchmark tests
func BenchmarkMD5(b *testing.B) {

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HashSum("md5", testData)
	}
}

func BenchmarkSHA256(b *testing.B) {

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HashSum("sha256", testData)
	}
}
