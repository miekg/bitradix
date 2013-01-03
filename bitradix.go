// Package bitradix implements a radix tree that branches on the bits in a 64 bits key.
// The value that can be stored is an unsigned 32 bit integer.
//                                                                                                  
// A radix tree is defined in:
//    Donald R. Morrison. "PATRICIA -- practical algorithm to retrieve
//    information coded in alphanumeric". Journal of the ACM, 15(4):514-534,
//    October 1968
package bitradix

import (
	"strconv"
)

// Radix implements a radix tree. 
type Radix struct {
	branch   [2]*Radix // branch[0] is left branch for 0, and branch[1] the right for 1
	key      uint64    // the key stored
	keyset   bool      // true if the key has been set
	Value    uint32    // The value stored.
	internal bool      // internal node
}

// New returns an empty, initialized Radix tree.
func New() *Radix {
	return &Radix{[2]*Radix{nil, nil}, 0, false, 0, false}
}

// Insert inserts a new value in the tree r. It returns the inserted node.
// r must be the root of the tree.
func (r *Radix) Insert(n uint64, v uint32) *Radix {
	return r.insert(n, v, 63)
}

// Implement insert
func (r *Radix) insert(n uint64, v uint32, bit uint) *Radix {
	// if bit == 0 ? TODO(mg) When does that happen
	switch r.internal {
	case true:
		// Internal node, no key. With branches, walk the branches.
		return r.branch[bitK(n, bit)].insert(n, v, bit-1)
	case false:
		// external node, (optional) key, no branches
		if !r.keyset {
			r.keyset = true
			r.key = n
			r.Value = v
			return r
		}

		// match keys and create new branches, and go from there
		r.branch[0], r.branch[1] = New(), New()

		r.internal = true
		r.keyset = false
		bitcurrent := bitK(r.key, bit)
		bitnew := bitK(n, bit)

		if bitcurrent == bitnew {
			// "fill" the correct node, with the current key - and call ourselves
			r.branch[bitcurrent].key = r.key
			r.branch[bitcurrent].keyset = true
			r.key = 0
			return r.branch[bitcurrent].insert(n, v, bit-1)
		}
		// bitcurrent = 0, bitnew == 1 or vice versa
		r.branch[bitcurrent].key = r.key
		r.branch[bitcurrent].keyset = true
		r.branch[bitnew].key = n
		r.branch[bitnew].keyset = true
		r.key = 0
		return r.branch[bitnew]
	}
	panic("bitradix: not reached")
}

func (r *Radix) String() string {
	return r.str("")
}

func (r *Radix) str(indent string) (s string) {
	s += indent
	if r.keyset {
		s += strconv.Itoa(int(r.key)) + "\n" + indent
	} else {
		s += "<nil>\n" + indent
	}
	if r.internal {
		s += "0: " + r.branch[0].str(indent+" ")
		s += "1: " + r.branch[1].str(indent+" ")
	}
	return s
}

// From: http://stackoverflow.com/questions/2249731/how-to-get-bit-by-bit-data-from-a-integer-value-in-c

// Return bit k from n. We count from the right, MSB left.
// So k = 0 is the last bit on the left and k = 63 is the first bit on the right.
func bitK(n uint64, k uint) byte {
	return byte((n & (1 << k)) >> k)
}
