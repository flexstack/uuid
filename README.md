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

// Scan a SQL UUID
var u uuid.UUID
err := db.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&u)
if u.IsNil() {
    // Handle nil row
}
```


## Credit

This package is a fork of [github.com/google/uuid](https://github.com/gofrs/uuid) with the following changes:

- 2x improvement to `FromString`, `UnmarshalText`, and `UnmarshalJSON` performance
- Adds base58 encoding.
- Allows people to set a default format (i.e. base58, hash, canonical)
- Scans nil UUIDs from SQL databases as nil UUIDs (00000000-0000-0000-0000-000000000000) instead of `nil`.
- Fixes issue with [TimestampFromV7](https://github.com/gofrs/uuid/issues/128) not being spec compliant.
- Removed v1, v3, v5 UUIDs.
- Removed support for braced and URN string formats.
