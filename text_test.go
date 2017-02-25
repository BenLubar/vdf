package vdf_test

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/BenLubar/vdf"
)

func TestText(t *testing.T) {
	for _, name := range []string{
		"hello",
		"hello_eof",
		"hello_quotes",
		"cond",
	} {
		name := name // shadow

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			in, err := ioutil.ReadFile(filepath.Join("testdata", "in_"+name+".txt"))
			if err != nil {
				t.Fatal("couldn't read input: ", err)
			}
			out, err := ioutil.ReadFile(filepath.Join("testdata", "out_"+name+".txt"))
			if err != nil {
				t.Fatal("couldn't read expected output: ", err)
			}

			var n vdf.Node
			err = n.UnmarshalText(in)
			if err != nil {
				t.Fatal("couldn't parse: ", err)
			}

			out1, err := n.MarshalText()
			if err != nil {
				t.Error("couldn't serialize: ", err)
			} else if !bytes.Equal(in, out1) {
				t.Error("serialized version differs!")
				t.Logf("in:  %q", in)
				t.Logf("out: %q", out1)
			}

			n.ClearFormatting()

			out2, err := n.MarshalText()
			if err != nil {
				t.Error("couldn't clean and serialize: ", err)
			} else if !bytes.Equal(out, out2) {
				t.Error("cleaned and serialized version differs!")
				t.Logf("in:  %q", out)
				t.Logf("out: %q", out2)
			}
		})
	}
}
