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
	zero     *Radix // left branch
	one      *Radix // right branch
	key      uint64 // the key stored
	keyset   bool   // true if the key has been set
	Value    uint32 // The value stored.
	internal bool   // internal node
}

// New returns an empty, initialized Radix tree.
func New() *Radix {
	return &Radix{nil, nil, 0, false, 0, false}
}

// Insert inserts a new value in the tree r. It returns the inserted node.
// r must be the root of the tree.
func (r *Radix) Insert(n uint64, v uint32) *Radix {
	return r.insert(n, v, 63)
}

// Implement insert
func (r *Radix) insert(n uint64, v uint32, bit uint) *Radix {
	// if bit == 1 ? TODO(mg)
	switch r.internal {
	case true:
		// internal node, no key, with branches, walk the branches
		switch bitK(n, bit) {
		case 0:
			return r.zero.insert(n, v, bit-1)
		case 1:
			return r.one.insert(n, v, bit-1)
		}
	case false:
		// external node, possible key, no branches
		if !r.keyset {
			r.keyset = true
			r.key = n
			r.Value = v
			return r
		}

		// match keys and create new branches, and go from there

		// branch
		r.zero = &Radix{nil, nil, 0, false, 0, false}
		r.one = &Radix{nil, nil, 0, false, 0, false}
		r.internal = true
		r.keyset = false
		bitcurrent := bitK(r.key, bit)
		bitnew := bitK(n, bit)

		if bitcurrent == bitnew {
			// "fill" the correct node, with the current key - and reenter the function
			if bitcurrent == 0 {
				r.zero.key = r.key
				r.zero.keyset = true
				r.key = 0
				return r.zero.insert(n, v, bit-1)
			}
			if bitcurrent == 1 {
				r.one.key = r.key
				r.one.keyset = true
				r.key = 0
				return r.one.insert(n, v, bit-1)
			}
		} else {
			// not equal, branch and fill and return	
			if bitcurrent == 0 {
				// bitnew == 1
				r.zero.key = r.key
				r.zero.keyset = true
				r.one.key = n
				r.one.keyset = true
				r.key = 0
				return r.one
			}
			if bitcurrent == 1 {
				// bitnew == 0
				r.one.key = r.key
				r.one.keyset = true
				r.zero.key = n
				r.zero.keyset = true
				r.key = 0
				return r.zero
			}
		}
	}
	return nil
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
		s += "0: " + r.zero.str(indent + " ")
		s += "1: " + r.one.str(indent + " ")
	}
	return s
}

// From: http://stackoverflow.com/questions/2249731/how-to-get-bit-by-bit-data-from-a-integer-value-in-c

// Return bit k from n. We count from the right, MSB left.
// So k = 0 is the last bit on the left and k = 63 is the first bit on the right.
func bitK(n uint64, k uint) byte {
	return byte((n & (1 << k)) >> k)
}
