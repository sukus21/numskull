package parser

import (
	"fmt"
	"numskull/token"
	"numskull/utils"
)

//Error things
type parserError_t struct {
	msg string
}

//Error message
func (err parserError_t) Error() string {
	return err.msg
}

type programContext struct {
	jumplinepos         int
	jumplinedestination int
}

//Preprocess program
func ParseProgram(raw string) ([]float64, bool) {

	//Different channels I need
	lines := make(chan string)
	Tokens := make(chan []float64)
	errors := make(chan error)

	//Actually preprocess program
	go logErrors(errors)
	go programSeperate(raw, lines, errors)
	go TokenizeLines(lines, Tokens, errors)
	return validateTokens(Tokens, errors)
}

//Seperate program per line
func programSeperate(program string, lines chan<- string, errors chan<- error) {

	currentLine := ""
	var prevChar byte = 0

	//thing
	for i := 0; i < len(program); i++ {

		char := program[i]

		//Ignore these
		if char == '\r' {
			continue
		}

		//End line
		if (char == '/' || char == '*') && prevChar == '/' || char == '\n' {

			//Was this a comment?
			if char == '/' && prevChar == '/' {
				currentLine = currentLine[:len(currentLine)-1]

				//Skip to next newline
				for char != '\n' {
					char = program[i]
					i++
				}
				i--
			}

			//Multiline comment?
			if char == '*' && prevChar == '/' {

				//Place a space
				currentLine = currentLine[:len(currentLine)-1]
				currentLine += " "

				//Wait for comment closure
				prevChar = 0
				for i++; i < len(program); i++ {
					char = program[i]

					//End of comment?
					if prevChar == '*' && char == '/' {
						break
					}

					//Newline?
					if char == '\n' {
						lines <- currentLine
						currentLine = ""
					}

					prevChar = char
				}

				continue
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
func TokenizeLines(lines <-chan string, Tokens chan<- []float64, errors chan<- error) {

	//Repeat for as long as there are lines
	linecount := 0
	for msg := range lines {
		linecount++
		pos := 0
		line := make([]float64, 0, 32)
		for {

			//Read Token
			tok, num, err := readToken(msg, &pos, linecount)
			if err != nil {
				fmt.Println(err.Error())
			}

			//Add Token to slice
			line = append(line, float64(tok))

			//Newline, end of current line
			if tok == token.Newline {
				break
			}

			//Was this a number?
			if tok == token.Number {
				line = append(line, num)
			}
		}

		//Send this slice through channel
		Tokens <- line
	}

	//Alright, we done
	close(Tokens)
}

//Make sure this stuff is valid code, and construct finished program
func validateTokens(Tokens <-chan []float64, errors chan<- error) ([]float64, bool) {

	//Finished program
	program := make([]float64, 0, 1024)
	curlies := make([]programContext, 0, 64)
	squares := make([]programContext, 0, 64)

	success := true
	expectBracket := false
	lastLine := 0
	linecount := 0
	lineStart := -1

	for toks := range Tokens {
		linecount++

		//Is there anything on this line
		if len(toks) == 0 || toks[0] == float64(token.Newline) {
			continue
		}

		//Expecting start bracket?
		tok := token.Token(toks[0])
		if expectBracket {
			expectBracket = false
			if tok != token.CurlyStart && tok != token.SquareStart {
				var err parserError_t
				err.msg = fmt.Sprintf("Line %d-%d: Expected starting bracket, got %s.\n", lastLine, linecount, tok.GetTokenName())
				errors <- err
				success = false
				continue
			}
		}

		//Is this an end bracket?
		if tok == token.CurlyEnd || tok == token.SquareEnd {
			next := token.Token(toks[1])

			//Expect newline
			if len(toks) > 2 || next != token.Newline {
				var err parserError_t
				err.msg = fmt.Sprintf("Line %d: Expected newline, got %s.\n", linecount, next.GetTokenName())
				errors <- err
				success = false
				continue
			}

			//Pop thing off the relevant stack
			var cnts *[]programContext
			if tok == token.CurlyEnd {
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
			if tok == token.SquareEnd {
				program = append(program, float64(token.SquareEnd), float64(cnt.jumplinepos))
			}
			program[cnt.jumplinedestination] = float64(len(program))
			continue
		}

		//Then this should be a number
		if tok != token.Number {
			var err parserError_t
			err.msg = fmt.Sprintf("Line %d: Expected number, got %s.\n", linecount, tok.GetTokenName())
			errors <- err
			success = false
			continue
		}

		//This IS a number
		lineStart = len(program)
		program = append(program, toks[0], toks[1])
		pos := 2

		for stayIn := true; stayIn; {
			tok = token.Token(toks[pos])
			pos++
			switch tok {

			//Unexpected newline
			case token.Newline:
				var err parserError_t
				err.msg = fmt.Sprintf("Line %d: Unexpected end of line.", linecount)
				errors <- err
				stayIn = false
				success = false

			//Lefthand chaining
			case token.ChainPlus, token.ChainMinus:

				//Read new Token
				next := token.Token(toks[pos])
				pos++
				if next != token.Number {
					var err parserError_t
					err.msg = fmt.Sprintf("Line %d: Expected number, got %s.\n", linecount, next.GetTokenName())
					errors <- err
					success = false
					stayIn = false
					break
				}

				//Push operand and number into program
				program = append(program, float64(tok), float64(token.Number), toks[pos])
				pos++
				continue

			//No righthand required
			case token.Decrement, token.Increment, token.PrintChar, token.PrintNumber, token.ReadInput:

				//Read next Token
				stayIn = false
				next := token.Token(toks[pos])
				pos++

				//Was this not a newline?
				if next != token.Newline {
					var err parserError_t
					err.msg = fmt.Sprintf("Line %d: Expected newline, got %s.\n", linecount, next.GetTokenName())
					errors <- err
					success = false
					break
				}

				//Write Token
				program = append(program, float64(tok))

			//Righthand required
			case token.Assign, token.Add, token.Sub, token.Multiply, token.Divide:

				//Read next Token
				stayIn = false
				next := token.Token(toks[pos])
				pos++

				//Expect number
				if next != token.Number {
					var err parserError_t
					err.msg = fmt.Sprintf("Line %d: Expected number, got %s.\n", linecount, next.GetTokenName())
					errors <- err
					success = false
					break
				}
				num := toks[pos]
				pos++

				//Expect newline
				next = token.Token(toks[pos])
				pos++
				if next != token.Newline {
					var err parserError_t
					err.msg = fmt.Sprintf("Line %d: Expected newline, got %s.\n", linecount, next.GetTokenName())
					errors <- err
					success = false
					break
				}

				//Push operand and number into program
				program = append(program, float64(tok), float64(token.Number), num)
				pos++

			//Expect righthand AND start bracket
			case token.Equals, token.Different, token.GreaterThan, token.GreaterEquals, token.LessThan, token.LessEquals:

				//Read next Token
				stayIn = false
				next := token.Token(toks[pos])
				pos++

				//Expect number
				if next != token.Number {
					var err parserError_t
					err.msg = fmt.Sprintf("Line %d: Expected number, got %s.\n", linecount, next.GetTokenName())
					errors <- err
					success = false
					break
				}
				num := toks[pos]
				pos++

				//Expect start bracket
				next = token.Token(toks[pos])
				pos++
				if next != token.CurlyStart && next != token.SquareStart {

					//Is it a newline?
					if next == token.Newline {
						//All is good
						expectBracket = true
						lastLine = linecount
					} else {
						var err parserError_t
						err.msg = fmt.Sprintf("Line %d: Expected newline, got %s.\n", linecount, next.GetTokenName())
						errors <- err
						success = false
						break
					}
				}

				//Push operand and number into program
				program = append(program, float64(tok), float64(token.Number), num, 0)
				pos++

				//Push current context according to bracket type
				if next == token.CurlyStart {
					curlies = append(curlies, programContext{
						jumplinepos:         lineStart,
						jumplinedestination: len(program) - 1,
					})
				} else if next == token.SquareStart {
					squares = append(squares, programContext{
						jumplinepos:         lineStart,
						jumplinedestination: len(program) - 1,
					})
				}

			//What on earth did you send me?
			default:
				stayIn = false
				var err parserError_t
				err.msg = fmt.Sprintf("Line %d: Unexpected %s found, expected operation.", linecount, tok.GetTokenName())
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
func readToken(text string, pos *int, line int) (token.Token, float64, error) {

	//Read a "word"
	word, done := readWord(text, pos)
	if done {
		return token.Newline, 0, nil
	}

	//Is this a number?
	num, err := utils.BytesliceToNumber([]byte(word))
	if err == nil {

		//Yes it is!
		return token.Number, num, nil
	}

	//Then what is it?
	switch word {

	//Arithmetic
	case "--":
		return token.Decrement, 0, nil
	case "++":
		return token.Increment, 0, nil
	case "+=":
		return token.Add, 0, nil
	case "-=":
		return token.Sub, 0, nil
	case "*=":
		return token.Multiply, 0, nil
	case "/=":
		return token.Divide, 0, nil

	//IO
	case "\"":
		return token.ReadInput, 0, nil
	case "!":
		return token.PrintNumber, 0, nil
	case "#":
		return token.PrintChar, 0, nil

	//Conditions
	case "?=":
		return token.Equals, 0, nil
	case "?!":
		return token.Different, 0, nil
	case "?>":
		return token.GreaterThan, 0, nil
	case "?>=":
		return token.GreaterEquals, 0, nil
	case "?<":
		return token.LessThan, 0, nil
	case "?<=":
		return token.LessEquals, 0, nil

	//Bracket
	case "{":
		return token.CurlyStart, 0, nil
	case "}":
		return token.CurlyEnd, 0, nil
	case "[":
		return token.SquareStart, 0, nil
	case "]":
		return token.SquareEnd, 0, nil

	//Others
	case "=":
		return token.Assign, 0, nil
	case "-":
		return token.ChainMinus, 0, nil
	case "+":
		return token.ChainPlus, 0, nil

	//Default
	default:
		return token.Invalid, 0, parserError_t{
			fmt.Sprintf("Line %d: Unknown operation: \"%s\"", line, word),
		}
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