package main

import (
	"fmt"
	"strconv"
)

func formatEnvfile(environmentName string, variables []FormatVariable) string {
	result := fmt.Sprintf("# environment: %s\n", environmentName)
	for _, variable := range variables {
		b := []byte{}
		b = strconv.AppendQuote(b, fmt.Sprintf("%v", variable.Value))
		result += fmt.Sprintf("%s=%s\n", variable.Name, string(b))
	}
	return result
}
