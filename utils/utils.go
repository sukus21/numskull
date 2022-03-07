package utils

import "strconv"

type utilsError_t struct {
	msg string
}

func (err utilsError_t) Error() string {
	return err.msg
}

//Converts a slice of bytes to a float64
func BytesliceToNumber(numData []byte) (float64, error) {
	var err utilsError_t

	//Empty string
	if len(numData) == 0 {
		err.msg = "Number has no length."
		return 0, err
	}

	//Only a - sign
	if len(numData) == 1 && numData[0] == '-' {
		err.msg = "Invalid number, just a - sign."
		return 0, err
	}

	//Unclosed comma
	lastChar := numData[len(numData)-1]
	if lastChar == ',' || lastChar == '.' {
		err.msg = "Comma with no value at the end"
		return 0, err
	}

	//Convert to number
	number, errC := strconv.ParseFloat(string(numData), 64)
	if errC != nil {
		return 0, errC
	}

	//Conversion succesful :)
	return number, nil
}
