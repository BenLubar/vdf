// Package vdf implements Valve Data Format, also known as KeyValues.
//
// VDF is documented on the Valve Developer Community wiki:
// https://developer.valvesoftware.com/wiki/KeyValues
//
// This package attempts to replicate the functionality of the Source SDK 2013
// version of KeyValues.
//
// https://github.com/ValveSoftware/source-sdk-2013/blob/master/mp/src/tier1/KeyValues.cpp
// https://github.com/ValveSoftware/source-sdk-2013/blob/master/mp/src/tier1/kvpacker.cpp
package vdf

import (
	"strings"
	"unicode"
)

// Node is the basic building block of VDF.
//
// All methods on a Node are either an accessor or a mutator. Accessors are
// safe to call from multiple goroutines at the same time and are safe to call
// on a nil receiver. Mutators require that the Node and all of its children
// are only being accessed by the goroutine calling the mutator.
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

// notNil returns a pointer to the zero value of Node if n is nil.
func (n *Node) notNil() *Node {
	if n != nil {
		return n
	}
	return &blankNode
}

// Name returns the name of this Node.
//
// Name is an accessor.
func (n *Node) Name() string { return n.notNil().name }

// SetName sets the name of this node.
//
// SetName is a mutator.
func (n *Node) SetName(name string) {
	n.name = name
	if n.cf != nil && n.cf.unquotedKey && (strings.IndexFunc(name, unicode.IsSpace) != -1 || strings.ContainsAny(name, "\"{}")) {
		n.cf.unquotedKey = false
	}
}

// Condition returns the condition of this Node. A node with no condition is
// represented by an empty string.
//
// Condition is an accessor.
func (n *Node) Condition() string { return n.notNil().condition }

// SetCondition sets the condition of this Node. Putting whitespace, double
// quotes, or curly braces in a condition will panic.
//
// SetCondition is a mutator.
func (n *Node) SetCondition(condition string) {
	if strings.IndexFunc(condition, unicode.IsSpace) != -1 {
		panic("vdf: condition cannot contain spaces")
	}
	if strings.ContainsAny(condition, "\"{}") {
		panic("vdf: condition cannot contain \", {, or }")
	}
	n.condition = condition
	if n.cf != nil && n.cf.condition == "" {
		n.cf.condition = " "
	}
}

// ClearFormatting resets the Node and its children to use standard formatting
// in MarshalText. The formatting is only set by UnmarshalText.
//
// ClearFormatting is a mutator.
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
