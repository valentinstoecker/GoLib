package types

import "fmt"

type NcPointer struct {
	To NcType
}

func (t *NcPointer) TypeDef() string {
	return fmt.Sprintf("*%s", t.To.TypeDef())
}
