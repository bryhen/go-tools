package main_test

import (
	"bytes"
	"io"
	"log"
	"math"
	"testing"
)

var (
	msg = []byte("I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!I'm doing the copying wowwww!")
	out = []byte{}
)

type Wr struct {
	Out []byte
}

func (w *Wr) Write(bs []byte) (int, error) {
	w.Out = append(w.Out, bs...)
	return len(bs), nil
}

func BenchmarkRead(b *testing.B) {
	w := &Wr{}
	r := bytes.NewReader(msg)

	for range b.N {
		bs, _ := io.ReadAll(r)
		w.Write(bs)
	}
	log.Println(string(w.Out[len(w.Out)-int(math.Min(float64(len(w.Out))/2, 100)):]))
}

func BenchmarkCopy(b *testing.B) {
	w := &Wr{}
	r := bytes.NewReader(msg)

	for range b.N {
		io.Copy(w, r)
	}
	log.Println(string(w.Out[len(w.Out)-int(math.Min(float64(len(w.Out))/2, 100)):]))
}
