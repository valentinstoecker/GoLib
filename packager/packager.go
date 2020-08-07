package packager

import (
	"encoding/binary"
	"io"
)

var chanBufSize = 16

// Packager prepends written data with the size as uint64 (Big-Endian)
func Packager(w io.Writer) (chan<- []byte, chan<- error) {
	in := make(chan []byte, chanBufSize)
	err := make(chan error)
	go func() {
		size := make([]byte, 8)
		for buf := range in {
			binary.BigEndian.PutUint64(size, uint64(len(buf)))
			_, e := w.Write(size)
			if e != nil {
				err <- e
				return
			}
			_, e = w.Write(buf)
			if e != nil {
				err <- e
				return
			}
		}
	}()

	return in, err
}

// UnPackager is the reverse of Packager
func UnPackager(r io.Reader) <-chan []byte {
	out := make(chan []byte, chanBufSize)

	go func() {
		var err error
		size := make([]byte, 8)
		for {
			_, err = io.ReadFull(r, size)
			if err != nil {
				break
			}
			toRead := binary.BigEndian.Uint64(size)
			buf := make([]byte, toRead)
			_, err = io.ReadFull(r, buf)
			if err != nil {
				break
			}
			out <- buf
		}
		if err != io.EOF {
			panic(err)
		}
		close(out)
	}()

	return out
}
