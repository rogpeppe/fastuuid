# fastuuid
--
    import "github.com/rogpeppe/fastuuid"

Package fastuuid provides fast UUID generation of 192 bit universally unique
identifiers.

It also provides simple support for 128-bit RFC-4122-like UUID strings.

Note that the generated UUIDs are not unguessable - each UUID generated from a
Generator is adjacent to the previously generated UUID.

It ignores RFC 4122.

By way of comparison with another popular UUID-generation package,
github.com/satory/go.uuid, here are some benchmarks:

    BenchmarkNext-4              	130710132	         9.01 ns/op
    BenchmarkHex128-4            	13261449	        78.3 ns/op
    BenchmarkContended-4         	41227394	        29.5 ns/op
    BenchmarkSatoriNext-4        	 1000000	      1039 ns/op
    BenchmarkSatoriHex128-4      	 1000000	      1044 ns/op
    BenchmarkSatoriContended-4   	 1729365	       699 ns/op

## Usage

#### func  Hex128

```go
func Hex128(uuid [24]byte) string
```
Hex128 returns an RFC4122-like representation of the first 128 bits of the given
UUID. For example:

    f81d4fae-7dec-11d0-a765-00a0c91e6bf6.

It does not bother to set the version or variant bits of the UUID, as they only
serve to reduce the randomness.

If you want unpredictable UUIDs, you might want to consider hashing the uuid
(using SHA256, for example) before passing it to Hex128.

#### func  ValidHex128

```go
func ValidHex128(id string) bool
```
ValidHex128 reports whether id is a valid UUID as returned by Hex128 and various
other UUID packages, such as github.com/satori/go.uuid's NewV4 function.

Note that it does not allow upper case hex.

#### type Generator

```go
type Generator struct {
}
```

Generator represents a UUID generator that generates UUIDs in sequence from a
random starting point.

#### func  MustNewGenerator

```go
func MustNewGenerator() *Generator
```
MustNewGenerator is like NewGenerator but panics on failure.

#### func  NewGenerator

```go
func NewGenerator() (*Generator, error)
```
NewGenerator returns a new Generator. It can fail if the crypto/rand read fails.

#### func (*Generator) Hex128

```go
func (g *Generator) Hex128() string
```
Hex128 is a convenience method that returns Hex128(g.Next()).

#### func (*Generator) Next

```go
func (g *Generator) Next() [24]byte
```
Next returns the next UUID from the generator. Only the first 8 bytes can differ
from the previous UUID, so taking a slice of the first 16 bytes is sufficient to
provide a somewhat less secure 128 bit UUID.

It is OK to call this method concurrently.
