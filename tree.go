package vdf

import "strings"

func (n *Node) advanceSimple(simple bool) *Node {
	for next := n; next != nil; next = next.next {
		if simple && next.child == nil {
			return next
		}
		if !simple && next.value == nil {
			return next
		}
	}
	return nil
}
func (n *Node) FirstChild() *Node   { return n.notNil().child }
func (n *Node) FirstSubTree() *Node { return n.notNil().child.advanceSimple(false) }
func (n *Node) FirstValue() *Node   { return n.notNil().child.advanceSimple(true) }
func (n *Node) NextChild() *Node    { return n.notNil().next }
func (n *Node) NextSubTree() *Node  { return n.notNil().next.advanceSimple(false) }
func (n *Node) NextValue() *Node    { return n.notNil().next.advanceSimple(true) }

func (n *Node) FirstByName(name string) *Node {
	for c := n.FirstChild(); c != nil; c = c.NextChild() {
		if strings.EqualFold(c.Name(), name) {
			return c
		}
	}
	return nil
}
func (n *Node) NextByName(name string) *Node {
	for c := n.NextChild(); c != nil; c = c.NextChild() {
		if strings.EqualFold(c.Name(), name) {
			return c
		}
	}
	return nil
}

func (n *Node) Append(c *Node) {
	if c.parent != nil {
		panic("vdf: cannot append a node that already has a parent")
	}

	f := &n.child
	var l *Node
	for *f != nil {
		l = *f
		f = &(*f).next
	}
	c.parent = n
	*f = c
	c.prev = l

	if n.value != nil && n.cf != nil {
		n.cf.between += "\n{\n"
		n.cf.after += "\n}\n"
	}
	n.value = nil
}

func (n *Node) Remove() {
	if n == nil || n.parent == nil {
		return
	}

	next := n.next

	if n.next != nil {
		n.next.prev = n.prev
		n.next = nil
	}

	if n.prev != nil {
		n.prev.next = next
		n.prev = nil
	} else {
		n.parent.child = next
	}
}
