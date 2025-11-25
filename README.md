# Shane's Awesome Shopping Kart (SASK)

Welcome to probably the greatest shopping cart that this planet, nay, galaxy,
has ever seen.

## Requirements

- Go 1.25 or later
- No external services required (uses in-memory storage)

## Installation
```bash
git clone https://github.com/shaneHowearth/kart.git
cd kart
go mod download
```

#### API

The API server can be run with
```
$ go run cmd/main.go
```

That will start the server listening on port 8080.

Docker configuration has not been included.

### Domains
There are two domains, [product](https://github.com/shaneHowearth/kart/blob/main/product) and [order](https://github.com/shaneHowearth/kart/blob/main/order).

Product has been pre-seeded with a small number of products, which can be seen
in the [seed file](https://github.com/shaneHowearth/kart/blob/main/product/datastore/seed.go)

### Testing

**End to end tests:**
- [Product E2E test](https://github.com/shaneHowearth/kart/blob/main/runProductTests.sh)
- [Order E2E test](https://github.com/shaneHowearth/kart/blob/main/runOrderTests.sh)

**Unit tests**
```bash
$ go test -v ./...
```

## Bonus round.

There is a [Promotional Coupons](https://github.com/shaneHowearth/kart/tree/main/cmd/coupons)
search tool, that can search text files for valid or invalid promotional codes.

### Usage
```bash
go run cmd/coupons/main.go   [file2] [file3]...
```

### Requirements

- Files must be uncompressed plain text (`.gz` files are rejected)
- Files must contain uppercase codes only
- Pattern matching is case-insensitive (patterns are converted to uppercase)

### Performance

- First search: ~7 seconds (searching 3 1GB files on an M4 MBP)
- Subsequent searches: <1ms (cached in SQLite)
- Cache file: `promotion_data.db` (created in current directory)

To clear the cache:
```bash
rm promotion_data.db
```
