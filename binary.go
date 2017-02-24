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
	var buf bytes.Buffer

	if err := n.writeAsBinary(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (n *Node) writeAsBinary(w io.Writer) error {
	for c := n; c != nil; c = c.NextChild() {
		var err error
		switch c.value.(type) {
		case nil:
			_, err = w.Write([]byte{ptNone})
		case string:
			_, err = w.Write([]byte{ptString})
		case int32:
			_, err = w.Write([]byte{ptInt})
		case float32:
			_, err = w.Write([]byte{ptFloat})
		case uint32:
			_, err = w.Write([]byte{ptPtr})
		case []uint16:
			_, err = w.Write([]byte{ptWString})
		case color.NRGBA:
			_, err = w.Write([]byte{ptColor})
		case uint64:
			_, err = w.Write([]byte{ptUint64})
		default:
			panic("invalid vdf.Node")
		}
		if err != nil {
			return err
		}

		name := c.name
		if i := strings.IndexByte(name, 0); i != -1 {
			name = name[:i]
		}
		if _, err = io.WriteString(w, name); err != nil {
			return err
		}
		if _, err = w.Write([]byte{0}); err != nil {
			return err
		}

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
			panic("unreachable")
		}
		if err != nil {
			return err
		}
	}
	if n != nil {
		if _, err := w.Write([]byte{ptNullMarker}); err != nil {
			return err
		}
	}
	return nil
}

func (n *Node) UnmarshalBinary(b []byte) error {
	*n = Node{}
	return n.readAsBinary(bufio.NewReader(bytes.NewReader(b)))
}

func (n *Node) readAsBinary(r *bufio.Reader) error {
	packType, err := r.ReadByte()
	if err != nil {
		return err
	}

	for c := n; packType != ptNullMarker; {
		name, err := r.ReadString(0)
		if err != nil {
			return err
		}
		c.SetName(strings.TrimSuffix(name, "\x00"))

		switch packType {
		case ptNone:
			var sub Node
			c.Append(&sub)
			if err := sub.readAsBinary(r); err != nil {
				return err
			}
		case ptString:
			v, err := r.ReadString(0)
			if err != nil {
				return err
			}
			c.value = strings.TrimSuffix(v, "\x00")
		case ptWString:
			var length uint16
			if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
				return err
			}
			v := make([]uint16, length)
			for i := range v {
				if err := binary.Read(r, binary.LittleEndian, &v[i]); err != nil {
					return err
				}
			}
			c.value = v
		case ptInt:
			var v int32
			if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
				return err
			}
			c.value = v
		case ptUint64:
			var v uint64
			if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
				return err
			}
			c.value = v
		case ptFloat:
			var v float32
			if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
				return err
			}
			c.value = v
		case ptColor:
			var v color.NRGBA
			if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
				return err
			}
			c.value = v
		case ptPtr:
			var v uint32
			if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
				return err
			}
			c.value = v
		default:
			return fmt.Errorf("vdf: unknown pack type %d", packType)
		}

		packType, err = r.ReadByte()
		if err != nil {
			return err
		}

		if packType == ptNullMarker {
			break
		}

		var peer Node
		peer.parent = c.parent
		c.next = &peer
		peer.prev = c
		c = &peer
	}

	return nil
}
