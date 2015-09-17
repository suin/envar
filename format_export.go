package main

import (
	"fmt"
)

func formatExport(environmentName string, variables []FormatVariable) string {
	result := fmt.Sprintf("# environment: %s\n", environmentName)
	for _, variable := range variables {
		var value string
		switch variable.Type {
		case Null:
			value = ""
		case String:
			value = quoteString(variable.Value)
		default:
			value = fmt.Sprintf("%v", variable.Value)
		}
		result += fmt.Sprintf("export %s=%s\n", variable.Name, value)
	}
	return result
}
