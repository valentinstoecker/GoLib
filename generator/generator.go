package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	err := generate(os.Args[1], os.Args[2])
	if err != nil {
		panic(err)
	}
}

func generate(path, typeName string) error {
	fmt.Println(typeName)
	in, err := os.Open(path)
	defer in.Close()
	if err != nil {
		return err
	}
	if files, err := in.Readdir(0); err == nil {
		for _, file := range files {
			if isTGo(file.Name()) && !file.IsDir() {
				in, err := os.Open(filepath.Join(path, file.Name()))
				if err != nil {
					return err
				}
				newName := strings.TrimSuffix(file.Name(), ".tgo") + typeName + ".go"
				out, err := os.Create(filepath.Join(path, newName))
				if err != nil {
					return err
				}
				generateFile(typeName, in, out)
				fmt.Println("+ " + filepath.Join(path, file.Name()))
			}
		}
		return nil
	}
	if !isTGo(path) {
		return errors.New(path + " is not a dir or .tgo file")
	}
	out, err := os.Create(strings.TrimSuffix(path, ".tgo") + typeName + ".go")
	if err != nil {
		return err
	}
	err = generateFile(typeName, in, out)
	if err != nil {
		return err
	}
	fmt.Println("+ " + path)
	return nil
}

func generateFile(typeName string, src io.Reader, dst io.Writer) error {
	var b strings.Builder
	io.Copy(&b, src)
	t, err := template.New("").Parse(b.String())
	if err != nil {
		return err
	}
	t.Execute(dst, struct {
		Type string
	}{
		typeName,
	})
	return nil
}

func isTGo(name string) bool {
	return strings.HasSuffix(name, ".tgo")
}
