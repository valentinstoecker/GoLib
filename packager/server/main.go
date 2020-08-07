package main

import (
	"fmt"
	"net"

	"github.com/valentinstoecker/GoLib/packager"
)

func handleCon(con net.Conn) {
	fmt.Println("Connection")
	defer fmt.Println("Done")

	in := packager.UnPackager(con)
	for range /*buf :=*/ in {
		//fmt.Println(buf)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":420")
	if err != nil {
		panic(err)
	}

	for {
		con, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		go handleCon(con)
	}
}
