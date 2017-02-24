package vdf

import (
	"fmt"
	"image/color"
	"strconv"
	"unicode/utf16"
)

func (n *Node) String() string {
	if n == nil || n.child != nil {
		return ""
	}

	switch v := n.value.(type) {
	case nil:
		return ""
	case string:
		return v
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'g', -1, 32)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case []uint16:
		return string(utf16.Decode(v))
	case color.NRGBA:
		return fmt.Sprintf("%d %d %d %d", v.R, v.G, v.B, v.A)
	case uint64:
		return strconv.FormatUint(v, 10)
	}
	panic("invalid vdf.Node")
}

func (n *Node) SetString(s string) {
	for n.child != nil {
		n.child.Remove()
	}

	n.value = s
}

func (n *Node) Int() int32 {
	if n == nil || n.child != nil {
		return 0
	}

	switch v := n.value.(type) {
	case nil:
		return 0
	case string:
		i, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return 0
		}
		return int32(i)
	case int32:
		return v
	case float32:
		return int32(v)
	case uint32:
		return int32(v)
	case []uint16:
		i, err := strconv.ParseInt(string(utf16.Decode(v)), 10, 32)
		if err != nil {
			return 0
		}
		return int32(i)
	case color.NRGBA:
		return 0
	case uint64:
		return int32(v)
	}
	panic("invalid vdf.Node")
}

func (n *Node) SetInt(i int32) {
	for n.child != nil {
		n.child.Remove()
	}

	n.value = i
}

func (n *Node) Float() float32 {
	if n == nil || n.child != nil {
		return 0
	}

	switch v := n.value.(type) {
	case nil:
		return 0
	case string:
		f, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return 0
		}
		return float32(f)
	case int32:
		return float32(v)
	case float32:
		return v
	case uint32:
		return float32(v)
	case []uint16:
		f, err := strconv.ParseFloat(string(utf16.Decode(v)), 32)
		if err != nil {
			return 0
		}
		return float32(f)
	case color.NRGBA:
		return 0
	case uint64:
		return float32(v)
	}
	panic("invalid vdf.Node")
}

func (n *Node) SetFloat(f float32) {
	for n.child != nil {
		n.child.Remove()
	}

	n.value = f
}

func (n *Node) Ptr() uint32 {
	if n == nil || n.child != nil {
		return 0
	}

	switch v := n.value.(type) {
	case nil:
		return 0
	case string:
		i, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return 0
		}
		return uint32(i)
	case int32:
		return uint32(v)
	case float32:
		return uint32(v)
	case uint32:
		return v
	case []uint16:
		i, err := strconv.ParseUint(string(utf16.Decode(v)), 10, 32)
		if err != nil {
			return 0
		}
		return uint32(i)
	case color.NRGBA:
		return 0
	case uint64:
		return uint32(v)
	}
	panic("invalid vdf.Node")
}

func (n *Node) SetPtr(i uint32) {
	for n.child != nil {
		n.child.Remove()
	}

	n.value = i
}

func (n *Node) WString() []uint16 {
	if n == nil || n.child != nil {
		return nil
	}

	switch v := n.value.(type) {
	case []uint16:
		c := make([]uint16, len(v))
		copy(c, v)
		return c
	default:
		return utf16.Encode([]rune(n.String()))
	}
}

func (n *Node) SetWString(s []uint16) {
	for n.child != nil {
		n.child.Remove()
	}

	c := make([]uint16, len(s))
	copy(c, s)
	n.value = c
}

func (n *Node) Color() color.NRGBA {
	if n == nil || n.child != nil {
		return color.NRGBA{}
	}

	switch v := n.value.(type) {
	case nil:
		return color.NRGBA{}
	case string:
		var c color.NRGBA
		_, err := fmt.Sscanf(v, "%d %d %d %d", &c.R, &c.G, &c.B, &c.A)
		if err != nil {
			return color.NRGBA{}
		}
		return c
	case int32:
		return color.NRGBA{}
	case float32:
		return color.NRGBA{}
	case uint32:
		return color.NRGBA{}
	case []uint16:
		var c color.NRGBA
		_, err := fmt.Sscanf(string(utf16.Decode(v)), "%d %d %d %d", &c.R, &c.G, &c.B, &c.A)
		if err != nil {
			return color.NRGBA{}
		}
		return c
	case color.NRGBA:
		return v
	case uint64:
		return color.NRGBA{}
	}
	panic("invalid vdf.Node")
}

func (n *Node) SetColor(c color.NRGBA) {
	for n.child != nil {
		n.child.Remove()
	}

	n.value = c
}

func (n *Node) Uint64() uint64 {
	if n == nil || n.child != nil {
		return 0
	}

	switch v := n.value.(type) {
	case nil:
		return 0
	case string:
		i, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0
		}
		return i
	case int32:
		return uint64(v)
	case float32:
		return uint64(v)
	case uint32:
		return uint64(v)
	case []uint16:
		i, err := strconv.ParseUint(string(utf16.Decode(v)), 10, 64)
		if err != nil {
			return 0
		}
		return i
	case color.NRGBA:
		return 0
	case uint64:
		return v
	}
	panic("invalid vdf.Node")
}

func (n *Node) SetUint64(i uint64) {
	for n.child != nil {
		n.child.Remove()
	}

	n.value = i
}
