package types

type NcBasicType struct {
	Name string
	Size int
}

func (t *NcBasicType) TypeDef() string {
	return t.Name
}
