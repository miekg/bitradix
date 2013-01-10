package bitradix

// Radix64 implements a radix tree with an uint64 as its key. The methods
// are identical to those of Radix32, except for the key length.
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
	return nil
}

func (r *Radix64) remove(uint64, bits int, v uint32, bit int) *Radix64 {
	return nil
}

func (r *Radix64) prune(b bool) {
	return
}

func (r *Radix64) find(n uint64, bits, bit int, last *Radix64) *Radix64 {
	switch r.Leaf() {
	case false:
		if r.bits != 0 {
			// Actual key, drag it along
			return r.branch[bitK64(n, bit)].find(n, bits, bit-1, r)
		}
		return r.branch[bitK64(n, bit)].find(n, bits, bit-1, last)
	case true:
		mask := uint64(mask64 << uint(r.bits))
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
