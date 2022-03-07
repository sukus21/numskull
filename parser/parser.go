package parser

import (
	"fmt"
	"numskull/utils"
)

const (
	token_invalid       token = iota
	token_number        token = iota
	token_chainPlus     token = iota
	token_chainMinus    token = iota
	token_assign        token = iota
	token_add           token = iota
	token_sub           token = iota
	token_multiply      token = iota
	token_divide        token = iota
	token_increment     token = iota
	token_decrement     token = iota
	token_printNumber   token = iota
	token_printChar     token = iota
	token_readInput     token = iota
	token_equals        token = iota
	token_different     token = iota
	token_lessThan      token = iota
	token_greaterThan   token = iota
	token_lessEquals    token = iota
	token_greaterEquals token = iota
	token_curlyStart    token = iota
	token_curlyEnd      token = iota
	token_squareStart   token = iota
	token_squareEnd     token = iota
	token_newline       token = iota
)

//Error things
type parserError_t struct {
	msg string
}

//Error message
func (err parserError_t) Error() string {
	return err.msg
}

//Token struct
type token int

func (tok token) getTokenName() string {
	switch tok {
	case token_add:
		return "addition"
	case token_sub:
		return "subtraction"
	case token_assign:
		return "assign"
	case token_multiply:
		return "multiplication"
	case token_divide:
		return "division"
	case token_increment:
		return "increment"
	case token_decrement:
		return "decrement"
	case token_chainPlus:
		return "chain plus"
	case token_chainMinus:
		return "chain minus"

	case token_equals:
		return "equals"
	case token_different:
		return "different"
	case token_lessEquals:
		return "less or equals"
	case token_lessThan:
		return "less than"
	case token_greaterEquals:
		return "greater or equals"
	case token_greaterThan:
		return "greater than"

	case token_number:
		return "number"
	case token_printChar:
		return "print char"
	case token_printNumber:
		return "print number"
	case token_readInput:
		return "input"

	case token_curlyStart:
		return "curly start"
	case token_curlyEnd:
		return "curly end"
	case token_squareStart:
		return "square start"
	case token_squareEnd:
		return "square end"
	case token_invalid:
		return "invalid"

	default:
		return "unknown"
	}
}

//Preprocess program
func ParseProgram(raw string) []float64 {

	//Variables
	lines := make(chan string)
	tokens := make(chan []float64)
	errors := make(chan error)

	//Actually preprocess program
	go programSeperate(raw, lines, errors)
	tokenizeLines(lines, tokens, errors)

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
		pos := 0
		for {
			tok, num, err := readToken(msg, &pos)
		}
	}
}

func validateTokens(tokens <-chan []float64, errors chan<- error) ([]float64, bool) {
	//Poop
	return nil, false
}

//Grab number
func readToken(text string, pos *int) (token, float64, error) {

	//Read a "word"
	word, done := readWord(text, pos)
	if done {
		return token_newline, 0, nil
	}

	//Is this a number?
	num, err := utils.BytesliceToNumber([]byte(word))
	if err == nil {

		//Yes it is!
		return token_number, num, nil
	}

	//Then what is it?
	switch word {

	//Arithmetic
	case "--":
		return token_decrement, 0, nil
	case "++":
		return token_increment, 0, nil
	case "+=":
		return token_add, 0, nil
	case "-=":
		return token_sub, 0, nil
	case "*=":
		return token_multiply, 0, nil
	case "/=":
		return token_divide, 0, nil

	//IO
	case "\"":
		return token_readInput, 0, nil
	case "!":
		return token_printNumber, 0, nil
	case "#":
		return token_printChar, 0, nil

	//Conditions
	case "?=":
		return token_equals, 0, nil
	case "?!":
		return token_different, 0, nil
	case "?>":
		return token_greaterThan, 0, nil
	case "?>=":
		return token_greaterEquals, 0, nil
	case "?<":
		return token_lessThan, 0, nil
	case "?<=":
		return token_lessEquals, 0, nil

	//Bracket
	case "{":
		return token_curlyStart, 0, nil
	case "}":
		return token_curlyEnd, 0, nil
	case "[":
		return token_squareStart, 0, nil
	case "]":
		return token_squareEnd, 0, nil

	//Others
	case "=":
		return token_assign, 0, nil
	case "-":
		return token_chainMinus, 0, nil
	case "+":
		return token_chainPlus, 0, nil

	//Default
	default:
		return token_invalid, 0, parserError_t{msg: "Unknown operation: \"" + word + "\""}
	}
}

func readWord(text string, pos *int) (string, bool) {
	char := getNextNonWhitespace(text, pos)
	if char == 0 {
		return "", true
	}
	output := string(char)

	//Handle this case
	if char == '-' {
		char = getNextChar(text, pos)
		output += string(char)
	}

	switch char {

	//Numeric?
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ',', '.':
		for {
			switch char = getNextChar(text, pos); char {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ',', '.':
				output += string(char)
			default:
				*pos--
				return output, false
			}
		}

	//Anything else
	default:
		for {
			switch char = getNextChar(text, pos); char {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ',', '.', 0, ' ', '\t', 0x0D, '\n':
				*pos--
				return output, false
			default:
				output += string(char)
			}
		}
	}
}

//Get the next program character that ISN'T whitespace
func getNextNonWhitespace(text string, pos *int) byte {
	for *pos < len(text) {
		if val := getNextChar(text, pos); !isWhitespace(val) {
			return val
		}
	}

	return 0
}

//Get the next character
func getNextChar(text string, pos *int) byte {
	if *pos == len(text) {
		return 0
	}

	char := text[*pos]
	*pos++
	return char
}

//Is the given character whitespace or not?
func isWhitespace(char byte) bool {

	//Is the character any of the given values?
	return char == ' ' || char == '\t' || char == '\n' || char == 0x0D
}
