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

// Package uuid provides implementations of the Universally Unique Identifier
// (UUID), as specified in RFC-4122 and the Peabody RFC Draft (revision 03).
//
// RFC-4122[1] provides the specification for versions 1, 3, 4, and 5. The
// Peabody UUID RFC Draft[2] provides the specification for the new k-sortable
// UUIDs, versions 6 and 7.
//
// DCE 1.1[3] provides the specification for version 2, but version 2 support
// was removed from this package in v4 due to some concerns with the
// specification itself. Reading the spec, it seems that it would result in
// generating UUIDs that aren't very unique. In having read the spec it seemed
// that our implementation did not meet the spec. It also seems to be at-odds
// with RFC 4122, meaning we would need quite a bit of special code to support
// it. Lastly, there were no Version 2 implementations that we could find to
// ensure we were understanding the specification correctly.
//
// [1] https://tools.ietf.org/html/rfc4122
// [2] https://datatracker.ietf.org/doc/html/draft-peabody-dispatch-new-uuid-format-03
// [3] http://pubs.opengroup.org/onlinepubs/9696989899/chap5.htm#tagcjh_08_02_01_01
package uuid

import (
	"fmt"

	"github.com/flexstack/uuid/base58"
)

// Size of a UUID in bytes.
const Size = 16

// UUID is an array type to represent the value of a UUID, as defined in RFC-4122.
type UUID [Size]byte

var DefaultFormat = FormatCanonical

// UUID versions.
const (
	_  byte = iota
	V1      // Version 1 (date-time and MAC address)
	V2      // Version 2 (date-time and MAC address, DCE security version)
	V3      // Version 3 (namespace name-based)
	V4      // Version 4 (random)
	V5      // Version 5 (namespace name-based)
	V6      // Version 6 (k-sortable timestamp and node ID) [peabody draft]
	V7      // Version 7 (k-sortable timestamp and random data) [peabody draft]
)

// UUID layout variants.
const (
	VariantNCS byte = iota
	VariantRFC4122
	VariantMicrosoft
	VariantFuture
)

// TimestampFromV7 returns the Timestamp embedded within a V7 UUID. This
// function returns an error if the UUID is any version other than 7.
//
// This is implemented based on revision 03 of the Peabody UUID draft, and may
// be subject to change pending further revisions. Until the final specification
// revision is finished, changes required to implement updates to the spec will
// not be considered a breaking change. They will happen as a minor version
// releases until the spec is final.
func TimestampFromV7(u UUID) (int64, error) {
	if u.Version() != 7 {
		return 0, fmt.Errorf("uuid: %s is version %d, not version 6", u, u.Version())
	}

	t := 0 |
		(int64(u[0]) << 40) |
		(int64(u[1]) << 32) |
		(int64(u[2]) << 24) |
		(int64(u[3]) << 16) |
		(int64(u[4]) << 8) |
		int64(u[5])

	return t, nil
}

// Nil is the nil UUID, as specified in RFC-4122, that has all 128 bits set to
// zero.
var Nil = UUID{}
var Omni = UUID{
	0xff, 0xff, 0xff, 0xff,
	0xff, 0xff,
	0xff, 0xff,
	0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
}

// Predefined namespace UUIDs.
var (
	NamespaceDNS  = Must(FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8"))
	NamespaceURL  = Must(FromString("6ba7b811-9dad-11d1-80b4-00c04fd430c8"))
	NamespaceOID  = Must(FromString("6ba7b812-9dad-11d1-80b4-00c04fd430c8"))
	NamespaceX500 = Must(FromString("6ba7b814-9dad-11d1-80b4-00c04fd430c8"))
)

// IsNil returns if the UUID is equal to the nil UUID
func (u UUID) IsNil() bool {
	return u == Nil
}

// Version returns the algorithm version used to generate the UUID.
func (u UUID) Version() byte {
	return u[6] >> 4
}

// Variant returns the UUID layout variant.
func (u UUID) Variant() byte {
	switch {
	case (u[8] >> 7) == 0x00:
		return VariantNCS
	case (u[8] >> 6) == 0x02:
		return VariantRFC4122
	case (u[8] >> 5) == 0x06:
		return VariantMicrosoft
	case (u[8] >> 5) == 0x07:
		fallthrough
	default:
		return VariantFuture
	}
}

// Bytes returns a byte slice representation of the UUID.
func (u UUID) Bytes() []byte {
	return u[:]
}

const hextable = "0123456789abcdef"

// encodeCanonical encodes the canonical RFC-4122 form of UUID u into the
// first 36 bytes dst. Loop unrolled for maximum performance.
func encodeCanonical(dst []byte, u UUID) {
	dst[8] = '-'
	dst[13] = '-'
	dst[18] = '-'
	dst[23] = '-'
	
	// Unrolled loop - each UUID byte becomes 2 hex chars
	// canonicalByteRange: [0, 2, 4, 6, 9, 11, 14, 16, 19, 21, 24, 26, 28, 30, 32, 34]
	
	c := u[0]; dst[0] = hextable[c>>4]; dst[1] = hextable[c&0x0f]
	c = u[1]; dst[2] = hextable[c>>4]; dst[3] = hextable[c&0x0f]
	c = u[2]; dst[4] = hextable[c>>4]; dst[5] = hextable[c&0x0f]
	c = u[3]; dst[6] = hextable[c>>4]; dst[7] = hextable[c&0x0f]
	c = u[4]; dst[9] = hextable[c>>4]; dst[10] = hextable[c&0x0f]
	c = u[5]; dst[11] = hextable[c>>4]; dst[12] = hextable[c&0x0f]
	c = u[6]; dst[14] = hextable[c>>4]; dst[15] = hextable[c&0x0f]
	c = u[7]; dst[16] = hextable[c>>4]; dst[17] = hextable[c&0x0f]
	c = u[8]; dst[19] = hextable[c>>4]; dst[20] = hextable[c&0x0f]
	c = u[9]; dst[21] = hextable[c>>4]; dst[22] = hextable[c&0x0f]
	c = u[10]; dst[24] = hextable[c>>4]; dst[25] = hextable[c&0x0f]
	c = u[11]; dst[26] = hextable[c>>4]; dst[27] = hextable[c&0x0f]
	c = u[12]; dst[28] = hextable[c>>4]; dst[29] = hextable[c&0x0f]
	c = u[13]; dst[30] = hextable[c>>4]; dst[31] = hextable[c&0x0f]
	c = u[14]; dst[32] = hextable[c>>4]; dst[33] = hextable[c&0x0f]
	c = u[15]; dst[34] = hextable[c>>4]; dst[35] = hextable[c&0x0f]
}

func encodeHash(dst []byte, u UUID) {
	// Unrolled loop for hash format (no dashes)
	// hashByteRange: [0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30]
	
	c := u[0]; dst[0] = hextable[c>>4]; dst[1] = hextable[c&0x0f]
	c = u[1]; dst[2] = hextable[c>>4]; dst[3] = hextable[c&0x0f]
	c = u[2]; dst[4] = hextable[c>>4]; dst[5] = hextable[c&0x0f]
	c = u[3]; dst[6] = hextable[c>>4]; dst[7] = hextable[c&0x0f]
	c = u[4]; dst[8] = hextable[c>>4]; dst[9] = hextable[c&0x0f]
	c = u[5]; dst[10] = hextable[c>>4]; dst[11] = hextable[c&0x0f]
	c = u[6]; dst[12] = hextable[c>>4]; dst[13] = hextable[c&0x0f]
	c = u[7]; dst[14] = hextable[c>>4]; dst[15] = hextable[c&0x0f]
	c = u[8]; dst[16] = hextable[c>>4]; dst[17] = hextable[c&0x0f]
	c = u[9]; dst[18] = hextable[c>>4]; dst[19] = hextable[c&0x0f]
	c = u[10]; dst[20] = hextable[c>>4]; dst[21] = hextable[c&0x0f]
	c = u[11]; dst[22] = hextable[c>>4]; dst[23] = hextable[c&0x0f]
	c = u[12]; dst[24] = hextable[c>>4]; dst[25] = hextable[c&0x0f]
	c = u[13]; dst[26] = hextable[c>>4]; dst[27] = hextable[c&0x0f]
	c = u[14]; dst[28] = hextable[c>>4]; dst[29] = hextable[c&0x0f]
	c = u[15]; dst[30] = hextable[c>>4]; dst[31] = hextable[c&0x0f]
}

type Format string

const (
	FormatCanonical Format = "canonical"
	FormatHash      Format = "hash"
	FormatBase58    Format = "base58"
)

// Format returns a string representation of the UUID in the specified format.
// If no format is specified, DefaultFormat is used.
func (u UUID) Format(format ...Format) string {
	f := DefaultFormat
	if len(format) > 0 {
		f = format[0]
	}

	switch f {
	case FormatCanonical:
		return u.String()
	case FormatHash:
		dst := make([]byte, 32)
		encodeHash(dst, u)
		return string(dst)
	case FormatBase58:
		return base58.Encode(u[:])
	default:
		return base58.Encode(u[:])
	}
}

// Returns a string representation of the UUID in the form of
// a canonical RFC-4122 string.
func (u UUID) String() string {
	dst := make([]byte, 36)
	encodeCanonical(dst, u)
	return string(dst)
}

// SetVersion sets the version bits.
func (u *UUID) SetVersion(v byte) {
	u[6] = (u[6] & 0x0f) | (v << 4)
}

// SetVariant sets the variant bits.
func (u *UUID) SetVariant(v byte) {
	switch v {
	case VariantNCS:
		u[8] = (u[8]&(0xff>>1) | (0x00 << 7))
	case VariantRFC4122:
		u[8] = (u[8]&(0xff>>2) | (0x02 << 6))
	case VariantMicrosoft:
		u[8] = (u[8]&(0xff>>3) | (0x06 << 5))
	case VariantFuture:
		fallthrough
	default:
		u[8] = (u[8]&(0xff>>3) | (0x07 << 5))
	}
}

// Must is a helper that wraps a call to a function returning (UUID, error)
// and panics if the error is non-nil. It is intended for use in variable
// initializations such as
//
//	var packageUUID = uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426655440000"))
func Must(u UUID, err error) UUID {
	if err != nil {
		panic(err)
	}
	return u
}
