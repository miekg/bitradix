package bitradix

import (
	"net"
	"reflect"
	"testing"
)

// Test multiple insertions under the same key

func TestInsert(t *testing.T) {
	tests := map[uint64]uint32{
		0x08: 2012,
		0x04: 2010,
		0x09: 2013,
	}
	r := New()
	for key, value := range tests {
		if x := r.Insert(key, value); x.Value != value {
			t.Logf("Expected %d, got %d for %d (node type %v)\n", value, x.Value, key, x.Internal())
			t.Fail()
		}
	}
}

func TestInsertIdempotent(t *testing.T) {
	r := New()
	r.Insert(0x08, 2012)
	r.Insert(0x08, 2013)
	if x, _ := r.Find(0x08); x.Value != 2013 {
		t.Logf("Expected %d, got %d for %d\n", 2013, x.Value, 0x08)
		t.Fail()
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
		if x, _ := r.Find(k); x.Value != v {
			t.Logf("Expected %d, got %d for %d (node type %v)\n", v, x.Value, k, x.Internal())
			t.Fail()
		}
	}
}

// Test with "real-life" ip addresses
func ipToUint(t *testing.T, ip net.IP) (i uint64) {
	ip = ip.To4()
	fv := reflect.ValueOf(&i).Elem()
	fv.SetUint(uint64(uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[+3])))
	t.Logf("Bit %08b\n", i)
	return
}

func addRoute(t *testing.T, r *Radix, s string, asn uint32) {
	t.Logf("Route %s, AS %d\n", s, asn)
	_, ipnet, _ := net.ParseCIDR(s)
	r.Insert(ipToUint(t, ipnet.IP), asn)
}

func findRoute(t *testing.T, r *Radix, s string) uint32 {
	_, ipnet, _ := net.ParseCIDR(s)
	t.Logf("Search %s\n", s)
	node, _ := r.Find(ipToUint(t, ipnet.IP)) // discard step
	//	if node == nil {
	//		return 0
	//	}
	return node.Value
}

func TestFindIP(t *testing.T) {
	r := New()
	// not a map to have influence on the order
	addRoute(t, r, "10.0.0.2/8", 10)
	addRoute(t, r, "10.20.0.0/14", 20)
	addRoute(t, r, "10.21.0.0/16", 21)
	addRoute(t, r, "192.168.0.0/16", 192)
	addRoute(t, r, "192.168.2.0/24", 1922)

	testips := map[string]uint32{
		"10.20.1.2/32":   20,
		"10.22.1.2/32":   20,
		"10.19.0.1/32":   10,
		"10.21.0.1/32":   21,
		"192.168.2.3/32": 1922,
		"230.0.0.1/32":   0,
	}

	for ip, asn := range testips {
		if x := findRoute(t, r, ip); asn != x {
			t.Logf("Expected %d, got %d for %s\n", asn, x, ip)
			t.Fail()
		}
	}
}

func TestFindIPShort(t *testing.T) {
	r := New()
	// not a map to have influence on the order
	addRoute(t, r, "10.0.0.2/8", 10)
	addRoute(t, r, "10.20.0.0/14", 20)

	t.Logf("%s\n", r.string())

	testips := map[string]uint32{
		"10.20.1.2/32": 20,
		"10.19.0.1/32": 10,
	}

	for ip, asn := range testips {
		if x := findRoute(t, r, ip); asn != x {
			t.Logf("Expected %d, got %d for %s\n", asn, x, ip)
			t.Fail()
		}
	}
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
