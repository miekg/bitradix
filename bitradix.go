// Package bitradix implements a radix tree that branches on the bits in a 32 bits key.
// The value that can be stored is an unsigned 32 bit integer.
//                                                                                                  
// A radix tree is defined in:
//    Donald R. Morrison. "PATRICIA -- practical algorithm to retrieve
//    information coded in alphanumeric". Journal of the ACM, 15(4):514-534,
//    October 1968
package bitradix

import (
	"fmt"
)

// With help from:
// http://faculty.simpson.edu/lydia.sinapova/www/cmsc250/LN250_Weiss/L08-Radix.htm

const bitSize = 32 // length in bits of the key

// Radix implements a radix tree.
type Radix struct {
	branch   [2]*Radix // branch[0] is left branch for 0, and branch[1] the right for 1
	key      uint32    // the key under which this value is stored
	bits     int       // the number of significant bits, if 0 the key has not been set.
	Value    uint32    // The value stored.
	internal bool      // internal node
}

// New returns an empty, initialized Radix tree.
func New() *Radix {
	return &Radix{[2]*Radix{nil, nil}, 0, 0, 0, false}
}

// Key returns the key under which this node is stored.
func (r *Radix) Key() uint32 {
	return r.key
}

// Bits returns the number of significant bits for the key.
// A value of zero indicates a key that has not been set.
func (r *Radix) Bits() int {
	return r.bits
}

// Internal returns true is r is an internal node, when false is returned
// the node is a leaf node.
func (r *Radix) Internal() bool {
	return r.internal
}

// Insert inserts a new value n in the tree r. The first bits bits of n are significant
// and used to store the value v.
// It returns the inserted node, r must be the root of the tree.
func (r *Radix) Insert(n uint32, bits int, v uint32) *Radix {
	println("INSERTING", n, v)
	return r.insert(n, bits, v, bitSize-1)
}

// Remove removes a value from the tree r. It returns the node removed, or nil
// when nothing is found. r must be the root of the tree.
func (r *Radix) Remove(n uint32, bits int) *Radix {
	return nil
}

// Find searches the tree for the key n, where the first bits bits of n 
// are significant. It returns the node found.
func (r *Radix) Find(n uint32, bits int) *Radix {
	return r.find(n, bits, bitSize-1)
}

// Do traverses the tree r in breadth-first order. For each visited node,
// the function f is called with the current node and the branch taken
// (0 for the zero, 1 for the one branch, -1 is used for the root node).
func (r *Radix) Do(f func(*Radix, int)) {
	q := make(queue, 0)

	q.Push(&node{r, -1})
	x := q.Pop()
	for x != nil {
		f(x.Radix, x.branch)
		for i, b := range x.Radix.branch {
			if b != nil {
				println("NODE", i)
				q.Push(&node{b, i})
			}
		}
		x = q.Pop()
	}
}

func (r *Radix) String() string {
	return r.str(" ")
}

func (r *Radix) str(j string) string {
	if r == nil {
		return ""
	}
	s := fmt.Sprintf("%s %032b -> %d\n", j, r.key, r.Value)
	s += r.branch[0].str(j + "")
	s += r.branch[1].str(j + "")
	return s
}

// Implement insert
func (r *Radix) insert(n uint32, bits int, v uint32, bit int) *Radix {
	switch r.internal {
	case true:
		if bitSize-bits == bit { // we need to store a value here
			println("Store here internal")
			// TODO(mg): check previous value?
			r.key = n
			r.bits = bits
			r.Value = v
			// keep it internal
			return r
		}
		// Internal node, no key. With branches, walk the branches.
		return r.branch[bitK(n, bit)].insert(n, bits, v, bit-1)
	case false:
		println("external")
		// External node, (optional) key, no branches
		if r.bits == 0 { // nothing here yet, put something in
			println("nothing here yet, put something in ", n, bits)
			r.bits = bits
			r.key = n
			r.Value = v
			return r
		}
		if bitSize-bits == bit { // seen all bits, put something here
			println("seen all bits, put something here")
			if r.bits != 0 {
				println("something here ALREADY")
			}
			r.bits = bits
			r.key = n
			r.Value = v
			return r
		}

		// create new branches, and go from there
		r.branch[0], r.branch[1] = New(), New()
		r.internal = true // becomes an internal node by definition

		bcur := bitK(r.key, bit)
		bnew := bitK(n, bit)
		fmt.Printf("r.key %032b %d\n", r.key, bit)
		fmt.Printf("n     %032b %d\n", n, bit)

		println("bcur", bcur, "bnew", bnew)

		switch x := bitSize - r.bits; true {
		case x == bit: // current node needs to stay here
			println("Current node needs to be kept here")
			// put new stuff in the branch below
			r.branch[bnew].key = n
			r.branch[bnew].Value = v
			r.branch[bnew].bits = bits
			return r.branch[bnew]
		case x < bit: // current node can be put one level down
			println("Moving nodes down", r.key, r.bits, r.Value)
			if bcur == bnew {
				println("equal bits")
				// "fill" the correct node, with the current key - and call ourselves
				r.branch[bcur].key = r.key
				r.branch[bcur].Value = r.Value
				r.branch[bcur].bits = r.bits
				r.bits = 0
				r.key = 0
				r.Value = 0
				return r.branch[bnew].insert(n, bits, v, bit-1)
			}
			// bcur = 0, bnew == 1 or vice versa
			println("not equal, each in own branch")
			r.branch[bcur].key = r.key
			r.branch[bcur].Value = r.Value
			r.branch[bcur].bits = r.bits
			r.branch[bnew].key = n
			r.branch[bnew].Value = v
			r.branch[bnew].bits = bits
			r.key = 0
			r.Value = 0
			r.bits = 0
			return r.branch[bnew]
		case x > bit: // node is at the wrong spot
			panic("bitradix: node put too far down")
		}

	}
	panic("bitradix: not reached")
}

func (r *Radix) find(n uint32, bits, bit int) *Radix {
	switch r.internal {
	case true:
		// Internal node, no key, continue in the right branch
		return r.branch[bitK(n, bit)].find(n, bits, bit-1)
	case false:
		return r
	}
	panic("bitradix: not reached")
}

// From: http://stackoverflow.com/questions/2249731/how-to-get-bit-by-bit-data-from-a-integer-value-in-c

// Return bit k from n. We count from the right, MSB left.
// So k = 0 is the last bit on the left and k = 63 is the first bit on the right.
func bitK(n uint32, k int) byte {
	return byte((n & (1 << uint(k))) >> uint(k))
}
