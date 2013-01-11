// Package bitradix implements a radix tree that branches on the bits of a 32 or
// 64 bits unsigned integer key.
//                                                                                                  
// A radix tree is defined in:
//    Donald R. Morrison. "PATRICIA -- practical algorithm to retrieve
//    information coded in alphanumeric". Journal of the ACM, 15(4):514-534,
//    October 1968
// 
// This website provides some background information on Radix trees.
// http://faculty.simpson.edu/lydia.sinapova/www/cmsc250/LN250_Weiss/L08-Radix.htm
package bitradix

const (
	bitSize32 = 32
	bitSize64 = 64
	mask32    = 0xFFFFFFFF
	mask64    = 0xFFFFFFFFFFFFFFFF
)

// Radix32 implements a radix tree with an uint32 as its key.
type Radix32 struct {
	branch [2]*Radix32 // branch[0] is left branch for 0, and branch[1] the right for 1
	parent *Radix32
	key    uint32      // the key under which this value is stored
	bits   int         // the number of significant bits, if 0 the key has not been set.
	Value  interface{} // The value stored.
}

// New32 returns an empty, initialized Radix32 tree.
func New32() *Radix32 {
	return &Radix32{[2]*Radix32{nil, nil}, nil, 0, 0, nil}
}

// Key returns the key under which this node is stored.
func (r *Radix32) Key() uint32 {
	return r.key
}

// Bits returns the number of significant bits for the key.
// A value of zero indicates a key that has not been set.
func (r *Radix32) Bits() int {
	return r.bits
}

// Leaf returns true is r is an leaf node, when false is returned
// the node is a non-leaf node.
func (r *Radix32) Leaf() bool {
	return r.branch[0] == nil && r.branch[1] == nil
}

// Insert inserts a new value n in the tree r (possibly silently overwriting an existing value). 
// The first bits bits of n are significant and used to store the value v.
// It returns the inserted node, r must be the root of the tree.
func (r *Radix32) Insert(n uint32, bits int, v interface{}) *Radix32 {
	return r.insert(n, bits, v, bitSize32-1)
}

// Remove removes a value from the tree r. It returns the node removed, or nil
// when nothing is found, r must be the root of the tree.
func (r *Radix32) Remove(n uint32, bits int) *Radix32 {
	return r.remove(n, bits, bitSize32-1)
}

// Find searches the tree for the key n, where the first bits bits of n 
// are significant. It returns the node found.
func (r *Radix32) Find(n uint32, bits int) *Radix32 {
	return r.find(n, bits, bitSize32-1, nil)
}

// Do traverses the tree r in breadth-first order. For each visited node,
// the function f is called with the current node, the level of the node
// (starting with 0 for the root), and the branch taken
// (0 for the zero, 1 for the one branch, -1 is used for the root node).
func (r *Radix32) Do(f func(*Radix32, int, int)) {
	q := make(queue32, 0)

	level := 0 // TODO(mg): Does level really works as intended??
	q.Push(&node32{r, level, -1})
	x := q.Pop()
	for x != nil {
		f(x.Radix32, x.level, x.branch)
		for i, b := range x.Radix32.branch {
			if b != nil {
				q.Push(&node32{b, level, i})
			}
		}
		level++
		x = q.Pop()
	}
}

// Implement insert
func (r *Radix32) insert(n uint32, bits int, v interface{}, bit int) *Radix32 {
	switch r.Leaf() {
	case false:
		if r.parent != nil && r.bits == 0 { // nothing here yet, put something in
			r.bits = bits
			r.key = n
			r.Value = v
			return r
		}
		/*
			if bitSize32-bits == bit { // we need to store a value here
				println("STORING", bitSize32, bits, bit)
				// TODO(mg): check previous value? And then what?
				r.key = n
				r.bits = bits
				r.Value = v
				return r
			}
		*/
		// Non-leaf node, no key, one or two branches
		k := bitK32(n, bit)
		if r.branch[k] == nil {
			r.branch[k] = New32() // create missing branch
			r.branch[k].parent = r
		}
		return r.branch[k].insert(n, bits, v, bit-1)
	case true:
		// External node, (optional) key, no branches
		if r.bits == 0 { // nothing here yet, put something in
			r.bits = bits
			r.key = n
			r.Value = v
			return r
		}
		if r.key == n { // equal keys, overwrite
//			fmt.Printf("OVERWRITING, old value %d %d new %d %d\n", r.Value, r.bits, v, bits)
			r.bits = bits
			r.key = n
			r.Value = v
			return r
		}
		bcur := bitK32(r.key, bit)
		bnew := bitK32(n, bit)
		if bcur == bnew {
			r.branch[bcur] = New32()
			r.branch[bcur].parent = r
			// "fill" the correct node, with the current key - and call ourselves
			r.branch[bcur].key = r.key
			r.branch[bcur].Value = r.Value
			r.branch[bcur].bits = r.bits
			r.bits = 0
			r.key = 0
			r.Value = 0
			return r.branch[bnew].insert(n, bits, v, bit-1)
		}
		// not equal
		// keep current node, and branch of in child
		r.branch[bnew] = New32()
		r.branch[bnew].parent = r
		return r.branch[bnew].insert(n, bits, v, bit-1)
	}
	panic("bitradix: not reached")
}

// Walk the tree searching for n, keep the last node that has a key in tow.
// This is the node we should retreat to when we find and delete our node.
func (r *Radix32) remove(n uint32, bits, bit int) *Radix32 {
	if r.bits > 0 && r.bits == bits {
		// possible hit
		mask := uint32(mask32 << (bitSize32 - uint(r.bits)))
		if r.key&mask == n&mask {
			// save r in r1
			r1 := &Radix32{[2]*Radix32{nil, nil}, nil, r.key, r.bits, r.Value}
			r.prune(true)
			return r1
		}
	}
	k := bitK32(n, bit)
	if r.Leaf() || r.branch[k] == nil { // dead end
		return nil
	}
	return r.branch[bitK32(n, bit)].remove(n, bits, bit-1)
}

// Prune the tree, when b is true the current node is deleted.
func (r *Radix32) prune(b bool) {
	if b {
		if r.parent == nil {
			r.bits = 0
			r.key = 0
			r.Value = nil
			return
		}
		// we are a node, we have a parent, so the parent is a non-leaf node
		if r.parent.branch[0] == r {
			// kill that branch
			r.parent.branch[0] = nil
		}
		if r.parent.branch[1] == r {
			r.parent.branch[1] = nil
		}
		r.parent.prune(false)
		return
	}
	if r == nil {
		return
	}
	if r.bits != 0 {
		// fun stops
		return
	}
	// Does I have one or two childeren, if one, move my self up one node
	// Also the child must be a leaf node!
	b0 := r.branch[0]
	b1 := r.branch[1]
	if b0 != nil && b1 != nil {
		// two branches, we cannot replace ourselves with a child
		return
	}
	if b0 != nil {
		if !b0.Leaf() {
			return
		}
		// move b0 into this node	
		r.key = b0.key
		r.bits = b0.bits
		r.Value = b0.Value
		r.branch[0] = b0.branch[0]
		r.branch[1] = b0.branch[1]
	}
	if b1 != nil {
		if !b1.Leaf() {
			return
		}
		// move b1 into this node
		r.key = b1.key
		r.bits = b1.bits
		r.Value = b1.Value
		r.branch[0] = b1.branch[0]
		r.branch[1] = b1.branch[1]
	}
	r.parent.prune(false)
}

// Search the tree, when "seeing" a node with a key, store that
// node, when we don't find anything within the allowed bit bits
// we return that one.
func (r *Radix32) find(n uint32, bits, bit int, last *Radix32) *Radix32 {
	k := bitK32(n, bit)
	if r == nil {
		return last
	}
	switch r.Leaf() {
	case false:
		if r.bits != 0 {
			// TODO(mg) double check, think this is correct i.e using bits, not r.bits
			// is is our key
			mask := uint32(mask32 << (bitSize32 - uint(bits)))
			if r.key&mask == n&mask {
				return r
			}
			// See if it fits our key and mask
			mask = uint32(mask32 << (bitSize32 - uint(r.bits)))
			if r.key&mask == n&mask { // bits are the same, drag the key with us
				// but only if the mask is larger
				if last == nil || r.bits > last.bits {
					//fmt.Printf("MORE %032b %d\n", r.key, r.Value)
					return r.branch[k].find(n, bits, bit-1, r)
				}
			}
			return r.branch[k].find(n, bits, bit-1, last)
		}
		if r.branch[k] == nil {
			return last
		}
		return r.branch[k].find(n, bits, bit-1, last)
	case true:
		mask := uint32(mask32 << (bitSize32 - uint(r.bits)))
//		fmt.Printf("mask %032b %d\n", mask, r.bits)
		if r.key&mask == n&mask {
			if last == nil || r.bits > last.bits {
//			println("HIT")
			return r
		}
		}
//		fmt.Printf("%032b %d\n", n, r.Value)
		return last
	}
	panic("bitradix: not reached")
}

// From: http://stackoverflow.com/questions/2249731/how-to-get-bit-by-bit-data-from-a-integer-value-in-c

// Return bit k from n. We count from the right, MSB left.
// So k = 0 is the last bit on the left and k = 31 is the first bit on the right.
func bitK32(n uint32, k int) byte {
	return byte((n & (1 << uint(k))) >> uint(k))
}
