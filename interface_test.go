package vdf_test

import (
	"encoding"

	"github.com/BenLubar/vdf"
)

var _ encoding.TextMarshaler = (*vdf.Node)(nil)
var _ encoding.TextUnmarshaler = (*vdf.Node)(nil)
var _ encoding.BinaryMarshaler = (*vdf.Node)(nil)
var _ encoding.BinaryUnmarshaler = (*vdf.Node)(nil)
