package main

import (
	"fmt"
	"strconv"
)

func quoteString(v interface{}) string {
	b := []byte{}
	b = strconv.AppendQuote(b, fmt.Sprintf("%v", v))
	return string(b)
}
