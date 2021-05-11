package main

import (
	"fmt"

	"github.com/valentinstoecker/GoLib/servermanager"
)

func main() {
	in, out := servermanager.NewClient()
	done := make(chan interface{})
	go func() {
		for res := range out {
			fmt.Printf("%s", string(res))
		}
		done <- nil
	}()
	var str string
	for {
		_, err := fmt.Scanln(&str)
		if err != nil {
			fmt.Println()
			break
		}
		str += "\n"
		fmt.Println("---------")
		fmt.Print(str)
		fmt.Println("---------")
		in <- []byte(str)
	}
	<-done
}
