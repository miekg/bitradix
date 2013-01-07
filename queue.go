package bitradix

type node struct {
	*Radix
	branch int	// -1 root, 0 left branch, 1 right branch
}

type queue []*node

// Push adds a node to the queue.
func (q *queue) Push(n *node) {
	*q = append(*q, n)
}

// Pop removes and returns a node from the queue in first to last order.
func (q *queue) Pop() *node {
	lq := len(*q)
	if lq == 0 {
		return nil
	}
	n := (*q)[0]
	switch lq {
	case 1:
		*q = (*q)[:0]
	default:
		*q = (*q)[1:lq]
	}
	return n
}
