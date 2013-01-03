// Package bitradix implements a radix tree that branches on the bits in a key.
// The value that can be stored is an unsigned 32 bit integer
//                                                                                                  
// A radix tree is defined in:
//    Donald R. Morrison. "PATRICIA -- practical algorithm to retrieve
//    information coded in alphanumeric". Journal of the ACM, 15(4):514-534,
//    October 1968
package bitradix

import (
	"fmt"
)

// Radix implements a radix tree. 
type Radix struct {
	zero     *Radix // left branch
	one      *Radix // right branch
	key      uint64 // the key stored
	keyset   bool   // true if the key has been set
	Value    string // The value stored.
	internal bool   // internal node
}

// New returns an empty, initialized Radix tree.
func New() *Radix {
	return &Radix{nil, nil, 0, false, "", false}
}

// Insert inserts a new value in the tree r. It returns the inserted node.
// r must be the root of the tree.
func (r *Radix) Insert(n uint64, v string) *Radix {
	return r.insert(n, v, 63)
}

// Implement insert
func (r *Radix) insert(n uint64, v string, bit uint) *Radix {
	fmt.Printf("key %064b\n", n)
	// if bit == 1 ? TODO(mg)
	switch r.internal {
	case true:
		// internal node, no key, with branches
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
		// match keys and see how many extra nodes we need to store these



	}
	return nil
}

// From: http://stackoverflow.com/questions/2249731/how-to-get-bit-by-bit-data-from-a-integer-value-in-c

// Return bit k from n. We count from the right, MSB left.
// So k = 0 is the last bit on the left and k = 63 is the first bit on the right.
func bitK(n uint64, k uint) byte {
	return byte((n & (1 << k)) >> k)
}
