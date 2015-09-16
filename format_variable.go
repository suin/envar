package main

type FormatVariable struct {
	Name  string
	Type  VariableType
	Value interface{}
}

type FormatVariables []FormatVariable

func (p FormatVariables) Len() int {
	return len(p)
}

func (p FormatVariables) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p FormatVariables) Less(i, j int) bool {
	return p[i].Name < p[j].Name
}
