// Package fastuuid provides fast UUID generation of 192 bit
// universally unique identifiers. It does not provide
// formatting or parsing of the identifiers (it is assumed
// that a simple hexadecimal or base64 representation
// is sufficient, for which adequate functionality exists elsewhere).
//
// Note that the generated UUIDs are not unguessable - each
// UUID generated from a Generator is adjacent to the
// previously generated UUID.
//
// It ignores RFC 4122.
package fastuuid

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"sync/atomic"
)

// Generator represents a UUID generator that
// generates UUIDs in sequence from a random starting
// point.
type Generator struct {
	// The constant seed. The first 8 bytes of this are
	// copied into counter and then ignored thereafter.
	seed    [24]byte
	counter uint64
}

// NewGenerator returns a new Generator.
// It can fail if the crypto/rand read fails.
func NewGenerator() (*Generator, error) {
	var g Generator
	_, err := rand.Read(g.seed[:])
	if err != nil {
		return nil, errors.New("cannot generate random seed: " + err.Error())
	}
	g.counter = binary.LittleEndian.Uint64(g.seed[:8])
	return &g, nil
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

// Next returns the next UUID from the generator.
// Only the first 8 bytes can differ from the previous
// UUID, so taking a slice of the first 16 bytes
// is sufficient to provide a somewhat less secure 128 bit UUID.
//
// It is OK to call this method concurrently.
func (g *Generator) Next() [24]byte {
	x := atomic.AddUint64(&g.counter, 1)
	uuid := g.seed
	binary.LittleEndian.PutUint64(uuid[:8], x)
	return uuid
}
