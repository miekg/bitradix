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
