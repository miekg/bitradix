package bitradix

// Radix64 implements a radix tree with an uint64 as its key.
type Radix64 struct {
	branch [2]*Radix64 // branch[0] is left branch for 0, and branch[1] the right for 1
	key    uint64      // the key under which this value is stored
	bits   int         // the number of significant bits, if 0 the key has not been set.
	Value  uint32      // The value stored.
	leaf   bool        // leaf node
}

// New64 returns an empty, initialized Radix64 tree.
func New64() *Radix64 {
	return &Radix64{[2]*Radix64{nil, nil}, 0, 0, 0, true}
}

// Key returns the key under which this node is stored.
func (r *Radix64) Key() uint64 {
	return r.key
}

// Bits returns the number of significant bits for the key.
// A value of zero indicates a key that has not been set.
func (r *Radix64) Bits() int {
	return r.bits
}

// Leaf returns true is r is an leaf node, when false is returned
// the node is a non-leaf node.
func (r *Radix64) Leaf() bool {
	return r.leaf
}

// Insert inserts a new value n in the tree r. The first bits bits of n are significant
// and used to store the value v.
// It returns the inserted node, r must be the root of the tree.
func (r *Radix64) Insert(n uint64, bits int, v uint32) *Radix64 {
	return r.insert(n, bits, v, bitSize64-1)
}

// Remove removes a value from the tree r. It returns the node removed, or nil
// when nothing is found. r must be the root of the tree.
func (r *Radix64) Remove(n uint64, bits int) *Radix64 {
	return nil
}

// Find searches the tree for the key n, where the first bits bits of n 
// are significant. It returns the node found.
func (r *Radix64) Find(n uint64, bits int) *Radix64 {
	return r.find(n, bits, bitSize64-1, nil)
}

// Do traverses the tree r in breadth-first order. For each visited node,
// the function f is called with the current node and the branch taken
// (0 for the zero, 1 for the one branch, -1 is used for the root node).
func (r *Radix64) Do(f func(*Radix64, int)) {
	q := make(queue64, 0)

	q.Push(&node64{r, -1})
	x := q.Pop()
	for x != nil {
		f(x.Radix64, x.branch)
		for i, b := range x.Radix64.branch {
			if b != nil {
				q.Push(&node64{b, i})
			}
		}
		x = q.Pop()
	}
}

// See Radix32.insert for the docs.
func (r *Radix64) insert(n uint64, bits int, v uint32, bit int) *Radix64 {
	return nil
}

// See Radix32.find for the docs.
func (r *Radix64) find(n uint64, bits, bit int, last *Radix64) *Radix64 {
	switch r.leaf {
	case false:
		if r.bits != 0 {
			// Actual key, drag it along
			return r.branch[bitK64(n, bit)].find(n, bits, bit-1, r)
		}
		return r.branch[bitK64(n, bit)].find(n, bits, bit-1, last)
	case true:
		mask := uint64(0xFFFFFFFFFFFFFFFF << uint(r.bits))
		if r.key&mask == n&mask {
			return r
		}
		return last
	}
	panic("bitradix: not reached")
}

// Return bit k from n. We count from the right, MSB left.
// So k = 0 is the last bit on the left and k = 63 is the first bit on the right.
func bitK64(n uint64, k int) byte {
	return byte((n & (1 << uint(k))) >> uint(k))
}
