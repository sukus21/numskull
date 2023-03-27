package utils

import (
	"strconv"
	"strings"
)

// Converts a slice of bytes to a float64
func BytesliceToNumber(numData []byte) (float64, error) {
	str := strings.Trim(string(numData), " \t\n\r")
	return strconv.ParseFloat(string(str), 64)
}
