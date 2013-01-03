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

// Radix implements a radix tree. 
type Radix struct {
	branch   [2]*Radix // branch[0] is left branch for 0, and branch[1] the right for 1
	Key      uint64    // The key under which this value is stored.
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

// Remove removes a value from the tree r. It returns the node removed, or nil
// when nothing is found. r must be the root of the tree.
func (r *Radix) Remove(n uint64) *Radix {
	return nil
}

// Find searches the tree for the key n. It returns the node found. Note that
// the returned node does not have to be an exact match, it may also be the best
// matching parent. Comparing r.Key with n tells the caller if the match was exact.
func (r *Radix) Find(n uint64) *Radix {
	return r.find(n, 63)
}

func (r *Radix) String() string {
	return r.str("")
}

// Implement insert
func (r *Radix) insert(n uint64, v uint32, bit uint) *Radix {
	// if bit == 0 ? TODO(mg) When does that happen
	switch r.internal {
	case true:
		// Internal node, no key. With branches, walk the branches.
		return r.branch[bitK(n, bit)].insert(n, v, bit-1)
	case false:
		// External node, (optional) key, no branches
		if !r.keyset {
			r.keyset = true
			r.Key = n
			r.Value = v
			return r
		}

		// match keys and create new branches, and go from there
		r.branch[0], r.branch[1] = New(), New()

		r.internal = true
		r.keyset = false
		bcur := bitK(r.Key, bit)
		bnew := bitK(n, bit)

		if bcur == bnew {
			// "fill" the correct node, with the current key - and call ourselves
			r.branch[bcur].Key = r.Key
			r.branch[bcur].keyset = true
			r.Key = 0
			return r.branch[bcur].insert(n, v, bit-1)
		}
		// bcur = 0, bnew == 1 or vice versa
		r.branch[bcur].Key = r.Key
		r.branch[bcur].keyset = true
		r.branch[bnew].Key = n
		r.branch[bnew].keyset = true
		r.Key = 0
		return r.branch[bnew]
	}
	panic("bitradix: not reached")
}

func (r *Radix) find(n uint64, bit uint) *Radix {
	// if bit == 0, return the current node?? Also see comment in r.insert()
	switch r.internal {
	case true:
		// Internal node, no key
		return r.branch[bitK(n, bit)].find(n, bit-1)
	case false:
		// External node, (optional) key, no branches
		return r
	}
	panic("bitradix: not reached")
}

func (r *Radix) str(indent string) (s string) {
	s += indent
	if r.keyset {
		s += strconv.Itoa(int(r.Key)) + "\n" + indent
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
