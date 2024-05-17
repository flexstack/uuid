// Copyright (C) 2013-2018 by Maxim Bublis <b@codemonkey.ru>
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package uuid

import (
	"errors"
	"fmt"

	"github.com/flexstack/uuid/base58"
)

// FromBytes returns a UUID generated from the raw byte slice input.
// It will return an error if the slice isn't 16 bytes long.
func FromBytes(input []byte) (UUID, error) {
	u := UUID{}
	err := u.UnmarshalBinary(input)
	return u, err
}

// FromBytesOrNil returns a UUID generated from the raw byte slice input.
// Same behavior as FromBytes(), but returns uuid.Nil instead of an error.
func FromBytesOrNil(input []byte) UUID {
	uuid, err := FromBytes(input)
	if err != nil {
		return Nil
	}
	return uuid
}

var errInvalidFormat = errors.New("uuid: invalid UUID format")

var hexLookupTable = func() [256]byte {
	var table [256]byte
	for i := range table {
		table[i] = 255 // Default to 255 for invalid characters
	}
	// Valid hexadecimal characters
	for i := '0'; i <= '9'; i++ {
		table[i] = byte(i - '0')
	}
	for i := 'a'; i <= 'f'; i++ {
		table[i] = byte(i - 'a' + 10)
	}
	for i := 'A'; i <= 'F'; i++ {
		table[i] = byte(i - 'A' + 10)
	}
	return table
}()

var canonicalByteRange = [16]byte{
	0, 2, 4, 6,
	9, 11,
	14, 16,
	19, 21,
	24, 26, 28, 30, 32, 34,
}

var hashByteRange = [16]byte{
	0, 2, 4, 6,
	8, 10,
	12, 14,
	16, 18,
	20, 22, 24, 26, 28, 30,
}

// Parse parses the UUID stored in the string text. Parsing and supported
// formats are the same as UnmarshalText.
func (u *UUID) Parse(s string) error {
	switch len(s) {
	case 22: // base58
		if err := base58.UnmarshalString(u[:], s); err != nil {
			return err
		}
		return nil

	case 32: // hash
		for i := 0; i < 32; i += 2 {
			v1 := hexLookupTable[s[i]]
			v2 := hexLookupTable[s[i+1]]
			if v1|v2 == 255 {
				return errInvalidFormat
			}
			u[i/2] = (v1 << 4) | v2
		}
		return nil

	case 36: // canonical
		if s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
			return fmt.Errorf("uuid: incorrect UUID format in string %q", s)
		}
		for i, x := range canonicalByteRange {
			v1 := hexLookupTable[s[x]]
			v2 := hexLookupTable[s[x+1]]
			if v1|v2 == 255 {
				return errInvalidFormat
			}
			u[i] = (v1 << 4) | v2
		}
		return nil

	default:
		return fmt.Errorf("uuid: incorrect UUID length %d in string %q", len(s), s)
	}
}

// FromString returns a UUID parsed from the input string.
// Input is expected in a form accepted by UnmarshalText.
func FromString(text string) (UUID, error) {
	var u UUID
	err := u.Parse(text)
	return u, err
}

// FromStringOrNil returns a UUID parsed from the input string.
// Same behavior as FromString(), but returns uuid.Nil instead of an error.
func FromStringOrNil(input string) UUID {
	uuid, err := FromString(input)
	if err != nil {
		return Nil
	}
	return uuid
}

// MarshalText implements the encoding.TextMarshaler interface.
// Creates a string representation of the UUID in the format specified by
// DefaultFormat.
func (u UUID) MarshalText() ([]byte, error) {
	switch DefaultFormat {
	case FormatCanonical:
		var buf [36]byte
		encodeCanonical(buf[:], u)
		return buf[:], nil
	case FormatHash:
		var buf [32]byte
		encodeHash(buf[:], u)
		return buf[:], nil
	default:
		var buf [22]byte
		copy(buf[:], base58.Encode(u[:]))
		return buf[:], nil
	}
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// Following formats are supported:
//
//	"6ba7b810-9dad-11d1-80b4-00c04fd430c8" (canonical)
//	"6ba7b8109dad11d180b400c04fd430c8" (hash)
//	"1C9z3nFjeJ44HMBeuqGNxt" (base58)
func (u *UUID) UnmarshalText(b []byte) error {
	switch len(b) {
	case 22: // base58
		if err := base58.UnmarshalBytes(u[:], b); err != nil {
			return err
		}
		return nil

	case 32: // hash
		for i := 0; i < 32; i += 2 {
			v1 := hexLookupTable[b[i]]
			v2 := hexLookupTable[b[i+1]]
			if v1|v2 == 255 {
				return errInvalidFormat
			}
			u[i/2] = (v1 << 4) | v2
		}
		return nil

	case 36: // canonical
		if b[8] != '-' || b[13] != '-' || b[18] != '-' || b[23] != '-' {
			return fmt.Errorf("uuid: incorrect UUID format in string %q", b)
		}
		for i, x := range canonicalByteRange {
			v1 := hexLookupTable[b[x]]
			v2 := hexLookupTable[b[x+1]]
			if v1|v2 == 255 {
				return errInvalidFormat
			}
			u[i] = (v1 << 4) | v2
		}
		return nil

	default:
		return fmt.Errorf("uuid: incorrect UUID length %d in string %q", len(b), b)
	}
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (u UUID) MarshalBinary() ([]byte, error) {
	return u[:], nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
// It will return an error if the slice isn't 16 bytes long.
func (u *UUID) UnmarshalBinary(data []byte) error {
	if len(data) != Size {
		return fmt.Errorf("uuid: UUID must be exactly 16 bytes long, got %d bytes", len(data))
	}
	copy(u[:], data)
	return nil
}
