package main

import (
	"fmt"
	"strconv"
	"strings"
)

func formatDocker(environmentName string, variables []FormatVariable) string {
	arguments := []string{}
	for _, variable := range variables {
		var value string
		switch variable.Type {
		case Null:
			value = ""
		case String:
			b := []byte{}
			b = strconv.AppendQuote(b, fmt.Sprintf("%v", variable.Value))
			value = string(b)
		default:
			value = fmt.Sprintf("%v", variable.Value)
		}
		arguments = append(arguments, fmt.Sprintf("-e %s=%s", variable.Name, value))
	}
	return strings.Join(arguments, " ")
}
