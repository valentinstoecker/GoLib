package types

import "fmt"

type NcArray struct {
	Of   NcType
	Size int
}

func (t *NcArray) TypeDef() string {
	return fmt.Sprintf("[%d]%s", t.Size, t.Of.TypeDef())
}
