package main

import (
	"fmt"
	"strings"
)

func formatDocker(environmentName string, variables []FormatVariable) string {
	arguments := []string{}
	for _, variable := range variables {
		var value string
		switch variable.Type {
		case Null:
			value = ""
		default:
			value = fmt.Sprintf("%v", variable.Value)
		}
		arguments = append(arguments, "-e "+quoteString(fmt.Sprintf("%s=%s", variable.Name, value)))
	}
	return strings.Join(arguments, " ")
}
