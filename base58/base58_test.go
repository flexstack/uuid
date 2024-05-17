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
