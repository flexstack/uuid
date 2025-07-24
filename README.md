# uuid

Faster, more flexible v4 and v7 UUIDs in Go

## Features

- [x] Generate and parse v4 and v7 UUIDs
- [x] Canonical, hash, and base58 encoding
- [x] Select a default string format (i.e. base58, hash, canonical)
- [x] SQL scanning and JSON marshaling
- [x] The fastest UUID parsing available in Golang

## Installation

```bash
go get github.com/flexstack/uuid
```

## Usage

```go
import "github.com/flexstack/uuid"

// Optionally set a default format
uuid.DefaultFormat = uuid.FormatBase58

// Generate a new v4 UUID
u := uuid.Must(uuid.NewV4())

// Generate a new v7 UUID
u := uuid.Must(uuid.NewV7())

// Parse a UUID
u, err := uuid.FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

// Parse a UUID from a byte slice
u, err := uuid.FromBytes([]byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8})

// Format a UUID
asHash := u.Format(uuid.FormatHash)
asBase58 := u.Format(uuid.FormatBase58)
asCanonical := u.Format(uuid.FormatCanonical)

// Scan a SQL UUID
var u uuid.UUID
err := db.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&u)
if u.IsNil() {
    // Handle nil UUID
}
```

## Setting a default format

Changing the default format will affect how UUIDs are marshaled to strings from `MarshalText`, and `MarshalJSON`.


```go
import (
	"encoding/json"
	"fmt"

	"github.com/flexstack/uuid"
)

func main() {
	uuid.DefaultFormat = uuid.FormatBase58

	u := uuid.FromStringOrNil("ffffffff-ffff-ffff-ffff-ffffffffffff")

	// Marshal to base58
	m := map[string]uuid.UUID{"id": u}
	b, _ := json.Marshal(m)
	fmt.Println(string(b)) // {"id": "YcVfxkQb6JRzqk5kF2tNLv"}
}
```

## Credit

This package is a fork of [github.com/gofrs/uuid](https://github.com/gofrs/uuid) with the following changes:

- Adds base58 encoding.
- Allows people to set a default format (i.e. base58, hash, canonical)
- Scans nil UUIDs from SQL databases as nil UUIDs (00000000-0000-0000-0000-000000000000) instead of `nil`.
- Fixes issue with [TimestampFromV7](https://github.com/gofrs/uuid/issues/128) not being spec compliant.
- Removed v1, v3, v5 UUIDs.
- Removed support for braced and URN string formats.

## Performance optimizations

This library includes additional performance optimizations beyond the original fork:

- **Zero allocations** for all parsing operations
- **Optimized hex encoding/decoding** with lookup tables and unrolled loops
- **Optimized base58 decoding** with stack allocation and loop unrolling (~29% faster)

## Benchmarks

MacBook Air (15-inch, M2, 2023) Apple M2, 24GB RAM, MacOS 15.3.2

### UUID generation
```
BenchmarkNewV4         1516407	       790.7 ns/op	      16 B/op	       1 allocs/op
BenchmarkNewV7         1816982	       659.3 ns/op	      16 B/op	       1 allocs/op
```

### String Operations
```
BenchmarkString/canonical     	59380742	        20.03 ns/op	      48 B/op	       1 allocs/op
BenchmarkString/hash          	57661926	        20.09 ns/op	      32 B/op	       1 allocs/op
BenchmarkString/base58        	 5236279	        231.6 ns/op	      24 B/op	       1 allocs/op

BenchmarkFromBytes             504783348	        2.380 ns/op	       0 B/op	       0 allocs/op

BenchmarkFromString/canonical  153610305	         7.834 ns/op	     0 B/op	       0 allocs/op
BenchmarkFromString/hash       158399199	         7.480 ns/op	     0 B/op	       0 allocs/op
BenchmarkFromString/base58      24494169	        48.91 ns/op	     0 B/op	       0 allocs/op
```

## Contributing

Read the [CONTRIBUTING.md](CONTRIBUTING.md) guide to learn how to contribute to this project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
