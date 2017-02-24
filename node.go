package vdf

type Node struct {
	condition string
	name      string
	parent    *Node
	prev      *Node
	next      *Node
	child     *Node
	// one of:
	// - interface nil
	// - string
	// - int32
	// - float32
	// - uint32
	// - []uint16
	// - color.NRGBA
	// - uint64
	value interface{}
	cf    *customFormat
}

var blankNode Node

func (n *Node) notNil() *Node {
	if n != nil {
		return n
	}
	return &blankNode
}

func (n *Node) Name() string        { return n.notNil().name }
func (n *Node) SetName(name string) { n.name = name }

func (n *Node) Condition() string             { return n.notNil().condition }
func (n *Node) SetCondition(condition string) { n.condition = condition }

func (n *Node) ClearFormatting() {
	n.cf = nil

	for c := n.FirstChild(); c != nil; c = c.NextChild() {
		c.ClearFormatting()
	}
}

type customFormat struct {
	before        string
	condition     string
	between       string
	after         string
	unquotedKey   bool
	unquotedValue bool
}
