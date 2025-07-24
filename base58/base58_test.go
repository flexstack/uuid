package base58

import (
	"crypto/rand"
	"encoding/hex"
	"testing"
)

type testValues struct {
	dec []byte
	enc string
}

func initTestPairs(n int) []testValues {
	var testPairs = make([]testValues, 0, n)

	// pre-make the test pairs, so it doesn't take up benchmark time...
	for i := 0; i < n; i++ {
		data := make([]byte, 16)
		rand.Read(data)
		testPairs = append(testPairs, testValues{dec: data, enc: Encode(data)})
	}

	return testPairs
}

func TestEncDecLoop(t *testing.T) {
	var b = make([]byte, 16)
	for i := 0; i < 100; i++ {
		rand.Read(b)
		fe := Encode(b)

		fd, ferr := Decode(fe)
		if ferr != nil {
			t.Errorf("fast error: %v", ferr)
		}

		if hex.EncodeToString(b) != hex.EncodeToString(fd) {
			t.Errorf("decoding err: %s != %s", hex.EncodeToString(b), hex.EncodeToString(fd))
		}
	}
}

func BenchmarkEncode(b *testing.B) {
	testPairs := initTestPairs(b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Encode(testPairs[i].dec)
	}
}

func BenchmarkDecode(b *testing.B) {
	testPairs := initTestPairs(b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Decode(testPairs[i].enc)
	}
}

var testCases = []string{
	"1C9z3nFjeJ44HMBeuqGNxt",
	"6ba7b8109dad11d180b400c04f",
	"Xk7pWZaRRFkqbVa3ma7F5f",
	"11111111111111111111EJ",
	"zzzzzzzzzzzzzzzzzzzzzz",
}

func BenchmarkUnmarshalBytesNew(b *testing.B) {
	dst := make([]byte, 16)
	src := []byte(testCases[0])

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = UnmarshalBytes(dst, src)
	}
}

func BenchmarkUnmarshalBytesNewMultiple(b *testing.B) {
	dst := make([]byte, 16)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			_ = UnmarshalBytes(dst, []byte(tc))
		}
	}
}
