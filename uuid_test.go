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
	"bytes"
	"fmt"
	"testing"
)

func TestUUID(t *testing.T) {
	t.Run("IsNil", testUUIDIsNil)
	t.Run("Bytes", testUUIDBytes)
	t.Run("StdString", testUUIDStdString)
	t.Run("String", testUUIDString)
	t.Run("StringDefaultFormat", testUUIDStringDefaultFormat)
	t.Run("Version", testUUIDVersion)
	t.Run("Variant", testUUIDVariant)
	t.Run("SetVersion", testUUIDSetVersion)
	t.Run("SetVariant", testUUIDSetVariant)
	t.Run("Format", testUUIDFormat)
}

func testUUIDIsNil(t *testing.T) {
	u := UUID{}
	got := u.IsNil()
	want := true
	if got != want {
		t.Errorf("%v.IsNil() = %t, want %t", u, got, want)
	}
}

func testUUIDBytes(t *testing.T) {
	got := codecTestUUID.Bytes()
	want := codecTestData
	if !bytes.Equal(got, want) {
		t.Errorf("%v.Bytes() = %x, want %x", codecTestUUID, got, want)
	}
}

func testUUIDStdString(t *testing.T) {
	got := NamespaceDNS.String()
	want := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	if got != want {
		t.Errorf("%v.String() = %q, want %q", NamespaceDNS, got, want)
	}
}

func testUUIDString(t *testing.T) {
	got := NamespaceDNS.Format()
	want := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	if got != want {
		t.Errorf("%v.String() = %q, want %q", NamespaceDNS, got, want)
	}
	if NamespaceDNS.String() != FromStringOrNil(want).String() {
		t.Errorf("%v.StdString() = %q, want %q", NamespaceDNS, NamespaceDNS.String(), FromStringOrNil(want).String())
	}
}

func testUUIDStringDefaultFormat(t *testing.T) {
	DefaultFormat = FormatBase58
	got := NamespaceDNS.Format()
	want := FromStringOrNil("6ba7b810-9dad-11d1-80b4-00c04fd430c8").Format(FormatBase58)
	if got != want {
		t.Errorf("%v.String() = %q, want %q", NamespaceDNS, got, want)
	}
	if NamespaceDNS.String() != FromStringOrNil(want).String() {
		t.Errorf("%v.StdString() = %q, want %q", NamespaceDNS, NamespaceDNS.String(), FromStringOrNil(want).String())
	}
	DefaultFormat = FormatCanonical
}

func testUUIDVersion(t *testing.T) {
	u := UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	if got, want := u.Version(), V1; got != want {
		t.Errorf("%v.Version() == %d, want %d", u, got, want)
	}
}

func testUUIDVariant(t *testing.T) {
	tests := []struct {
		u    UUID
		want byte
	}{
		{
			u:    UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			want: VariantNCS,
		},
		{
			u:    UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			want: VariantRFC4122,
		},
		{
			u:    UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			want: VariantMicrosoft,
		},
		{
			u:    UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xe0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			want: VariantFuture,
		},
	}
	for _, tt := range tests {
		if got := tt.u.Variant(); got != tt.want {
			t.Errorf("%v.Variant() == %d, want %d", tt.u, got, tt.want)
		}
	}
}

func testUUIDSetVersion(t *testing.T) {
	u := UUID{}
	want := V4
	u.SetVersion(want)
	if got := u.Version(); got != want {
		t.Errorf("%v.Version() == %d after SetVersion(%d)", u, got, want)
	}
}

func testUUIDSetVariant(t *testing.T) {
	variants := []byte{
		VariantNCS,
		VariantRFC4122,
		VariantMicrosoft,
		VariantFuture,
	}
	for _, want := range variants {
		u := UUID{}
		u.SetVariant(want)
		if got := u.Variant(); got != want {
			t.Errorf("%v.Variant() == %d after SetVariant(%d)", u, got, want)
		}
	}
}

func testUUIDFormat(t *testing.T) {
	val := Must(FromString("12345678-90ab-cdef-1234-567890abcdef"))
	tests := []struct {
		u    UUID
		f    Format
		want string
	}{
		{u: val, f: FormatCanonical, want: "12345678-90ab-cdef-1234-567890abcdef"},
		{u: val, f: FormatHash, want: "1234567890abcdef1234567890abcdef"},
		{u: val, f: FormatBase58, want: "3FP9ScdoVGyKrjtWQjQxDc"},
	}
	for _, tt := range tests {
		got := tt.u.Format(tt.f)
		if tt.want != got {
			t.Errorf(`Format("%s") got %s, want %s`, tt.f, got, tt.want)
		}
	}
}

func TestMust(t *testing.T) {
	sentinel := fmt.Errorf("uuid: sentinel error")
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("did not panic, want %v", sentinel)
		}
		err, ok := r.(error)
		if !ok {
			t.Fatalf("panicked with %T, want error (%v)", r, sentinel)
		}
		if err != sentinel {
			t.Fatalf("panicked with %v, want %v", err, sentinel)
		}
	}()
	fn := func() (UUID, error) {
		return Nil, sentinel
	}
	Must(fn())
}

func TestTimestampFromV7(t *testing.T) {
	tests := []struct {
		u       UUID
		want    int64
		wanterr bool
	}{
		{u: Must(NewV4()), wanterr: true},
		{u: Must(FromString("00000000-0000-7000-0000-000000000000")), want: 0x000000000000},
		{u: Must(FromString("018a8fec-3ced-7164-995f-93c80cbdc575")), want: 0x018a8fec3ced},
		{u: Must(FromString("ffffffff-ffff-7fff-ffff-ffffffffffff")), want: 0xffffffffffff},
	}
	for _, tt := range tests {
		got, err := TimestampFromV7(tt.u)

		switch {
		case tt.wanterr && err == nil:
			t.Errorf("TimestampFromV7(%v) want error, got %v", tt.u, got)

		case tt.want != got:
			t.Errorf("TimestampFromV7(%v) got %v, want %v", tt.u, got, tt.want)
		}
	}
}

func BenchmarkNewV4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Must(NewV4())
	}
}

func BenchmarkNewV7(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Must(NewV7())
	}
}
