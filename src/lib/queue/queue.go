package queue

import ()

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//  STRUCTS & TYPES
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// --> Node <---
// -------------------
// DESCRIPTION:
// -------------------
// Basic node stuct which holds an int value
type Node struct {
	Value int
}

// --> Queue <---
// -------------------
// DESCRIPTION:
// -------------------
// FIFO Queue struct
type Queue struct {
	nodes []*Node
	size  int
	head  int
	tail  int
	count int
}

// func NewQueue(size int) *Queue
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// Initializes a new Queue
func NewQueue(size int) *Queue {
	return &Queue{
		nodes: make([]*Node, size),
		size:  size,
	}
}

// func (q *Queue) Push(n *Node)
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// Pushes an item onto the Queue
func (q *Queue) Push(n *Node) {
	if q.head == q.tail && q.count > 0 {
		nodes := make([]*Node, len(q.nodes)+q.size)
		copy(nodes, q.nodes[q.head:])
		copy(nodes[len(q.nodes)-q.head:], q.nodes[:q.head])
		q.head = 0
		q.tail = len(q.nodes)
		q.nodes = nodes
	}
	q.nodes[q.tail] = n
	q.tail = (q.tail + 1) % len(q.nodes)
	q.count++
}

// func (q *Queue) Pop() *Node
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// Pops from the Queue
func (q *Queue) Pop() *Node {
	if q.count == 0 {
		return nil
	}
	node := q.nodes[q.head]
	q.head = (q.head + 1) % len(q.nodes)
	q.count--
	return node
}
