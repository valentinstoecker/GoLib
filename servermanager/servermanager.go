package servermanager

import (
	"fmt"
	"net"
	"os/exec"

	"github.com/valentinstoecker/GoLib/packager"
)

type queueElem struct {
	data []byte
	next *queueElem
}

type queue struct {
	head *queueElem
	tail *queueElem
}

func (q *queue) Append(buf []byte) {
	newEl := &queueElem{
		data: buf,
		next: nil,
	}
	if q.head == nil {
		q.head = newEl
	}
	q.tail = newEl
}

func (q *queue) Pop() []byte {
	if q.head == nil {
		return nil
	}
	buf := q.head.data
	if q.head == q.tail {
		q.tail = nil
	}
	q.head = q.head.next
	return buf
}

func NewServer() error {
	listener, err := net.ListenTCP("tcp4", &net.TCPAddr{Port: 1337})
	if err != nil {
		return err
	}
	cmd := exec.Command("/bin/bash")
	cmdMap := make(map[string]*queue)
	q := &queue{}
	cmdMap[""] = q
	out, _ := cmd.StdoutPipe()
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := out.Read(buf)
			if err != nil {
				break
			}
			if n > 0 {
				q.Append(buf[:n])
			}
			buf = make([]byte, 4096)
		}
	}()
	inP, _ := cmd.StdinPipe()
	cmd.Start()
	for {
		con, err := listener.AcceptTCP()
		if err != nil {
			continue
		}
		in := packager.UnPackager(con)
		out, _ := packager.Packager(con)
		go func() {
			selected := ""
			for {
				if q, ok := cmdMap[selected]; ok {
					buf := q.Pop()
					if buf == nil {
						continue
					}
					out <- buf
				}
			}
		}()
		for inV := range in {
			inP.Write(inV)
		}
	}
}

func NewClient() (chan<- []byte, <-chan []byte) {
	con, err := net.Dial("tcp4", "127.0.0.1:1337")
	if err != nil {
		fmt.Println(err)
	}
	in, _ := packager.Packager(con)
	out := packager.UnPackager(con)
	return in, out
}
