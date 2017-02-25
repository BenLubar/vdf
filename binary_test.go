package vdf_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/BenLubar/vdf"
)

func TestBinaryReEncode(t *testing.T) {
	in, err := ioutil.ReadFile("testdata/UserGameStatsSchema_630.bin")
	if err != nil {
		t.Fatal(err)
	}

	var n vdf.Node
	err = n.UnmarshalBinary(in)
	if err != nil {
		t.Fatal(err)
	}

	out, err := n.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(in, out) {
		t.Error("Byte slices differ!")
		t.Logf("in:  % x", in)
		t.Logf("out: % x", out)
	}
}
