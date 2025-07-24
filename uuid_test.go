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
	"sync"
	"testing"
	"time"
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

func TestConcurrentGeneration(t *testing.T) {
	const numGoroutines = 100
	const numUUIDs = 100

	var wg sync.WaitGroup
	results := make([][]UUID, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			uuids := make([]UUID, numUUIDs)
			for j := 0; j < numUUIDs; j++ {
				u, err := NewV4()
				if err != nil {
					t.Errorf("NewV4() failed in goroutine %d: %v", idx, err)
					return
				}
				uuids[j] = u
			}
			results[idx] = uuids
		}(i)
	}

	wg.Wait()

	// Check all UUIDs are unique across all goroutines
	seen := make(map[UUID]bool)
	for i, uuids := range results {
		if uuids == nil {
			continue // goroutine failed
		}
		for j, u := range uuids {
			if seen[u] {
				t.Errorf("Duplicate UUID found: %s (goroutine %d, index %d)", u, i, j)
			}
			seen[u] = true
		}
	}
}

func TestV7ConcurrentGeneration(t *testing.T) {
	const numGoroutines = 10
	const numUUIDs = 100

	var wg sync.WaitGroup
	results := make([][]UUID, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			uuids := make([]UUID, numUUIDs)
			for j := 0; j < numUUIDs; j++ {
				u, err := NewV7()
				if err != nil {
					t.Errorf("NewV7() failed in goroutine %d: %v", idx, err)
					return
				}
				uuids[j] = u
			}
			results[idx] = uuids
		}(i)
	}

	wg.Wait()

	// Check all UUIDs are unique and timestamps are reasonable
	seen := make(map[UUID]bool)
	for i, uuids := range results {
		if uuids == nil {
			continue
		}
		for j, u := range uuids {
			if seen[u] {
				t.Errorf("Duplicate V7 UUID found: %s (goroutine %d, index %d)", u, i, j)
			}
			seen[u] = true

			// Verify timestamp extraction works
			if _, err := TimestampFromV7(u); err != nil {
				t.Errorf("Failed to extract timestamp from V7 UUID %s: %v", u, err)
			}
		}
	}
}

func TestCustomEpochFunc(t *testing.T) {
	fixedTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	gen := &Gen{
		epochFunc: func() time.Time { return fixedTime },
	}

	u1, err := gen.NewV7()
	if err != nil {
		t.Fatalf("NewV7() failed: %v", err)
	}

	u2, err := gen.NewV7()
	if err != nil {
		t.Fatalf("NewV7() failed: %v", err)
	}

	ts1, _ := TimestampFromV7(u1)
	ts2, _ := TimestampFromV7(u2)

	expectedTS := uint64(fixedTime.UnixMilli())
	if uint64(ts1) != expectedTS {
		t.Errorf("First V7 UUID has timestamp %d, want %d", ts1, expectedTS)
	}
	if uint64(ts2) != expectedTS {
		t.Errorf("Second V7 UUID has timestamp %d, want %d", ts2, expectedTS)
	}

	// With same timestamp, clock sequence should increment
	if u1 == u2 {
		t.Errorf("Two V7 UUIDs with same timestamp should be different due to clock sequence")
	}
}

func TestV7RapidGeneration(t *testing.T) {
	// Generate many UUIDs rapidly in same millisecond to test clock sequence behavior
	uuids := make([]UUID, 1000)
	for i := 0; i < 1000; i++ {
		u, err := NewV7()
		if err != nil {
			t.Fatalf("NewV7() failed at iteration %d: %v", i, err)
		}
		uuids[i] = u
	}

	// All should be unique
	seen := make(map[UUID]bool)
	for i, u := range uuids {
		if seen[u] {
			t.Errorf("Duplicate V7 UUID at index %d: %s", i, u)
		}
		seen[u] = true
	}

	// Timestamps should be monotonic (non-decreasing)
	var lastTS int64
	for i, u := range uuids {
		ts, err := TimestampFromV7(u)
		if err != nil {
			t.Errorf("Failed to extract timestamp from UUID %d: %v", i, err)
			continue
		}
		if ts < lastTS {
			t.Errorf("V7 timestamp went backwards at index %d: %d < %d", i, ts, lastTS)
		}
		lastTS = ts
	}
}

func TestV4Correctness(t *testing.T) {
	for i := 0; i < 1000; i++ {
		u, err := NewV4()
		if err != nil {
			t.Fatalf("NewV4() failed: %v", err)
		}

		// Check version
		if u.Version() != V4 {
			t.Errorf("UUID %s has version %d, want %d", u, u.Version(), V4)
		}

		// Check variant
		if u.Variant() != VariantRFC4122 {
			t.Errorf("UUID %s has variant %d, want %d", u, u.Variant(), VariantRFC4122)
		}

		// Check not nil
		if u.IsNil() {
			t.Errorf("UUID should not be nil: %s", u)
		}

		// Check string format
		s := u.String()
		if len(s) != 36 {
			t.Errorf("UUID string %s has length %d, want 36", s, len(s))
		}

		// Check parsing back
		parsed, err := FromString(s)
		if err != nil {
			t.Errorf("Failed to parse UUID string %s: %v", s, err)
		}
		if parsed != u {
			t.Errorf("Parsed UUID %s != original %s", parsed, u)
		}
	}
}

func TestV7Correctness(t *testing.T) {
	var lastTime int64

	for i := 0; i < 100; i++ {
		u, err := NewV7()
		if err != nil {
			t.Fatalf("NewV7() failed: %v", err)
		}

		// Check version
		if u.Version() != V7 {
			t.Errorf("UUID %s has version %d, want %d", u, u.Version(), V7)
		}

		// Check variant
		if u.Variant() != VariantRFC4122 {
			t.Errorf("UUID %s has variant %d, want %d", u, u.Variant(), VariantRFC4122)
		}

		// Check timestamp extraction
		ts, err := TimestampFromV7(u)
		if err != nil {
			t.Errorf("Failed to extract timestamp from V7 UUID %s: %v", u, err)
		}

		// Check timestamp is reasonable (within last hour and future)
		now := time.Now().UnixMilli()
		if ts < now-3600000 || ts > now+1000 {
			t.Errorf("V7 UUID %s has unreasonable timestamp %d (now=%d)", u, ts, now)
		}

		// Check timestamps are non-decreasing (monotonic)
		if ts < lastTime {
			t.Errorf("V7 UUID timestamp went backwards: %d < %d", ts, lastTime)
		}
		lastTime = ts

		time.Sleep(1 * time.Millisecond) // Ensure different timestamps
	}
}

func TestUUIDUniqueness(t *testing.T) {
	seen := make(map[UUID]bool)

	// Test V4 uniqueness
	for i := 0; i < 10000; i++ {
		u, err := NewV4()
		if err != nil {
			t.Fatalf("NewV4() failed: %v", err)
		}
		if seen[u] {
			t.Errorf("Duplicate V4 UUID generated: %s", u)
		}
		seen[u] = true
	}

	// Test V7 uniqueness
	for i := 0; i < 1000; i++ {
		u, err := NewV7()
		if err != nil {
			t.Fatalf("NewV7() failed: %v", err)
		}
		if seen[u] {
			t.Errorf("Duplicate V7 UUID generated: %s", u)
		}
		seen[u] = true
	}
}

func TestStringFormats(t *testing.T) {
	u, err := NewV4()
	if err != nil {
		t.Fatalf("NewV4() failed: %v", err)
	}

	// Test canonical format
	canonical := u.Format(FormatCanonical)
	if len(canonical) != 36 {
		t.Errorf("Canonical format has wrong length: %d", len(canonical))
	}
	if canonical[8] != '-' || canonical[13] != '-' || canonical[18] != '-' || canonical[23] != '-' {
		t.Errorf("Canonical format missing dashes: %s", canonical)
	}

	// Test hash format
	hash := u.Format(FormatHash)
	if len(hash) != 32 {
		t.Errorf("Hash format has wrong length: %d", len(hash))
	}

	// Test base58 format
	base58 := u.Format(FormatBase58)
	if len(base58) == 0 {
		t.Errorf("Base58 format is empty")
	}

	// Test that String() matches canonical
	if u.String() != canonical {
		t.Errorf("String() != Format(FormatCanonical): %s vs %s", u.String(), canonical)
	}
}
