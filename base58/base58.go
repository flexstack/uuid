package base58

// Alphabet is a a b58 alphabet.
type Alphabet struct {
	decode [128]int8
	encode [58]byte
}

// Bitcoin base58 alphabet.
var encode = [58]byte{
	'1', '2', '3', '4', '5', '6', '7', '8', '9', 'A',
	'B', 'C', 'D', 'E', 'F', 'G', 'H', 'J', 'K', 'L',
	'M', 'N', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W',
	'X', 'Y', 'Z', 'a', 'b', 'c', 'd', 'e', 'f', 'g',
	'h', 'i', 'j', 'k', 'm', 'n', 'o', 'p', 'q', 'r',
	's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
}
var decode = [128]uint64{
	'1': 0, '2': 1, '3': 2, '4': 3, '5': 4,
	'6': 5, '7': 6, '8': 7, '9': 8, 'A': 9,
	'B': 10, 'C': 11, 'D': 12, 'E': 13, 'F': 14,
	'G': 15, 'H': 16, 'J': 17, 'K': 18, 'L': 19,
	'M': 20, 'N': 21, 'P': 22, 'Q': 23, 'R': 24,
	'S': 25, 'T': 26, 'U': 27, 'V': 28, 'W': 29,
	'X': 30, 'Y': 31, 'Z': 32, 'a': 33, 'b': 34,
	'c': 35, 'd': 36, 'e': 37, 'f': 38, 'g': 39,
	'h': 40, 'i': 41, 'j': 42, 'k': 43, 'm': 44,
	'n': 45, 'o': 46, 'p': 47, 'q': 48, 'r': 49,
	's': 50, 't': 51, 'u': 52, 'v': 53, 'w': 54,
	'x': 55, 'y': 56, 'z': 57,
}
var padLeft = [22]string{
	"",
	"1",
	"11",
	"111",
	"1111",
	"11111",
	"111111",
	"1111111",
	"11111111",
	"111111111",
	"1111111111",
	"11111111111",
	"111111111111",
	"1111111111111",
	"11111111111111",
	"111111111111111",
	"1111111111111111",
	"11111111111111111",
	"111111111111111111",
	"1111111111111111111",
	"11111111111111111111",
	"111111111111111111111",
}

var uuidSize = 16

func Decode(str string) ([]byte, error) {
	dst := make([]byte, uuidSize)
	if err := UnmarshalString(dst, str); err != nil {
		return nil, err
	}

	return dst, nil
}

func UnmarshalString(dst []byte, str string) error {
	return UnmarshalBytes(dst, []byte(str))
}

func UnmarshalBytes(dst, src []byte) error {
	// Use stack allocation for better performance
	var outi [4]uint32

	// Optimized for the common case of 22-byte base58 UUID
	if len(src) == 22 {
		// Unrolled loop for base58 decoding
		// Process all 22 characters with partially unrolled loop
		var c uint64

		// Unroll by 2 for better performance
		for i := 0; i < 22; i += 2 {
			// First character
			c = decode[src[i]]
			t3 := uint64(outi[3])*58 + c
			c = t3 >> 32
			outi[3] = uint32(t3)

			t2 := uint64(outi[2])*58 + c
			c = t2 >> 32
			outi[2] = uint32(t2)

			t1 := uint64(outi[1])*58 + c
			c = t1 >> 32
			outi[1] = uint32(t1)

			t0 := uint64(outi[0])*58 + c
			outi[0] = uint32(t0)

			// Second character (if exists)
			if i+1 < 22 {
				c = decode[src[i+1]]
				t3 = uint64(outi[3])*58 + c
				c = t3 >> 32
				outi[3] = uint32(t3)

				t2 = uint64(outi[2])*58 + c
				c = t2 >> 32
				outi[2] = uint32(t2)

				t1 = uint64(outi[1])*58 + c
				c = t1 >> 32
				outi[1] = uint32(t1)

				t0 = uint64(outi[0])*58 + c
				outi[0] = uint32(t0)
			}
		}
	} else {
		// Fallback for non-standard lengths
		for i := 0; i < len(src); i++ {
			c := decode[src[i]]

			for j := 3; j >= 0; j-- {
				t := uint64(outi[j])*58 + c
				c = t >> 32
				outi[j] = uint32(t)
			}
		}
	}

	// Unrolled output conversion
	dst[0] = byte(outi[0] >> 24)
	dst[1] = byte(outi[0] >> 16)
	dst[2] = byte(outi[0] >> 8)
	dst[3] = byte(outi[0])

	dst[4] = byte(outi[1] >> 24)
	dst[5] = byte(outi[1] >> 16)
	dst[6] = byte(outi[1] >> 8)
	dst[7] = byte(outi[1])

	dst[8] = byte(outi[2] >> 24)
	dst[9] = byte(outi[2] >> 16)
	dst[10] = byte(outi[2] >> 8)
	dst[11] = byte(outi[2])

	dst[12] = byte(outi[3] >> 24)
	dst[13] = byte(outi[3] >> 16)
	dst[14] = byte(outi[3] >> 8)
	dst[15] = byte(outi[3])

	return nil
}

func Encode(bin []byte) string {
	// A UUID will result in a base58 string of at most 22 characters.
	// This calculation is specific to 128-bit numbers (UUIDs).
	const maxEncodedSize = 22
	out := [22]byte{}
	var outIndex int = maxEncodedSize - 1 // Start filling from the end

	for i := 0; i < uuidSize; i++ {
		carry := uint32(bin[i])

		for j := maxEncodedSize - 1; j >= outIndex; j-- {
			carry += uint32(out[j]) * 256
			out[j] = byte(carry % 58)
			carry /= 58
		}

		for carry > 0 {
			outIndex--
			out[outIndex] = byte(carry % 58)
			carry /= 58
		}
	}

	for i := outIndex; i < maxEncodedSize; i++ {
		out[i] = encode[out[i]]
	}

	if outIndex == 0 {
		return string(out[:])
	}

	totalLen := 22 // Always 22 for padded result
	result := make([]byte, totalLen)

	// Fill padding with '1' characters
	for i := 0; i < outIndex; i++ {
		result[i] = '1'
	}

	copy(result[outIndex:], out[outIndex:])
	return string(result)
}
