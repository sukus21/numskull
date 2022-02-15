package main

import (
	"fmt"
	"os"
	"strconv"
)

const (
	usage_v string = "-v, --version         Prints program version number"
	usage_h string = "-h, --help <argument> Prints usage for given argument"
	usage_i string = "-i, --input <path>    File to read input from"
	usage_o string = "-o, --output <path>   File to print output to"
	usage_t string = "-t, --type            Tells program to read input file as text"
	usage_c string = "-c, --console         Force program output to console"
)

var program []byte
var readPos int
var memory map[float64]float64 = make(map[float64]float64)

var input []float64
var inputPos int

var consoleOutput bool = true
var readFromFile bool = false
var inputBinary bool = true
var writeToFile bool = false
var outputFile *os.File = nil

//Entrypoint, reads command line arguments
func main() {

	//Uncomment only for debugging
	//os.Args = []string{os.Args[0], "-i", "hello.bf", "-o", "out.txt", "program.txt"}

	//No arguments provided
	fmt.Println()
	if len(os.Args) == 1 {

		printUsage()
		closeProgram(1)
	}

	//Variables
	var forceConsole bool = false
	var inputRaw []byte

	//Read and process arguments
	var somethingDone bool = false
	for argPos := 1; argPos < len(os.Args); argPos++ {
		switch os.Args[argPos] {

		//Help :)
		case "-h", "-H", "--help":
			argPos++

			//No help thing provided
			if argPos == len(os.Args) {
				fmt.Println("Error: specify an argument to explain.")
				fmt.Println("Example:", os.Args[0], os.Args[argPos-1], "v")
				fmt.Println("Pass in no arguments to see list of available arguments.")
				closeProgram(0)
			}

			switch os.Args[argPos] {

			//Help for the version tag
			case "v", "V", "version":
				fmt.Println(usage_v)
				fmt.Println("That's it really.")

			//Help for the help tag
			case "h", "H", "help":
				fmt.Println(usage_h)
				fmt.Println()
				fmt.Println("Example: numskull", "--"+os.Args[argPos], "v")
				fmt.Println("Shows the help page for the parameter '-v'.")
				fmt.Println("The program stops once help has been printed.")

			//Help for the input tag
			case "i", "I", "input":
				fmt.Println(usage_i)
				if len(os.Args[argPos]) != 1 {
					os.Args[argPos] = "-" + os.Args[argPos]
				}
				fmt.Println("When reading after the end of file, the result is always -1.")
				fmt.Println("If no file is specified, input is given through the console.")
				fmt.Println("Input file will be read as binary by default, look up \"numskull --help -t\" for more info.")
				fmt.Println()
				fmt.Println("Example: numskull", "-"+os.Args[argPos], "\"C:\\numbers.bin\" \"C:\\program.nms\"")
				fmt.Println("Opens program.nms, and reads from numbers.bin when reading input.")

			//Help for the output tag
			case "o", "O", "output":
				if len(os.Args[argPos]) != 1 {
					os.Args[argPos] = "-" + os.Args[argPos]
				}
				fmt.Println(usage_o)
				fmt.Println("When outputting, you can choose to also write that output to a file.")
				fmt.Println("The file is treated as a byte array.")
				fmt.Println("If no file is specified, the output of the program is displayed in the console.")
				fmt.Println("If you still want console output AND saving to a file, use the \"-c\" argument.")
				fmt.Println("If the program stops due to an error, the file is still saved.")
				fmt.Println()
				fmt.Println("Example: numskull", "-"+os.Args[argPos], "\"C:\\numbers.bin\" \"C:\\program.nms\"")
				fmt.Println("Opens program.nms, and saves output to numbers.bin, once the program stops running.")

			//Help for the type tag
			case "t", "T", "type":
				if len(os.Args[argPos]) != 1 {
					os.Args[argPos] = "-" + os.Args[argPos]
				}
				fmt.Println(usage_t)
				fmt.Println("Normally when reading input from a file, the input will be read as binary.")
				fmt.Println("That is, one byte per input, each consisting of a number from 0 to 255.")
				fmt.Println("Some may not want this behaviour, so passing in", "-"+os.Args[argPos], "will make input read as text.")
				fmt.Println("Entries are read as numbers, seperated by whitespace (tabs, spaces, or newlines).")
				fmt.Println("An incorrectly formatted entry will cause an error when trying to read it.")
				fmt.Println("If the -i argument isn't present, this argument does nothing.")

			//Help for the console tag
			case "c", "C", "console":
				fmt.Println(usage_c)
				fmt.Println("When outputting to a file, console output is turned off by default.")
				fmt.Println("Use this tag to reenable it, while also writing the output to a file, using -o.")
				fmt.Println("If the -o argument isn't present, this argument does nothing.")

			default:
				fmt.Println("Error: unknown argument.")
				fmt.Println()
				printUsage()
				closeProgram(1)
			}

			//Exit
			closeProgram(0)

		//Print numskull version
		case "-v", "-V", "--version":
			fmt.Println("Numskull interpreter v0.1.0")
			fmt.Println("Supports Numskull version 1.0")
			somethingDone = true

		//Force console output
		case "-c", "-C", "--console":
			forceConsole = true

		//Text reading mode
		case "-t", "-T", "--type":
			inputBinary = false

		//Specify output file
		case "-o", "-O", "--output":
			argPos++

			//No file specified
			if argPos >= len(os.Args)-1 {
				fmt.Println("Error: no output file specified")
				fmt.Println(usage_o)
				break
			}

			//Update stuff
			var err error
			outputFile, err = os.Create(os.Args[argPos])
			if err != nil {
				fmt.Println("Error creating output file")
				fmt.Println(err.Error())
				closeProgram(1)
			}
			consoleOutput = false
			writeToFile = true

		//Read input file
		case "-i", "-I", "--input":
			argPos++

			//No file specified
			if argPos >= len(os.Args)-1 {
				fmt.Println("Error: no input file specified")
				fmt.Println(usage_i)
				break
			}

			//Open file
			raw, err := os.ReadFile(os.Args[argPos])
			if err != nil {
				fmt.Println("Error while opening input file")
				fmt.Println(err.Error())
				fmt.Println(usage_i)
				closeProgram(1)
			}

			//Read thing
			inputRaw = raw
			readFromFile = true

		//Unknown parameter
		default:
			if argPos == len(os.Args)-1 {
				break
			}

			fmt.Println("Error: Unknown argument: " + os.Args[argPos])
			printUsage()
			closeProgram(1)
		}
	}

	//No more arguments? Was anything achieved?
	finArg := os.Args[len(os.Args)-1]
	if finArg[0] == '-' && somethingDone {
		closeProgram(0)
	}

	//Try to read program file
	file, err := os.ReadFile(finArg)
	if err != nil {
		fmt.Println("Error opening program file")
		fmt.Println(err.Error())
		closeProgram(1)
	}
	program = file

	//Read input file
	if readFromFile {

		//Convert input from text to []float64
		if !inputBinary {
			input = make([]float64, len(inputRaw)/5)[:0]
			pos := 0
			for {
				numData := make([]byte, 64)[:0]
				foundChar := false
				foundComma := false
				foundWhitespace := false

				//Read ONE number
				for {
					if pos == len(inputRaw) || foundWhitespace {
						break
					}

					//Read char
					char := inputRaw[pos]
					pos++

					switch char {

					//Numbers
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						foundChar = true
						numData = append(numData, char)

					//Decimal indicators
					case '.', ',':
						if !foundChar {
							numData = append(numData, '0')
							foundChar = true
						}
						if foundComma {
							fmt.Println("Error converting number")
							fmt.Println("Double commas on entry", len(input)+1, ".")
							closeProgram(1)
						}
						foundComma = true
						numData = append(numData, '.')

					//Minus sign
					case '-':
						if foundChar {
							fmt.Println("Error converting number")
							fmt.Println("Unexpected character encountered '", string(char), "' on entry", len(input)+1)
							closeProgram(1)
						}
						foundChar = true
						numData = append(numData, char)

					//Whitespace
					case ' ', 0x0D, '\t', '\n':
						if !foundChar {
							continue
						} else {
							foundWhitespace = true
						}

					//Unknown character
					default:
						fmt.Println("Error converting number")
						fmt.Println("Unexpected character encountered '", string(char), "' on entry", len(input)+1)
						closeProgram(1)
					}
				}

				//Convert to number
				number, err := bytesliceToNumber(numData)
				if err != nil {
					fmt.Println("Error when converting input file")
					fmt.Println(err.Error())
					closeProgram(1)
				}
				input = append(input, number)

				//Is this all?
				if pos == len(inputRaw) {
					break
				}
			}
		}

		//Convert from []byte to []float64
		if inputBinary {
			input = make([]float64, len(inputRaw))
			for i := 0; i < len(inputRaw); i++ {
				input[i] = float64(inputRaw[i])
			}
		}
	}

	//Console?
	if forceConsole {
		consoleOutput = true
	}

	//Start executing it
	err = runProgram()
	fmt.Println()
	fmt.Println()
	if err != nil {

		//Print error
		fmt.Println(err.Error())
	} else {

		//Program finished :)
		fmt.Println("Program finished :)")
	}

	//Close file maybe
	if writeToFile {
		outputFile.Close()
	}
}

//Error things
type numskullError_t struct {
	msg string
}

//Error message
func (err numskullError_t) Error() string {
	return err.msg
}

//Grab number
func getExpression(memoryLookup bool) (float64, error) {

	numberData := make([]byte, 512)[:0]
	var foundChar bool = false
	var foundComma bool = false

	for breakLoop := false; readPos < len(program) && !breakLoop; {

		//Switch per character
		switch char := getNextChar(); char {

		//Regular numbers
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			numberData = append(numberData, char)
			foundChar = true

		//Negative stuff
		case '-':
			if !foundChar {
				numberData = append(numberData, char)
			} else {
				breakLoop = true
				readPos--
			}

		//Signify that this is a decimal number
		case '.', ',':

			//Append a 0 to the front
			if !foundChar {
				numberData = append(numberData, '0')
				foundChar = true
			}

			//Was another comma already seen?
			if foundComma {
				var err numskullError_t
				err.msg = "Double commas in number declaration."
				return 0, err
			}

			//Set sum variables
			numberData = append(numberData, '.')
			foundComma = true

		//Unknown character, throw an error
		default:

			//Straight up ignore this
			if char == '}' {
				break
			}

			//Jump maybe
			if char == ']' {
				if !foundChar {

					//Find matching opening bracket
					readPos -= 2
					for depth := 1; depth != 0 && readPos >= 0; {
						char := getPreviousNonWhitespace()
						if char == ']' {
							depth++
						} else if char == '[' {
							depth--
						}
					}

					//Is there not an opening bracket?
					if readPos >= len(program) {
						var err numskullError_t
						err.msg = "Expected an opening loop bracket '['. Never found one."
						return 0, err
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
					return 0, err
				}
			}

			//Ignore
			if isWhitespace(char) {
				if foundChar {
					breakLoop = true
					break
				}
				break
			}

			//We done readin bois
			readPos--
			breakLoop = true
		}
	}

	//Convert to number
	number, err := bytesliceToNumber(numberData)
	if err != nil {
		return 0, nil
	}

	//We have us a number!
	//Read from memory or nah?
	if memoryLookup {
		return memoryRead(number), nil
	}

	//Not supposed to look at memory, return number itself
	return number, nil
}

//Get the next program character that ISN'T whitespace
func getNextNonWhitespace() byte {
	for readPos < len(program) {
		if val := getNextChar(); !isWhitespace(val) {
			return val
		}
	}

	return 0
}

//Get the previous program character that ISN'T whitespace
func getPreviousNonWhitespace() byte {
	for readPos >= 0 {
		if val := getPreviousChar(); !isWhitespace(val) {
			return val
		}
	}

	//Nothing was found
	readPos = 0
	return 0
}

//Get the next character
func getNextChar() byte {
	if readPos == len(program) {
		return 0
	}

	char := program[readPos]
	readPos++
	return char
}

//Get the previous character
func getPreviousChar() byte {
	if readPos < 0 {
		readPos = 0
		return 0
	}

	char := program[readPos]
	readPos--
	return char
}

//Is the given character whitespace or not?
func isWhitespace(char byte) bool {

	//Is the character any of the given values?
	switch char {
	case ' ', '\n', '\t', 0x0D:
		return true
	}

	//It wasn't any of those, so not whitespace
	return false
}

//Is the given character newline or not?
func isNewline(char byte) bool {

	//Is the character any of the given values?
	switch char {
	case '\n', 0x0D:
		return true
	}

	//It wasn't any of those, so not newline
	return false
}

//Read from memory
func memoryRead(pos float64) float64 {

	//Does value exist in memory?
	if val, exists := memory[pos]; exists {

		//It does, return that!
		return val
	}

	//It doesn't, return number itself
	return pos
}

//Main program function
func runProgram() error {
	outOfLoop := false
	for !outOfLoop {

		//Left-hand operator
		lefthand, err := getExpression(false)
		if err != nil {
			return err
		}

		//Repeat this for as long as needed (lefthand modification)
		leaveInnerLoop := false
		for !leaveInnerLoop {
			leaveInnerLoop = true

			//What operation do we do?
			nchar := getNextNonWhitespace()
			switch nchar {

			//Assign a value
			case '=':

				//Get righthand expression
				righthand, err := getExpression(true)
				if err != nil {
					return err
				}

				//Value is valid, assign it
				memory[lefthand] = righthand

			//Modify and assign
			case '-', '+', '*', '/':
				next := getNextChar()

				//Decrement?
				if nchar == '-' && next == '-' {

					//Perform decrementation
					memory[lefthand] = memoryRead(lefthand) - 1
					break
				}

				//Increment?
				if nchar == '+' && next == '+' {

					//Perform incrementation
					memory[lefthand] = memoryRead(lefthand) + 1
					break
				}

				//Modification assignment?
				if next == '=' {

					//Get righthand expression
					righthand, err := getExpression(true)
					if err != nil {
						return err
					}

					//Value is valid, modify it
					switch nchar {
					case '+':
						memory[lefthand] = memoryRead(lefthand) + memoryRead(righthand)
					case '-':
						memory[lefthand] = memoryRead(lefthand) - memoryRead(righthand)
					case '*':
						memory[lefthand] = memoryRead(lefthand) * memoryRead(righthand)
					case '/':
						memory[lefthand] = memoryRead(lefthand) / memoryRead(righthand)
					}
					break
				}

				//Lefthand modification?
				if nchar == '+' || nchar == '-' {

					//Get valid expression
					readPos--
					val, err := getExpression(true)
					if err != nil {
						return err
					}

					//Negate offset maybe
					if nchar == '-' {
						val = -val
					}

					//Change lefthand operator
					lefthand += val
					leaveInnerLoop = false
				}

			//Print numeric value
			case '!':
				if consoleOutput {
					fmt.Print(memoryRead(lefthand))
				}
				if writeToFile {
					outputFile.WriteString(fmt.Sprint(memoryRead(lefthand)))
				}

			//Print ascii character
			case '#':
				if consoleOutput {
					fmt.Print(string(byte(memoryRead(lefthand))))
				}
				if writeToFile {
					outputFile.Write([]byte{(byte(memoryRead(lefthand)))})
				}

			//Read input
			case '"':

				//Read value
				val, err := getInput()
				if err != nil {
					return err
				}

				//Assign it to memory
				memory[lefthand] = val

			//Condition
			case '?':

				//What type of comparison?
				comptype := getNextChar()
				switch comptype {
				case '=': //Equals
				case '!': //NOT equals
				case '<': //Less than
					if getNextChar() == '=' {
						comptype = '('
					} else {
						readPos--
					}
				case '>': //Greater than
					if getNextChar() == '=' {
						comptype = ')'
					} else {
						readPos--
					}
				default: //No known comparison type
					var err numskullError_t
					err.msg = "Expected comparison type, got '" + string(comptype) + "'."
					return err
				}

				righthand, err := getExpression(true)
				if err != nil {
					return err
				}

				//Expect an opening bracket
				var openingBracket byte = getNextNonWhitespace()
				var closingBracket byte = 0
				switch openingBracket {
				case '{':
					closingBracket = '}'
				case '[':
					closingBracket = ']'
				default:
					var err numskullError_t
					err.msg = "Expected an opening bracket. Did not find one."
					return err
				}

				//evaluate condition
				expressionTrue := false
				switch comptype {
				case '=':
					expressionTrue = memoryRead(lefthand) == memoryRead(righthand)
				case '!':
					expressionTrue = memoryRead(lefthand) != memoryRead(righthand)
				case '<':
					expressionTrue = memoryRead(lefthand) < memoryRead(righthand)
				case '>':
					expressionTrue = memoryRead(lefthand) > memoryRead(righthand)
				case '(':
					expressionTrue = memoryRead(lefthand) <= memoryRead(righthand)
				case ')':
					expressionTrue = memoryRead(lefthand) >= memoryRead(righthand)
				}

				//Await closing bracket
				if !expressionTrue {

					for depth := 1; depth != 0 && readPos < len(program); {
						char := getNextNonWhitespace()
						if char == openingBracket {
							depth++
						} else if char == closingBracket {
							depth--
						}
					}

					//Is there not a closing bracket?
					if readPos >= len(program) && program[readPos-1] != closingBracket {
						var err numskullError_t
						err.msg = "Expected a closing bracket '" + string(closingBracket) + "'. Never found one."
						return err
					}
				}

			//Unknown
			default:

				//Read unknown action
				action := make([]byte, 512)[:0]
				readPos--
				for {
					char := getNextChar()
					if char == 0 || isWhitespace(char) {
						break
					}
					action = append(action, char)
				}

				var err numskullError_t
				err.msg = "Unknown action, \"" + string(action) + "\"."
				return err
			}
		}

		//Expect newline
		for char := byte(0); true; {
			char = getNextChar()

			//This was a newline
			if isNewline(char) || char == 0 {
				break
			} else if isWhitespace(char) {
				continue
			} else {
				//Read the rest of the line
				readPos--
				remainder := make([]byte, 512)[:0]
				for {
					char := getNextChar()
					if char == 0 || isNewline(char) {
						break
					}
					remainder = append(remainder, char)
				}

				//Return error
				var err numskullError_t
				err.msg = "Expected newline, got \"" + string(remainder) + "\"."
				return err
			}
		}

		//Is there even any more program left?
		if getNextNonWhitespace() == 0 {
			outOfLoop = true
		}
		readPos--
	}

	//Everything worked out
	return nil
}

//Prints program usage
func printUsage() {
	fmt.Println("Usage: numskull [-i file] [-t] [-o file] [-c] <program file>")
	fmt.Println("Useful options:")
	fmt.Println("\t", usage_h)
	fmt.Println("\t", usage_v)
	fmt.Println("\t", usage_i)
	fmt.Println("\t", usage_t)
	fmt.Println("\t", usage_o)
	fmt.Println("\t", usage_c)
}

//Converts a slice of bytes to a float64
func bytesliceToNumber(numData []byte) (float64, error) {
	var err numskullError_t

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

func getInput() (float64, error) {

	if readFromFile {

		//OOB read
		if inputPos >= len(input) {
			return -1, nil
		}

		//Read data at position and return
		data := input[inputPos]
		inputPos++
		return data, nil

	} else {

		//Read input from stdin
		var str string
		_, err := fmt.Scan(&str)
		if err != nil {
			return 0, err
		}

		//Convert this to float64 and return
		return bytesliceToNumber([]byte(str))
	}
}

//Close the program in a more elegant way.
func closeProgram(state int) {
	fmt.Println()
	if !writeToFile {
		fmt.Scan()
	}

	os.Exit(state)
}
