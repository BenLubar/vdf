package vdf_test

import (
	"encoding"
	"image/color"
	"reflect"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/BenLubar/vdf"
)

var _ encoding.TextMarshaler = (*vdf.Node)(nil)
var _ encoding.TextUnmarshaler = (*vdf.Node)(nil)
var _ encoding.BinaryMarshaler = (*vdf.Node)(nil)
var _ encoding.BinaryUnmarshaler = (*vdf.Node)(nil)

func TestContract(t *testing.T) {
	var accessors = []struct {
		name string
		f    func(t *testing.T, n *vdf.Node)
	}{
		{
			name: "Color",
			f: func(t *testing.T, n *vdf.Node) {
				n.Color()
			},
		},
		{
			name: "Condition",
			f: func(t *testing.T, n *vdf.Node) {
				n.Condition()
			},
		},
		{
			name: "FirstByName",
			f: func(t *testing.T, n *vdf.Node) {
				n.FirstByName("test")
			},
		},
		{
			name: "FirstChild",
			f: func(t *testing.T, n *vdf.Node) {
				n.FirstChild()
			},
		},
		{
			name: "FirstSubTree",
			f: func(t *testing.T, n *vdf.Node) {
				n.FirstSubTree()
			},
		},
		{
			name: "FirstValue",
			f: func(t *testing.T, n *vdf.Node) {
				n.FirstValue()
			},
		},
		{
			name: "Float",
			f: func(t *testing.T, n *vdf.Node) {
				n.Float()
			},
		},
		{
			name: "Int",
			f: func(t *testing.T, n *vdf.Node) {
				n.Int()
			},
		},
		{
			name: "MarshalBinary",
			f: func(t *testing.T, n *vdf.Node) {
				_, err := n.MarshalBinary()
				if err != nil {
					t.Error(err)
				}
			},
		},
		{
			name: "MarshalText",
			f: func(t *testing.T, n *vdf.Node) {
				_, err := n.MarshalText()
				if err != nil {
					t.Error(err)
				}
			},
		},
		{
			name: "Name",
			f: func(t *testing.T, n *vdf.Node) {
				n.Name()
			},
		},
		{
			name: "NextByName",
			f: func(t *testing.T, n *vdf.Node) {
				n.NextByName("test")
			},
		},
		{
			name: "NextChild",
			f: func(t *testing.T, n *vdf.Node) {
				n.NextChild()
			},
		},
		{
			name: "NextSubTree",
			f: func(t *testing.T, n *vdf.Node) {
				n.NextSubTree()
			},
		},
		{
			name: "NextValue",
			f: func(t *testing.T, n *vdf.Node) {
				n.NextValue()
			},
		},
		{
			name: "Ptr",
			f: func(t *testing.T, n *vdf.Node) {
				n.Ptr()
			},
		},
		{
			name: "String",
			f: func(t *testing.T, n *vdf.Node) {
				n.String()
			},
		},
		{
			name: "Uint64",
			f: func(t *testing.T, n *vdf.Node) {
				n.Uint64()
			},
		},
		{
			name: "WString",
			f: func(t *testing.T, n *vdf.Node) {
				n.WString()
			},
		},
	}
	var mutators = []struct {
		name string
		f    func(t *testing.T, n *vdf.Node)
	}{
		{
			name: "Append",
			f: func(t *testing.T, n *vdf.Node) {
				var c vdf.Node
				n.Append(&c)
			},
		},
		{
			name: "ClearFormatting",
			f: func(t *testing.T, n *vdf.Node) {
				n.ClearFormatting()
			},
		},
		{
			name: "Remove",
			f: func(t *testing.T, n *vdf.Node) {
				n.Remove()
			},
		},
		{
			name: "SetColor",
			f: func(t *testing.T, n *vdf.Node) {
				n.SetColor(color.NRGBA{255, 0, 0, 255})
			},
		},
		{
			name: "SetCondition",
			f: func(t *testing.T, n *vdf.Node) {
				n.SetCondition("$WIN32")
			},
		},
		{
			name: "SetFloat",
			f: func(t *testing.T, n *vdf.Node) {
				n.SetFloat(3.14)
			},
		},
		{
			name: "SetInt",
			f: func(t *testing.T, n *vdf.Node) {
				n.SetInt(42)
			},
		},
		{
			name: "SetName",
			f: func(t *testing.T, n *vdf.Node) {
				n.SetName("test")
			},
		},
		{
			name: "SetPtr",
			f: func(t *testing.T, n *vdf.Node) {
				n.SetPtr(0xdeadbeef)
			},
		},
		{
			name: "SetString",
			f: func(t *testing.T, n *vdf.Node) {
				n.SetString("Hello, World!")
			},
		},
		{
			name: "SetUint64",
			f: func(t *testing.T, n *vdf.Node) {
				n.SetUint64(1234567890123456789)
			},
		},
		{
			name: "SetWString",
			f: func(t *testing.T, n *vdf.Node) {
				n.SetWString([]uint16{0xd83d, 0xdca9})
			},
		},
		{
			name: "UnmarshalBinary",
			f: func(t *testing.T, n *vdf.Node) {
				err := n.UnmarshalBinary([]byte{8})
				if err != nil {
					t.Error(err)
				}
			},
		},
		{
			name: "UnmarshalText",
			f: func(t *testing.T, n *vdf.Node) {
				err := n.UnmarshalText([]byte("testing 123"))
				if err != nil {
					t.Error(err)
				}
			},
		},
	}
	typ := reflect.TypeOf((*vdf.Node)(nil))
	var nodeForAccessors vdf.Node
	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)
		if r, _ := utf8.DecodeRuneInString(m.Name); !unicode.IsUpper(r) {
			continue
		}
		found := false
		for _, c := range accessors {
			c := c // shadow
			if c.name == m.Name {
				t.Run("Accessor/"+m.Name, func(t *testing.T) {
					t.Parallel()
					c.f(t, &nodeForAccessors)
				})
				t.Run("Null/"+m.Name, func(t *testing.T) {
					t.Parallel()
					c.f(t, nil)
				})
				found = true
				break
			}
		}
		if found {
			continue
		}
		for _, c := range mutators {
			c := c // shadow
			if c.name == m.Name {
				t.Run("Mutator/"+m.Name, func(t *testing.T) {
					t.Parallel()
					c.f(t, new(vdf.Node))
				})
				found = true
				break
			}
		}
		if found {
			continue
		}
		t.Run("Missing/"+m.Name, func(t *testing.T) {
			t.Errorf("Missing test for method: %q", m.Name)
		})
	}
}
