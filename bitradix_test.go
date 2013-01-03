package bitradix

import (
	"testing"
)

func TestInsert(t *testing.T) {
	r := New()
	r.Insert(0x08, 2012)
	r.Insert(0x04, 2010)
//	r.Insert(0x09, 2013)
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
