package splitstore

import (
	"io"
	"math/rand"
	"time"
)

func Xor(buffs [][]byte) []byte {
	out := make([]byte, len(buffs[0]))
	for _, buf := range buffs {
		for i := range buf {
			out[i] = out[i] ^ buf[i]
		}
	}
	return out
}

func Create(buf []byte, parts int) [][]byte {
	src := rand.New(rand.NewSource(time.Now().Unix()))
	buffs := make([][]byte, parts)
	for i := range buffs {
		buffs[i] = make([]byte, len(buf))
		if i < len(buffs)-1 {
			io.ReadFull(src, buffs[i])
		} else {
			buffs[i] = buf
			buffs[i] = Xor(buffs)
		}
	}
	return buffs
}
