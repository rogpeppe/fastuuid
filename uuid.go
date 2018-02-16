// Package fastuuid provides fast UUID generation of 192 bit universally
// unique identifiers. It does not provide formatting or parsing of the
// identifiers (it is assumed that a simple hexadecimal or base64
// representation is sufficient, for which adequate functionality exists
// elsewhere).
//
// Note that the generated UUIDs are not unguessable - each UUID
// generated from a Generator is adjacent to the previously generated
// UUID.
//
// It ignores RFC 4122.
package fastuuid

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"sync"
)

// Generator represents a UUID generator that
// generates UUIDs in sequence from a random starting
// point.
type Generator struct {
	pool  sync.Pool
	mu    sync.Mutex
	base0 uint64
	base1 uint64
	seed  [24]byte
}

type counter struct {
	n     uint16
	limit uint16
	seed  [24]byte
}

const step = 32768

// NewGenerator returns a new Generator.
// It can fail if the crypto/rand read fails.
func NewGenerator() (*Generator, error) {
	var g Generator
	_, err := rand.Read(g.seed[:])
	if err != nil {
		return nil, errors.New("cannot generate random seed: " + err.Error())
	}
	g.pool.New = g.newCounter
	return &g, nil
}

func (g *Generator) newCounter() interface{} {
	g.mu.Lock()
	defer g.mu.Unlock()
	limit := g.base0 + step
	g.base0 = limit
	if limit < step {
		// Overflow.
		g.base1++
	}
	n := limit - step
	c := &counter{
		seed:  g.seed,
		n:     uint16(n & 0xffff),
		limit: uint16(limit & 0xffff),
	}
	xorBytesUint64(c.seed[:], n)
	// Reset the first two bytes. They will be filled in by Next.
	c.seed[0] = g.seed[0]
	c.seed[1] = g.seed[1]
	xorBytesUint64(c.seed[8:], g.base1)
	return c
}

// MustNewGenerator is like NewGenerator
// but panics on failure.
func MustNewGenerator() *Generator {
	g, err := NewGenerator()
	if err != nil {
		panic(err)
	}
	return g
}

// Next returns the next UUID from the generator. Only the first 8 bytes
// can differ from the previous UUID, so taking a slice of the first 16
// bytes is sufficient to provide a somewhat less secure 128 bit UUID.
//
// It is OK to call this method concurrently.
func (g *Generator) Next() (uuid [24]byte) {
	for {
		c := g.pool.Get().(*counter)
		n := c.n + 1
		if n == c.limit {
			// This counter has been exhausted. Abandon
			// it and get another one from the pool.
			continue
		}
		c.n = n
		uuid = c.seed
		g.pool.Put(c)
		uuid[0] ^= byte(n)
		uuid[1] ^= byte(n >> 8)
		return uuid
	}
}

func xorBytesUint64(into []byte, n uint64) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], n)
	// XOR all but the first two bytes. The first two bytes will be filled in by Next.
	for i, b := range buf[:] {
		into[i] ^= b
	}
}
