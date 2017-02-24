package vdf_test

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"os"
	"os/exec"
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

	if bytes.Equal(in, out) {
		return
	}

	err = ioutil.WriteFile("test_in", []byte(hex.Dump(in)), 0644)
	if err != nil {
		panic(err)
	}
	defer os.Remove("test_in")
	err = ioutil.WriteFile("test_out", []byte(hex.Dump(out)), 0644)
	if err != nil {
		panic(err)
	}
	defer os.Remove("test_out")

	txt, err := n.MarshalText()
	if err != nil {
		panic(err)
	}
	t.Log(string(txt))

	t.Error("Byte slices differ!")
	cmd := exec.Command("diff", "-u", "test_in", "test_out")
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}
