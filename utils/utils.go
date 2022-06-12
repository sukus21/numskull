package utils

import (
	"fmt"
	"strconv"
)

//Converts a slice of bytes to a float64
func BytesliceToNumber(numData []byte) (float64, error) {

	//Empty string
	if len(numData) == 0 {
		return 0, fmt.Errorf("number has no length")
	}

	//Only a - sign
	if len(numData) == 1 && numData[0] == '-' {
		return 0, fmt.Errorf("invalid number, just a - sign")
	}

	//Unclosed comma
	lastChar := numData[len(numData)-1]
	if lastChar == ',' || lastChar == '.' {
		return 0, fmt.Errorf("comma with no value at the end")
	}

	//Convert to number
	number, err := strconv.ParseFloat(string(numData), 64)
	if err != nil {
		return 0, err
	}

	//Conversion succesful :)
	return number, nil
}
