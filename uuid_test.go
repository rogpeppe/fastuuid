package fastuuid

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestUUID(t *testing.T) {
	// TODO fix this test now that we can return non-deterministic
	// counters.
	var buf [24]byte
	for i := range buf {
		buf[i] = byte(i) + 1
	}
	oldReader := rand.Reader
	rand.Reader = bytes.NewReader(buf[:])
	g, err := NewGenerator()
	rand.Reader = oldReader
	if err != nil {
		t.Fatalf("cannot make generator: %v", err)
	}
	uuid := g.Next()
	buf[0] = 1 ^ 1
	if uuid != buf {
		t.Fatalf("unexpected UUID; got %x; want %x", uuid, buf)
	}
	uuid = g.Next()
	buf[0] = 1 ^ 2
	if uuid != buf {
		t.Fatalf("unexpected next UUID; got %x; want %x", uuid, buf)
	}
}

func TestUniqueness(t *testing.T) {
	g := MustNewGenerator()
	mc := make(chan map[[24]byte]int)
	const nproc = 4
	for i := 0; i < nproc; i++ {
		go func() {
			m := make(map[[24]byte]int)
			for i := 0; i < step*10; i++ {
				uuid := g.Next()
				if old, ok := m[uuid]; ok {
					t.Errorf("non-unique uuid seq at %d, other %d", i, old)
				}
				m[uuid] = i
			}
			mc <- m
		}()
	}
	m := make(map[[24]byte]int)
	for i := 0; i < nproc; i++ {
		for uuid, iter := range <-mc {
			if old, ok := m[uuid]; ok {
				t.Errorf("non-unique uuid seq at %d, other %d", i, old)
			}
			m[uuid] = iter
		}
	}
}

func BenchmarkNext(b *testing.B) {
	g := MustNewGenerator()
	for i := 0; i < b.N; i++ {
		g.Next()
	}
}

func BenchmarkContended(b *testing.B) {
	g := MustNewGenerator()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			g.Next()
		}
	})
}
