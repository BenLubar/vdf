package vdf

import "encoding"

var _ encoding.TextMarshaler = (*Node)(nil)
var _ encoding.TextUnmarshaler = (*Node)(nil)
var _ encoding.BinaryMarshaler = (*Node)(nil)
var _ encoding.BinaryUnmarshaler = (*Node)(nil)
