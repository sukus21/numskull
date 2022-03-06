package parser

import (
	"fmt"
	"os"
)

const (
	token_invalid       = iota
	token_number        = iota
	token_chainPlus     = iota
	token_chainMinus    = iota
	token_assign        = iota
	token_add           = iota
	token_sub           = iota
	token_multiply      = iota
	token_divide        = iota
	token_increment     = iota
	token_decrement     = iota
	token_printNumber   = iota
	token_printChar     = iota
	token_readInput     = iota
	token_equals        = iota
	token_different     = iota
	token_lessThan      = iota
	token_greaterThan   = iota
	token_lessEquals    = iota
	token_greaterEquals = iota
	token_curlyStart    = iota
	token_curlyEnd      = iota
	token_squareStart   = iota
	token_squareEnd     = iota
)

//Error things
type parserError_t struct {
	msg string
}

//Error message
func (err numskullError_t) Error() string {
	return err.msg
}

//Token struct
type token int

//Preprocess program
func ParseProgram(raw string) []float64 {

	//Variables
	lines := make(chan string)
	tokens := make(chan []float64)
	errors := make(chan error)

	//Actually preprocess program
	go programSeperate(raw, lines, errors)
	go tokenizeLines(lines, tokens, errors)
	program, success := validateTokens(tokens, errors)

	//Check success state
	if !success {
		fmt.Println("uh oh, something went wrong :(")
	}

	//Return program
	return program
}

//Seperate program per line
func programSeperate(program string, lines chan<- string, errors chan<- error) {

	currentLine := ""
	var prevChar byte = 0

	//thing
	for i := 0; i < len(program); i++ {

		char := program[i]

		//Ignore these
		if char == 0x0D {
			continue
		}

		//End line
		if (char == '/' && prevChar == '/') || char == '\n' {

			//Was this a comment?
			if char == '/' {
				currentLine = currentLine[:len(currentLine)-1]
			}

			//Push line into output channel
			lines <- currentLine

			//Reset variables and continue loop
			currentLine = ""
			prevChar = 0
			continue
		}

		//Add character to current line
		currentLine += string(char)
		prevChar = char
	}

	//Add last line
	if len(currentLine) != 0 {
		lines <- currentLine
	}

	//Close channel
	close(lines)
}

//Handle lines
func tokenizeLines(lines <-chan string, tokens chan<- []float64, errors chan<- error) {

	//Repeat for as long as there are lines
	for msg := range lines {

		//Expect a number (or nothing)
		num, tok, err := readExpression(msg)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println(num, tok)
	}
}

func validateTokens(tokens <-chan []float64, errors chan<- error) ([]float64, bool) {
	//Poop
	return nil, false
}

//Grab number
func readExpression(text string, readPos *int) (float64, token, error) {

	expressionData := make([]byte, 512)[:0]
	var foundNumber bool = false
	var foundComma bool = false

	for breakLoop := false; *readPos < len(text) && !breakLoop; {

		//Switch per character
		switch char := getNextChar(); char {

		//Regular numbers
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			expressionData = append(expressionData, char)
			foundNumber = true

		//Negative stuff
		case '-':
			if !foundNumber {
				expressionData = append(expressionData, char)
			} else {
				breakLoop = true
				*readPos--
			}

		//Signify that this is a decimal number
		case '.', ',':

			//Append a 0 to the front
			if !foundNumber {
				expressionData = append(expressionData, '0')
				foundNumber = true
			}

			//Was another comma already seen?
			if foundComma {
				var err numskullError_t
				err.msg = "Double commas in number declaration."
				return 0, token_invalid, err
			}

			//Set sum variables
			expressionData = append(expressionData, '.')
			foundComma = true

		//Unknown character, throw an error
		default:

			//Straight up ignore this
			if char == '}' {
				break
			}

			//Jump maybe
			if char == ']' {
				if !foundNumber {

					//Find matching opening bracket
					*readPos -= 2
					for depth := 1; depth != 0 && *readPos >= 0; {
						char := getPreviousNonWhitespace()
						if char == ']' {
							depth++
						} else if char == '[' {
							depth--
						}
					}

					//Is there not an opening bracket?
					if *readPos >= len(text) {
						var err numskullError_t
						err.msg = "Expected an opening loop bracket '['. Never found one."
						return 0, token_invalid, err
					}

					//Jump to start of line
					char := byte(0)
					for {
						char = getPreviousChar()
						if isNewline(char) || char == 0 {
							break
						}
					}

					//Here we are
					break

				} else {

					//Unexpected character encountered
					var err numskullError_t
					err.msg = "Expected character encountered, '" + string(char) + "'."
					return 0, token_invalid, err
				}
			}

			//Ignore
			if isWhitespace(char) {
				if foundNumber {
					breakLoop = true
					break
				}
				break
			}

			//We done readin bois
			*readPos--
			breakLoop = true
		}
	}

	//Convert to number
	number, err := bytesliceToNumber(numberData)
	if err != nil {
		return 0, token_invalid, nil
	}

	//Not supposed to look at memory, return number itself
	return number, token_number, nil
}
