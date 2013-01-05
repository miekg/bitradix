// Package bitradix implements a radix tree that branches on the bits in a 64 bits key.
// The value that can be stored is an unsigned 32 bit integer.
//                                                                                                  
// A radix tree is defined in:
//    Donald R. Morrison. "PATRICIA -- practical algorithm to retrieve
//    information coded in alphanumeric". Journal of the ACM, 15(4):514-534,
//    October 1968
package bitradix

// With help from:
// http://faculty.simpson.edu/lydia.sinapova/www/cmsc250/LN250_Weiss/L08-Radix.htm

import (
	"strconv"
)

const bitSize = 64 // length in bits of the key

// Radix implements a radix tree.
type Radix struct {
	branch   [2]*Radix // branch[0] is left branch for 0, and branch[1] the right for 1
	key      uint64    // The key under which this value is stored.
	set      bool      // true if the key has been set
	Value    uint32    // The value stored.
	internal bool      // internal node
}

// New returns an empty, initialized Radix tree.
func New() *Radix {
	return &Radix{[2]*Radix{nil, nil}, 0, false, 0, false}
}

// Key returns the key under which this node is stored.
func (r *Radix) Key() uint64 {
	return r.key
}

// Set returns if the key has been set for this node. If set is false
// the meaning of the key is undefined.
func (r *Radix) Set() bool {
	return r.set
}

// Internal returns true is r is an internal node, when false is returned
// the node is a leaf node.
func (r *Radix) Internal() bool {
	return r.internal
}

// Insert inserts a new value in the tree r. It returns the inserted node.
// r must be the root of the tree.
func (r *Radix) Insert(n uint64, v uint32) *Radix {
	return r.insert(n, v, bitSize-1)
}

// Remove removes a value from the tree r. It returns the node removed, or nil
// when nothing is found. r must be the root of the tree.
func (r *Radix) Remove(n uint64) *Radix {
	return nil
}

// Find searches the tree for the key n. It returns the node found,
// and the number of branches taken. The later is the longest common
// prefix.
func (r *Radix) Find(n uint64) (*Radix, int) {
	return r.find(n, bitSize-1)
}

// Do traverses the tree r in depth-first order. For each visited node,
// the function f is called.
func (r *Radix) Do(f func(*Radix)) {
	f(r)
	for _, branch := range r.branch {
		branch.Do(f)
	}
}

// Implement insert
func (r *Radix) insert(n uint64, v uint32, bit uint) *Radix {
	switch r.internal {
	case true:
		// Internal node, no key. With branches, walk the branches.
		return r.branch[bitK(n, bit)].insert(n, v, bit-1)
	case false:
		// External node, (optional) key, no branches
		if !r.set {
			r.set = true
			r.key = n
			r.Value = v
			return r
		}

		// create new branches, and go from there
		r.branch[0], r.branch[1] = New(), New()
		// Current node, becomes an intermediate node
		r.internal = true
		r.set = false

		bcur := bitK(r.key, bit)
		bnew := bitK(n, bit)
		if bcur == bnew {
			// "fill" the correct node, with the current key - and call ourselves
			r.branch[bcur].key = r.key
			r.branch[bcur].Value = r.Value
			r.branch[bcur].set = true
			r.key = 0
			r.Value = 0
			if bit == 0 {
				r.branch[bnew].key = n
				r.branch[bnew].Value = v
				r.branch[bnew].set = true
				return r.branch[bnew]
			}
			return r.branch[bnew].insert(n, v, bit-1)
		}
		// bcur = 0, bnew == 1 or vice versa
		r.branch[bcur].key = r.key
		r.branch[bcur].Value = r.Value
		r.branch[bcur].set = true
		r.branch[bnew].key = n
		r.branch[bnew].Value = v
		r.branch[bnew].set = true
		r.key = 0
		r.Value = 0
		return r.branch[bnew]
	}
	panic("bitradix: not reached")
}

func (r *Radix) find(n uint64, bit uint) (*Radix, int) {
	switch r.internal {
	case true:
		// Internal node, no key, continue in the right branch
		return r.branch[bitK(n, bit)].find(n, bit-1)
	case false:
		return r, int(bitSize - bit)
	}
	panic("bitradix: not reached")
}

// From: http://stackoverflow.com/questions/2249731/how-to-get-bit-by-bit-data-from-a-integer-value-in-c

// Return bit k from n. We count from the right, MSB left.
// So k = 0 is the last bit on the left and k = 63 is the first bit on the right.
func bitK(n uint64, k uint) byte {
	return byte((n & (1 << k)) >> k)
}
