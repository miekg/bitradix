package bitradix

import (
	"testing"
)

func TestInsert(t *testing.T) {
	tests := map[uint8]uint32{
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

func TestFind(t *testing.T) {
	r := New()
	r.Insert(0x08, 2012)
	r.Insert(0x04, 2010)
	r.Insert(0x09, 2013)
	println(r.String())

	v1 := r.Find(0x08)
	println(v1.Key, v1.Value)
	v1 = r.Find(0x04)
	println(v1.Key, v1.Value)
}

type bittest struct {
	value uint8
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
