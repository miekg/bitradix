package bitradix

import (
	"testing"
)

func TestInsert(t *testing.T) {
	tests := map[uint64]uint32{
		0x08: 2012,
		0x04: 2010,
		0x09: 2013,
	}
	r := New()
	for key, value := range tests {
		if x := r.Insert(key, value); x.Value != value {
			t.Logf("Expected %d, got %d for %d\n", value, x.Value, key)
			t.Fail()
		}
	}
}

func TestFindExact(t *testing.T) {
	tests := map[uint64]uint32{
		0x08: 2012,
		0x04: 2010,
		0x09: 2013,
	}
	r := New()
	for k, v := range tests {
		r.Insert(k, v)
	}
	for k, v := range tests {
		if x := r.Find(k); x.Value != v {
			t.Logf("Expected %d, got %d for %d\n", v, x.Value, k)
			t.Fail()
		}
	}
}

func TestFind(t *testing.T) {
	r := New()
	r.Insert(0x08, 2001)	// This is a /n address 00...001000
	r.Insert(0x09, 2001)	// This is also a /n    00...001001

	// Longest common prefix
	x := r.Find(0xa)  // Look for /n 00..001010
	println("key", x.key, "value", x.Value)
}

type bittest struct {
	value uint64
	bit   uint
}

func TestBitK(t *testing.T) {
	tests := map[bittest]byte{
		bittest{0x40, 0}: 0,
		bittest{0x40, 6}: 1,
	}
	for test, expected := range tests {
		if x := bitK(test.value, test.bit); x != expected {
			t.Logf("Expected %d for %064b (bit #%d), got %d\n", expected, test.value, test.bit, x)
			t.Fail()
		}
	}
}
