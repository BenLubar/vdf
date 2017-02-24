package vdf

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

var escapeString = strings.NewReplacer("\\", "\\\\", "\n", "\\n", "\t", "\\t", "\v", "\\v", "\b", "\\b", "\r", "\\r", "\f", "\\f", "\a", "\\a", "'", "\\'", "\"", "\\\"")
var unescapeString = strings.NewReplacer("\\\\", "\\", "\\n", "\n", "\\t", "\t", "\\v", "\v", "\\b", "\b", "\\r", "\r", "\\f", "\f", "\\a", "\a", "\\'", "'", "\\\"", "\"")

func (n *Node) MarshalText() ([]byte, error) {
	if n == nil {
		return nil, nil
	}

	var buf bytes.Buffer
	if err := n.writeIndent(&buf, 0); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (n *Node) writeIndent(w io.Writer, indent int) error {
	if n.cf != nil {
		return n.writeCustom(w, indent)
	}
	if _, err := io.WriteString(w, strings.Repeat("\t", indent)); err != nil {
		return err
	}
	if err := writeString(w, n.name); err != nil {
		return err
	}
	if _, err := io.WriteString(w, " "); err != nil {
		return err
	}
	if n.value != nil {
		if err := writeString(w, n.String()); err != nil {
			return err
		}
		if n.condition != "" {
			if _, err := fmt.Fprintf(w, " [%s]", n.condition); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(w, "\n"); err != nil {
			return err
		}
		return nil
	}
	if n.condition != "" {
		if _, err := fmt.Fprintf(w, "[%s] ", n.condition); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(w, "{\n"); err != nil {
		return err
	}
	for c := n.FirstChild(); c != nil; c = c.NextChild() {
		if err := c.writeIndent(w, indent+1); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(w, strings.Repeat("\t", indent)); err != nil {
		return err
	}
	if _, err := io.WriteString(w, "}\n"); err != nil {
		return err
	}
	return nil
}

func (n *Node) writeCustom(w io.Writer, indent int) error {
	if _, err := io.WriteString(w, n.cf.before); err != nil {
		return err
	}
	if n.cf.unquotedKey {
		if _, err := io.WriteString(w, n.name); err != nil {
			return err
		}
	} else {
		if err := writeString(w, n.name); err != nil {
			return err
		}
	}
	if n.value != nil {
		if _, err := io.WriteString(w, n.cf.between); err != nil {
			return err
		}
		if n.cf.unquotedValue {
			if _, err := io.WriteString(w, n.String()); err != nil {
				return err
			}
		} else {
			if err := writeString(w, n.String()); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(w, n.cf.condition); err != nil {
			return err
		}
		if n.condition != "" {
			if _, err := fmt.Fprintf(w, "[%s]", n.condition); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(w, n.cf.after); err != nil {
			return err
		}
		return nil
	}
	if _, err := io.WriteString(w, n.cf.condition); err != nil {
		return err
	}
	if n.condition != "" {
		if _, err := fmt.Fprintf(w, "[%s]", n.condition); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(w, n.cf.between); err != nil {
		return err
	}
	for c := n.FirstChild(); c != nil; c = c.NextChild() {
		if err := c.writeIndent(w, indent+1); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(w, n.cf.after); err != nil {
		return err
	}
	return nil
}

func writeString(w io.Writer, s string) error {
	if _, err := io.WriteString(w, "\""); err != nil {
		return err
	}
	if _, err := escapeString.WriteString(w, s); err != nil {
		return err
	}
	if _, err := io.WriteString(w, "\""); err != nil {
		return err
	}
	return nil
}

type errClose string

func (errClose) Error() string {
	return "vdf: unexpected }"
}

func (n *Node) UnmarshalText(b []byte) error {
	*n = Node{}
	return n.readAsText(bufio.NewReader(bytes.NewReader(b)))
}

func (n *Node) readAsText(r *bufio.Reader) error {
	var last *Node
	current := n
	prefix, s, wasQuoted, wasConditional, err := readToken(r)
	if err != nil {
		return err
	}
	for {
		if wasConditional {
			return fmt.Errorf("vdf: unexpected conditional %q", s)
		}
		if !wasQuoted && s == "}" {
			return errClose(prefix)
		}
		if !wasQuoted && s == "{" {
			return fmt.Errorf("vdf: unexpected %s", s)
		}
		if current == nil {
			current = new(Node)
			current.parent = last.parent
			current.prev = last
			last.next = current
		}
		current.cf = new(customFormat)
		current.cf.before = prefix
		current.cf.unquotedKey = !wasQuoted
		current.name = s
		prefix, s, wasQuoted, wasConditional, err = readToken(r)
		if err != nil {
			return err
		}
		if wasConditional {
			current.cf.condition = prefix
			current.condition = strings.TrimSuffix(strings.TrimPrefix(s, "["), "]")
			prefix, s, wasQuoted, wasConditional, err = readToken(r)
			if err != nil {
				return err
			}
			if s != "{" || wasQuoted || wasConditional {
				return fmt.Errorf("vdf: missing {")
			}
		}

		if !wasQuoted && s == "{" {
			var suffix string
			suffix, err = readLineEnding(r)
			if err != nil {
				return err
			}
			current.cf.between = prefix + s + suffix

			var c Node
			c.parent = current
			err = c.readAsText(r)
			if p, ok := err.(errClose); ok {
				prefix = string(p)
			} else if err != nil {
				return err
			} else {
				return fmt.Errorf("vdf: missing }")
			}
			if c.cf != nil {
				current.child = &c
			}

			suffix, err = readLineEnding(r)
			if err != nil {
				return err
			}
			current.cf.after = prefix + "}" + suffix

			prefix, s, wasQuoted, wasConditional, err = readToken(r)
		} else {
			var suffix string
			suffix, err = readLineEnding(r)
			if err != nil {
				return err
			}
			current.cf.between = prefix
			current.cf.unquotedValue = !wasQuoted
			current.value = s
			current.cf.after = suffix

			prefix, s, wasQuoted, wasConditional, err = readToken(r)
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			if wasConditional {
				current.cf.condition = suffix + prefix
				current.condition = strings.TrimSuffix(strings.TrimPrefix(s, "["), "]")
				suffix, err = readLineEnding(r)
				if err != nil {
					return err
				}
				current.cf.after = suffix
				prefix, s, wasQuoted, wasConditional, err = readToken(r)
			}
		}

		last = current
		current = nil
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func readToken(r *bufio.Reader) (prefix, s string, wasQuoted, wasConditional bool, err error) {
	var prefixBuf []byte
	for {
		for {
			b, err := r.ReadByte()
			if err != nil {
				return string(prefixBuf), "", false, false, err
			}
			if !unicode.IsSpace(rune(b)) {
				if err = r.UnreadByte(); err != nil {
					return string(prefixBuf), "", false, false, err
				}
				break
			}
			prefixBuf = append(prefixBuf, b)
		}

		peek, err := r.Peek(2)
		if err == io.EOF {
			break
		}
		if err != nil {
			return string(prefixBuf), "", false, false, err
		}
		if peek[0] != '/' || peek[1] != '/' {
			break
		}

		if _, err := r.Discard(2); err != nil {
			return string(prefixBuf), "", false, false, err
		}
		prefixBuf = append(prefixBuf, '/', '/')

		line, err := r.ReadSlice('\n')
		prefixBuf = append(prefixBuf, line...)
		if err != nil {
			return string(prefixBuf), "", false, false, err
		}
	}
	prefix = string(prefixBuf)

	c, err := r.ReadByte()
	if err != nil {
		return prefix, "", false, false, err
	}

	var buf []byte

	if c == '"' {
		wasQuoted = true
		s, err = readQuoted(r)
		return
	}

	if c == '{' || c == '}' {
		return prefix, string(c), false, false, err
	}

	buf = append(buf, c)
	conditionalStart := false
	for {
		c, err = r.ReadByte()
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			s = string(buf)
			return
		}

		if c == '"' || c == '{' || c == '}' {
			if err = r.UnreadByte(); err != nil {
				s = string(buf)
				return
			}
			break
		}

		if c == '[' {
			conditionalStart = true
		}

		if c == ']' && conditionalStart {
			wasConditional = true
		}

		if unicode.IsSpace(rune(c)) {
			if err = r.UnreadByte(); err != nil {
				s = string(buf)
				return
			}
			break
		}

		buf = append(buf, c)
	}

	s = string(buf)
	return
}

func readLineEnding(r *bufio.Reader) (string, error) {
	var buf []byte
	for {
		b, err := r.ReadByte()
		if err == io.EOF {
			return string(buf), nil
		}
		if err != nil {
			return string(buf), err
		}
		if !unicode.IsSpace(rune(b)) {
			if err = r.UnreadByte(); err != nil {
				return string(buf), err
			}
			break
		}
		buf = append(buf, b)
		if b == '\n' {
			return string(buf), nil
		}
	}

	peek, err := r.Peek(2)
	if err == io.EOF {
		return string(buf), nil
	}
	if err != nil {
		return string(buf), err
	}
	if peek[0] != '/' || peek[1] != '/' {
		return string(buf), nil
	}

	if _, err := r.Discard(2); err != nil {
		return string(buf), nil
	}
	buf = append(buf, '/', '/')

	line, err := r.ReadSlice('\n')
	buf = append(buf, line...)
	if err == io.EOF {
		return string(buf), nil
	}
	return string(buf), err
}

func readQuoted(r *bufio.Reader) (string, error) {
	var buf []byte
	for {
		c, err := r.ReadByte()
		if err != nil {
			return "", err
		}

		if c == '"' {
			return string(buf), nil
		}

		if c == '\\' {
			c, err = r.ReadByte()
			if err != nil {
				return "", err
			}

			src := "\\" + string(c)
			if dst := unescapeString.Replace(src); src != dst {
				buf = append(buf, dst...)
				continue
			}

			if err = r.UnreadByte(); err != nil {
				return "", err
			}

			c = '\\'
		}

		buf = append(buf, c)
	}
}
