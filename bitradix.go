// Package bitradix implements a radix tree that branches on the bits in a 32 bits key.
// The value that can be stored is an unsigned 32 bit integer.
//                                                                                                  
// A radix tree is defined in:
//    Donald R. Morrison. "PATRICIA -- practical algorithm to retrieve
//    information coded in alphanumeric". Journal of the ACM, 15(4):514-534,
//    October 1968
package bitradix

// With help from:
// http://faculty.simpson.edu/lydia.sinapova/www/cmsc250/LN250_Weiss/L08-Radix.htm

const (
	bitSize32 = 32
	bitSize64 = 64
)

// Radix32 implements a radix tree with an uint32 as its key.
type Radix32 struct {
	branch   [2]*Radix32 // branch[0] is left branch for 0, and branch[1] the right for 1
	key      uint32    // the key under which this value is stored
	bits     int       // the number of significant bits, if 0 the key has not been set.
	Value    uint32    // The value stored.
	internal bool      // internal node
}

// Radix64 implements a radix tree with an uint64 as its key.
type Radix64 struct {
	branch   [2]*Radix64 // branch[0] is left branch for 0, and branch[1] the right for 1
	key      uint64    // the key under which this value is stored
	bits     int       // the number of significant bits, if 0 the key has not been set.
	Value    uint32    // The value stored.
	internal bool      // internal node
}

// New32 returns an empty, initialized Radix32 tree.
func New32() *Radix32 {
	return &Radix32{[2]*Radix32{nil, nil}, 0, 0, 0, false}
}

// New64 returns an empty, initialized Radix64 tree.
func New64() *Radix64 {
	return &Radix64{[2]*Radix64{nil, nil}, 0, 0, 0, false}
}

// Key returns the key under which this node is stored.
func (r *Radix32) Key() uint32 {
	return r.key
}

// Key returns the key under which this node is stored.
func (r *Radix64) Key() uint64 {
	return r.key
}

// Bits returns the number of significant bits for the key.
// A value of zero indicates a key that has not been set.
func (r *Radix32) Bits() int {
	return r.bits
}

// Bits returns the number of significant bits for the key.
// A value of zero indicates a key that has not been set.
func (r *Radix64) Bits() int {
	return r.bits
}

// Internal returns true is r is an internal node, when false is returned
// the node is a leaf node.
func (r *Radix32) Internal() bool {
	return r.internal
}

// Internal returns true is r is an internal node, when false is returned
// the node is a leaf node.
func (r *Radix64) Internal() bool {
	return r.internal
}

// Insert inserts a new value n in the tree r. The first bits bits of n are significant
// and used to store the value v.
// It returns the inserted node, r must be the root of the tree.
func (r *Radix32) Insert(n uint32, bits int, v uint32) *Radix32 {
	return r.insert(n, bits, v, bitSize32-1)
}

// Remove removes a value from the tree r. It returns the node removed, or nil
// when nothing is found. r must be the root of the tree.
func (r *Radix32) Remove(n uint32, bits int) *Radix32 {
	return nil
}

// Find searches the tree for the key n, where the first bits bits of n 
// are significant. It returns the node found.
func (r *Radix32) Find(n uint32, bits int) *Radix32 {
	return r.find(n, bits, bitSize32-1, nil)
}

// Do traverses the tree r in breadth-first order. For each visited node,
// the function f is called with the current node and the branch taken
// (0 for the zero, 1 for the one branch, -1 is used for the root node).
func (r *Radix32) Do(f func(*Radix32, int)) {
	q := make(queue32, 0)

	q.Push(&node32{r, -1})
	x := q.Pop()
	for x != nil {
		f(x.Radix32, x.branch)
		for i, b := range x.Radix32.branch {
			if b != nil {
				q.Push(&node32{b, i})
			}
		}
		x = q.Pop()
	}
}

// Implement insert
func (r *Radix32) insert(n uint32, bits int, v uint32, bit int) *Radix32 {
	switch r.internal {
	case true:
		if bitSize32-bits == bit { // we need to store a value here
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
		// External node, (optional) key, no branches
		if r.bits == 0 { // nothing here yet, put something in
			r.bits = bits
			r.key = n
			r.Value = v
			return r
		}
		if bitSize32-bits == bit { // seen all bits, put something here
			if r.bits != 0 {
				// println("something here ALREADY")
			}
			r.bits = bits
			r.key = n
			r.Value = v
			return r
		}

		// create new branches, and go from there
		r.branch[0], r.branch[1] = New32(), New32()
		r.internal = true // becomes an internal node by definition

		bcur := bitK(r.key, bit)
		bnew := bitK(n, bit)
		// fmt.Printf("r.key %032b %d\n", r.key, bit)
		// fmt.Printf("n     %032b %d\n", n, bit)

		switch x := bitSize32 - r.bits; true {
		case x == bit: // current node needs to stay here
			// put new stuff in the branch below
			r.branch[bnew].key = n
			r.branch[bnew].Value = v
			r.branch[bnew].bits = bits
			return r.branch[bnew]
		case x < bit: // current node can be put one level down
			if bcur == bnew {
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

// Search the tree, when "seeing" a node with a key, store that
// node, when we don't find anything within the allowed 
func (r *Radix32) find(n uint32, bits, bit int, last *Radix32) *Radix32 {
	switch r.internal {
	case true:
		if r.bits != 0 {
			// Actual key, drag it along
			return r.branch[bitK(n, bit)].find(n, bits, bit-1, r)
		}
		return r.branch[bitK(n, bit)].find(n, bits, bit-1, last)
	case false:
		mask := uint32(0xFFFFFFFF << uint(r.bits))
		if r.key&mask == n&mask {
			return r
		}
		return last
	}
	panic("bitradix: not reached")
}

// From: http://stackoverflow.com/questions/2249731/how-to-get-bit-by-bit-data-from-a-integer-value-in-c

// Return bit k from n. We count from the right, MSB left.
// So k = 0 is the last bit on the left and k = 63 is the first bit on the right.
func bitK(n uint32, k int) byte {
	return byte((n & (1 << uint(k))) >> uint(k))
}
