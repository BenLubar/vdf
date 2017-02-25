package vdf

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"image/color"
	"io"
	"strings"
)

// EPackType from KVPacker in Source SDK 2013:
// https://github.com/ValveSoftware/source-sdk-2013/blob/master/mp/src/public/tier1/kvpacker.h
const (
	ptNone       = 0
	ptString     = 1
	ptInt        = 2
	ptFloat      = 3
	ptPtr        = 4
	ptWString    = 5
	ptColor      = 6
	ptUint64     = 7
	ptNullMarker = 8
)

func (n *Node) MarshalBinary() ([]byte, error) {
	if n == nil {
		return nil, nil
	}

	var buf bytes.Buffer

	if err := n.writeAsBinary(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (n *Node) writeAsBinary(w io.Writer) error {
	for c := n; c != nil; c = c.NextChild() {
		if _, err := w.Write([]byte{packType(c.value)}); err != nil {
			return err
		}

		name := c.name
		if i := strings.IndexByte(name, 0); i != -1 {
			name = name[:i]
		}
		if _, err := io.WriteString(w, name); err != nil {
			return err
		}
		if _, err := w.Write([]byte{0}); err != nil {
			return err
		}

		var err error
		switch v := c.value.(type) {
		case nil:
			err = c.child.writeAsBinary(w)
		case string:
			if i := strings.IndexByte(v, 0); i != -1 {
				v = v[:i]
			}
			if _, err = io.WriteString(w, v); err != nil {
				return err
			}
			_, err = w.Write([]byte{0})
		case int32:
			err = binary.Write(w, binary.LittleEndian, &v)
		case float32:
			err = binary.Write(w, binary.LittleEndian, &v)
		case uint32:
			err = binary.Write(w, binary.LittleEndian, &v)
		case []uint16:
			err = binary.Write(w, binary.LittleEndian, uint16(len(v)))
			for i := range v {
				if err != nil {
					return err
				}
				err = binary.Write(w, binary.LittleEndian, v[i])
			}
		case color.NRGBA:
			err = binary.Write(w, binary.LittleEndian, &v)
		case uint64:
			err = binary.Write(w, binary.LittleEndian, &v)
		default:
			panic("invalid vdf.Node")
		}
		if err != nil {
			return err
		}
	}
	if _, err := w.Write([]byte{ptNullMarker}); err != nil {
		return err
	}
	return nil
}

// packType returns the EPackType for the value inside the interface.
func packType(v interface{}) byte {
	switch v.(type) {
	case nil:
		return ptNone
	case string:
		return ptString
	case int32:
		return ptInt
	case float32:
		return ptFloat
	case uint32:
		return ptPtr
	case []uint16:
		return ptWString
	case color.NRGBA:
		return ptColor
	case uint64:
		return ptUint64
	default:
		panic("invalid vdf.Node")
	}
}

func (n *Node) UnmarshalBinary(b []byte) error {
	*n = Node{}
	return n.readAsBinary(bufio.NewReader(bytes.NewReader(b)), nil)
}

func (n *Node) readAsBinary(r *bufio.Reader, parent *Node) error {
	pt, err := r.ReadByte()
	if err != nil {
		return err
	}

	for c := n; pt != ptNullMarker; {
		c.parent = parent
		name, err := r.ReadString(0)
		if err != nil {
			return err
		}
		c.SetName(strings.TrimSuffix(name, "\x00"))

		switch pt {
		case ptNone:
			var sub Node
			c.Append(&sub)
			sub.parent = nil
			err = sub.readAsBinary(r, c)
			if sub.parent == nil {
				c.child = nil
			}
		case ptString:
			var v string
			v, err = r.ReadString(0)
			c.value = strings.TrimSuffix(v, "\x00")
		case ptWString:
			var length uint16
			if err = binary.Read(r, binary.LittleEndian, &length); err != nil {
				return err
			}
			v := make([]uint16, length)
			for i := range v {
				if err = binary.Read(r, binary.LittleEndian, &v[i]); err != nil {
					return err
				}
			}
			c.value = v
		case ptInt:
			var v int32
			err = binary.Read(r, binary.LittleEndian, &v)
			c.value = v
		case ptUint64:
			var v uint64
			err = binary.Read(r, binary.LittleEndian, &v)
			c.value = v
		case ptFloat:
			var v float32
			err = binary.Read(r, binary.LittleEndian, &v)
			c.value = v
		case ptColor:
			var v color.NRGBA
			err = binary.Read(r, binary.LittleEndian, &v)
			c.value = v
		case ptPtr:
			var v uint32
			err = binary.Read(r, binary.LittleEndian, &v)
			c.value = v
		default:
			err = fmt.Errorf("vdf: unknown pack type %d", pt)
		}
		if err != nil {
			return err
		}

		pt, err = r.ReadByte()
		if err != nil {
			return err
		}

		if pt == ptNullMarker {
			break
		}

		var peer Node
		c.next = &peer
		peer.prev = c
		c = &peer
	}

	return nil
}
