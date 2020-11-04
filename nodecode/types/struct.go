package types

import "fmt"

type NcStruct []struct {
	Name string
	Data NcType
}

func (t *NcStruct) TypeDef() string {
	str := "struct{\n"
	for _, field := range *t {
		str += fmt.Sprintf("%s %s\n", field.Name, field.Data.TypeDef())
	}
	str += "}"
	return str
}
