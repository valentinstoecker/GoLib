package types

import "fmt"

type NcNamedType struct {
	Name string
	Type NcType
}

func (t *NcNamedType) TypeDef() string {
	return t.Name
}

func (t *NcNamedType) Code() string {
	return fmt.Sprintf("type %s %s", t.Name, t.Type.TypeDef())
}
