package main

import (
	"encoding/binary"
	"net"
	"time"

	"github.com/valentinstoecker/GoLib/packager"
)

func main() {
	con, err := net.Dial("tcp", "192.168.188.30:420")
	if err != nil {
		panic(err)
	}

	in, _ := packager.Packager(con)
	for i := 0; i < 1000; i++ {
		buf := make([]byte, 256*256*256)
		binary.LittleEndian.PutUint64(buf, uint64(i))
		in <- buf
	}
	close(in)
	time.Sleep(time.Second)
}
