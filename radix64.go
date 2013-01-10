package bitradix

// Radix64 implements a radix tree with an uint64 as its key.
// The methods are identical to those of Radix32, except for the key length used.
type Radix64 struct {
	branch [2]*Radix64 // branch[0] is left branch for 0, and branch[1] the right for 1
	parent *Radix64
	key    uint64 // the key under which this value is stored
	bits   int    // the number of significant bits, if 0 the key has not been set.
	Value  uint32 // The value stored.
}

func New64() *Radix64 {
	return &Radix64{[2]*Radix64{nil, nil}, nil, 0, 0, 0}
}

func (r *Radix64) Key() uint64 {
	return r.key
}

func (r *Radix64) Bits() int {
	return r.bits
}

func (r *Radix64) Leaf() bool {
	return r.branch[0] == nil && r.branch[1] == nil
}

func (r *Radix64) Insert(n uint64, bits int, v uint32) *Radix64 {
	return r.insert(n, bits, v, bitSize64-1)
}

func (r *Radix64) Remove(n uint64, bits int) *Radix64 {
	return r.remove(n, bits, bitSize64-1)
}

func (r *Radix64) Find(n uint64, bits int) *Radix64 {
	return r.find(n, bits, bitSize64-1, nil)
}

func (r *Radix64) Do(f func(*Radix64, int, int)) {
	q := make(queue64, 0)

	level := 0
	q.Push(&node64{r, level, -1})
	x := q.Pop()
	for x != nil {
		f(x.Radix64, x.level, x.branch)
		for i, b := range x.Radix64.branch {
			if b != nil {
				q.Push(&node64{b, level, i})
			}
		}
		level++
		x = q.Pop()
	}
}

func (r *Radix64) insert(n uint64, bits int, v uint32, bit int) *Radix64 {
	switch r.Leaf() {
	case false:
		if bitSize64-bits == bit { // we need to store a value here
			r.key = n
			r.bits = bits
			r.Value = v
			return r
		}
		// Non-leaf node, no key, one or two branches
		k := bitK64(n, bit)
		if r.branch[k] == nil {
			r.branch[k] = New64() // create missing branch
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
		if bitSize64-bits == bit { // seen all bits, put something here
			r.bits = bits
			r.key = n
			r.Value = v
			return r
		}

		bcur := bitK64(r.key, bit)
		bnew := bitK64(n, bit)

		switch x := bitSize64 - r.bits; true {
		case x == bit: // current node needs to stay here
			r.branch[bnew] = New64()
			r.branch[bnew].parent = r
			r.branch[bnew].key = n
			r.branch[bnew].Value = v
			r.branch[bnew].bits = bits
			return r.branch[bnew]
		case x < bit: // current node can be put one level down
			r.branch[bcur] = New64()
			r.branch[bcur].parent = r
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
			r.branch[bnew] = New64()
			r.branch[bnew].parent = r
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
		case x > bit:
			panic("bitradix: node put too far down")
		}
	}
	panic("bitradix: not reached")
	return nil
}

func (r *Radix64) remove(n uint64, bits, bit int) *Radix64 {
	if r.bits > 0 && r.bits == bits {
		// possible hit
		mask := uint64(mask64 << (bitSize64 - uint(r.bits)))
		if r.key&mask == n&mask {
			// save r in r1
			r1 := &Radix64{[2]*Radix64{nil, nil}, nil, r.key, r.bits, r.Value}
			r.prune(true)
			return r1
		}
	}
	k := bitK64(n, bit)
	if r.Leaf() || r.branch[k] == nil { // dead end
		return nil
	}
	return r.branch[bitK64(n, bit)].remove(n, bits, bit-1)
}

func (r *Radix64) prune(b bool) {
	if b {
		if r.parent == nil {
			r.bits = 0
			r.key = 0
			r.Value = 0
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

func (r *Radix64) find(n uint64, bits, bit int, last *Radix64) *Radix64 {
	k := bitK64(n, bit)
	switch r.Leaf() {
	case false:
		if r.branch[k] == nil {
			return r
		}
		if r.bits != 0 {
			// TODO(mg) double check, think this is correct i.e using bits
			mask := uint64(mask64 << (bitSize64 - uint(bits)))
			if r.key&mask == n&mask {
				return r
			}
			// A key, drag it along
			return r.branch[k].find(n, bits, bit-1, r)
		}
		return r.branch[k].find(n, bits, bit-1, last)
	case true:
		mask := uint64(mask64 << (bitSize64 - uint(r.bits)))
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
