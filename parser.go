package main

import (
	"fmt"
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
type token float64

func (tok token) getTokenName() string {
	switch tok {
	case token_add:
		return "+="
	case token_sub:
		return "-="
	case token_assign:
		return "="
	case token_multiply:
		return "*="
	case token_divide:
		return "/="
	case token_increment:
		return "++"
	case token_decrement:
		return "--"
	case token_chainPlus:
		return "+"
	case token_chainMinus:
		return "-"

	case token_equals:
		return "?="
	case token_different:
		return "?!"
	case token_lessEquals:
		return "?<="
	case token_lessThan:
		return "?<"
	case token_greaterEquals:
		return "?>="
	case token_greaterThan:
		return "?>"

	case token_number:
		return "number"
	case token_printChar:
		return "#"
	case token_printNumber:
		return "!"
	case token_readInput:
		return "\""

	case token_curlyStart:
		return "{"
	case token_curlyEnd:
		return "}"
	case token_squareStart:
		return "["
	case token_squareEnd:
		return "]"
	case token_newline:
		return "newline"
	case token_invalid:
		return "invalid"

	default:
		return "unknown"
	}
}

type programContext struct {
	jumplinepos         int
	jumplinedestination int
}

//Preprocess program
func ParseProgram(raw string) ([]float64, bool) {

	//Different channels I need
	lines := make(chan string)
	tokens := make(chan []float64)
	errors := make(chan error)

	//Actually preprocess program
	go logErrors(errors)
	go programSeperate(raw, lines, errors)
	go tokenizeLines(lines, tokens, errors)
	return validateTokens(tokens, errors)
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

				//Skip to next newline
				for char != '\n' {
					char = program[i]
					i++
				}
				i--
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
		pos := 0
		line := make([]float64, 0, 16)
		for {

			//Read token
			tok, num, err := readToken(msg, &pos)
			if err != nil {
				fmt.Println(err.Error())
			}

			//Add token to slice
			line = append(line, float64(tok))

			//Newline, end of current line
			if tok == token_newline {
				break
			}

			//Was this a number?
			if tok == token_number {
				line = append(line, num)
			}
		}

		//Send this slice through channel
		tokens <- line
	}

	//Alright, we done
	close(tokens)
}

//Make sure this stuff is valid code, and construct finished program
func validateTokens(tokens <-chan []float64, errors chan<- error) ([]float64, bool) {

	//Finished program
	program := make([]float64, 0, 1024)
	curlies := make([]programContext, 0, 64)
	squares := make([]programContext, 0, 64)

	success := true
	expectBracket := false
	lastLine := 0
	linecount := 0
	lineStart := -1

	for toks := range tokens {
		linecount++

		//Is there anything on this line
		if len(toks) == 0 || toks[0] == float64(token_newline) {
			continue
		}

		//Expecting start bracket?
		tok := token(toks[0])
		if expectBracket {
			expectBracket = false
			if tok != token_curlyStart && tok != token_squareStart {
				var err parserError_t
				err.msg = fmt.Sprintf("Line %d-%d: Expected starting bracket, got %s.\n", lastLine, linecount, tok.getTokenName())
				errors <- err
				success = false
				continue
			}
		}

		//Is this an end bracket?
		if tok == token_curlyEnd || tok == token_squareEnd {
			next := token(toks[1])

			//Expect newline
			if len(toks) > 2 || next != token_newline {
				var err parserError_t
				err.msg = fmt.Sprintf("Line %d: Expected newline, got %s.\n", linecount, next.getTokenName())
				errors <- err
				success = false
				continue
			}

			//Pop thing off the relevant stack
			var cnts *[]programContext
			if tok == token_curlyEnd {
				cnts = &curlies
			} else {
				cnts = &squares
			}

			//Check if stack is empty
			if len(*cnts) == 0 {
				var err parserError_t
				err.msg = fmt.Sprintf("Line %d: Unmatched }.\n", linecount)
				errors <- err
				success = false
				continue
			}

			//It is not
			cnt := []programContext(*cnts)[len(*cnts)-1]
			*cnts = []programContext(*cnts)[:len(*cnts)-1]

			//Was is a looping bracket?
			if tok == token_squareEnd {
				program = append(program, float64(token_squareEnd), float64(cnt.jumplinepos))
			}
			program[cnt.jumplinedestination] = float64(len(program))
			continue
		}

		//Then this should be a number
		if tok != token_number {
			var err parserError_t
			err.msg = fmt.Sprintf("Line %d: Expected number, got %s.\n", linecount, tok.getTokenName())
			errors <- err
			success = false
			continue
		}

		//This IS a number
		lineStart = len(program)
		program = append(program, toks[0], toks[1])
		pos := 2

		for stayIn := true; stayIn; {
			tok = token(toks[pos])
			pos++
			switch tok {

			//Unexpected newline
			case token_newline:
				var err parserError_t
				err.msg = fmt.Sprintf("Line %d: Unexpected end of line.", linecount)
				errors <- err
				stayIn = false
				success = false

			//Lefthand chaining
			case token_chainPlus, token_chainMinus:

				//Read new token
				next := token(toks[pos])
				pos++
				if next != token_number {
					var err parserError_t
					err.msg = fmt.Sprintf("Line %d: Expected number, got %s.\n", linecount, next.getTokenName())
					errors <- err
					success = false
					stayIn = false
					break
				}

				//Push operand and number into program
				program = append(program, float64(tok), float64(token_number), toks[pos])
				pos++
				continue

			//No righthand required
			case token_decrement, token_increment, token_printChar, token_printNumber, token_readInput:

				//Read next token
				stayIn = false
				next := token(toks[pos])
				pos++

				//Was this not a newline?
				if next != token_newline {
					var err parserError_t
					err.msg = fmt.Sprintf("Line %d: Expected newline, got %s.\n", linecount, next.getTokenName())
					errors <- err
					success = false
					break
				}

				//Write token
				program = append(program, float64(tok))

			//Righthand required
			case token_assign, token_add, token_sub, token_multiply, token_divide:

				//Read next token
				stayIn = false
				next := token(toks[pos])
				pos++

				//Expect number
				if next != token_number {
					var err parserError_t
					err.msg = fmt.Sprintf("Line %d: Expected number, got %s.\n", linecount, next.getTokenName())
					errors <- err
					success = false
					break
				}
				num := toks[pos]
				pos++

				//Expect newline
				next = token(toks[pos])
				pos++
				if next != token_newline {
					var err parserError_t
					err.msg = fmt.Sprintf("Line %d: Expected newline, got %s.\n", linecount, next.getTokenName())
					errors <- err
					success = false
					break
				}

				//Push operand and number into program
				program = append(program, float64(tok), float64(token_number), num)
				pos++

			//Expect righthand AND start bracket
			case token_equals, token_different, token_greaterThan, token_greaterEquals, token_lessThan, token_lessEquals:

				//Read next token
				stayIn = false
				next := token(toks[pos])
				pos++

				//Expect number
				if next != token_number {
					var err parserError_t
					err.msg = fmt.Sprintf("Line %d: Expected number, got %s.\n", linecount, next.getTokenName())
					errors <- err
					success = false
					break
				}
				num := toks[pos]
				pos++

				//Expect start bracket
				next = token(toks[pos])
				pos++
				if next != token_curlyStart && next != token_squareStart {

					//Is it a newline?
					if next == token_newline {
						//All is good
						expectBracket = true
						lastLine = linecount
					} else {
						var err parserError_t
						err.msg = fmt.Sprintf("Line %d: Expected newline, got %s.\n", linecount, next.getTokenName())
						errors <- err
						success = false
						break
					}
				}

				//Push operand and number into program
				program = append(program, float64(tok), float64(token_number), num, 0)
				pos++

				//Push current context according to bracket type
				if next == token_curlyStart {
					curlies = append(curlies, programContext{
						jumplinepos:         lineStart,
						jumplinedestination: len(program) - 1,
					})
				} else if next == token_squareStart {
					squares = append(squares, programContext{
						jumplinepos:         lineStart,
						jumplinedestination: len(program) - 1,
					})
				}

			//What on earth did you send me?
			default:
				stayIn = false
				var err parserError_t
				err.msg = fmt.Sprintf("Line %d: Unexpected %s found, expected operation.", linecount, tok.getTokenName())
				success = false
				errors <- err
			}
		}
	}

	//We done :)
	close(errors)
	return program, success
}

//Yup
func logErrors(errors <-chan error) {
	for err := range errors {
		fmt.Println(err.Error())
	}
}

//Grab number
func readToken(text string, pos *int) (token, float64, error) {

	//Read a "word"
	word, done := readWord(text, pos)
	if done {
		return token_newline, 0, nil
	}

	//Is this a number?
	num, err := BytesliceToNumber([]byte(word))
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

//Um
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
			case 0:
				return output, false
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
				if char != 0 {
					*pos--
				}
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
