package main

import (
	"fmt"
	"os"
)

const (
	usage_v string = "-v, --version         Prints program version number"
	usage_h string = "-h, --help <argument> Prints usage for given argument"
	usage_i string = "-i, --input <path>    File to read input from"
	usage_o string = "-o, --output <path>   File to print output to"
	usage_t string = "-t, --type            Tells program to read input file as text"
	usage_c string = "-c, --console         Force program output to console"
)

var program []float64
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
	//os.Args = []string{os.Args[0], "example programs/echo.nms"}

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
				fmt.Println("Example: numskull", os.Args[argPos-1], "v")
				fmt.Println("Shows the help page for the parameter '-v'.")
				fmt.Println("The program stops once help has been printed.")

			//Help for the input tag
			case "i", "I", "input":
				if len(os.Args[argPos]) != 1 {
					os.Args[argPos] = "-" + os.Args[argPos]
				}
				fmt.Println(usage_i)
				fmt.Println()
				fmt.Println("When reading after the end of file, the result is always -1.")
				fmt.Println("If this argument isn't present, input is given through the console.")
				fmt.Println("Input file will be read as binary by default, look up \"numskull --help -t\" for more info.")
				fmt.Println()
				fmt.Println("Example: numskull", "-"+os.Args[argPos], "numbers.bin program.nms")
				fmt.Println("Opens program.nms, and reads from numbers.bin when reading input.")

			//Help for the output tag
			case "o", "O", "output":
				if len(os.Args[argPos]) != 1 {
					os.Args[argPos] = "-" + os.Args[argPos]
				}
				fmt.Println(usage_o)
				fmt.Println()
				fmt.Println("When outputting, you can choose to also write that output to a file.")
				fmt.Println("The file is treated as a byte array.")
				fmt.Println("If this argument isn't present, the output of the program is displayed in the console.")
				fmt.Println("If you still want console output AND saving to a file, use the \"-c\" argument.")
				fmt.Println("If the program stops due to an error, the file is still saved.")
				fmt.Println()
				fmt.Println("Example: numskull", "-"+os.Args[argPos], "numbers.bin program.nms")
				fmt.Println("Opens program.nms, and saves output to numbers.bin, once the program stops running.")

			//Help for the type tag
			case "t", "T", "type":
				if len(os.Args[argPos]) != 1 {
					os.Args[argPos] = "-" + os.Args[argPos]
				}
				fmt.Println(usage_t)
				fmt.Println()
				fmt.Println("Normally when reading input from a file, the input will be read as binary.")
				fmt.Println("That is, one byte per input, each consisting of a number from 0 to 255.")
				fmt.Println("Some may not want this behaviour, so passing in", "-"+os.Args[argPos], "will make input read as text.")
				fmt.Println("Entries are read as numbers, seperated by whitespace (tabs, spaces, or newlines).")
				fmt.Println("An incorrectly formatted entry will cause an error when trying to read it.")
				fmt.Println("If the -i argument isn't present, this argument does nothing.")

			//Help for the console tag
			case "c", "C", "console":
				fmt.Println(usage_c)
				fmt.Println()
				fmt.Println("When outputting to a file, console output is turned off by default.")
				fmt.Println("Use this argument to reenable it, while also writing the output to a file, using -o.")
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
			fmt.Println("runs Numskull version 1.0")
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
				number, err := BytesliceToNumber(numData)
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

	success := true
	program, success = ParseProgram(string(file))
	if success {
		err = runProgram(program)
	}

	//Handle program output
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
func runProgram(program []float64) error {

	for readPos := 0; readPos < len(program); {

		tok := token(program[readPos])
		readPos++

		//Looping bracket?
		if tok == token_squareEnd {

			//Read jump point and jump
			readPos = int(program[readPos])
			continue
		}

		//It's a number
		lefthand := program[readPos]
		readPos++

		//Start chaining lefthands
		for {

			//Make sure this is a chainer
			tok = token(program[readPos])
			readPos++
			if tok != token_chainMinus && tok != token_chainPlus {
				break
			}
			readPos++

			//Chain lefthand
			if tok == token_chainMinus {
				lefthand -= memoryRead(program[readPos])
			} else {
				lefthand += memoryRead(program[readPos])
			}
			readPos++
		}

		//We got us an operation
		switch tok {
		case token_increment:
			memory[lefthand] = memoryRead(lefthand) + 1
		case token_decrement:
			memory[lefthand] = memoryRead(lefthand) - 1

		case token_assign:
			readPos++
			righthand := program[readPos]
			readPos++
			memory[lefthand] = memoryRead(righthand)
		case token_add:
			readPos++
			righthand := program[readPos]
			readPos++
			memory[lefthand] = memoryRead(lefthand) + memoryRead(righthand)
		case token_sub:
			readPos++
			righthand := program[readPos]
			readPos++
			memory[lefthand] = memoryRead(lefthand) - memoryRead(righthand)
		case token_multiply:
			readPos++
			righthand := program[readPos]
			readPos++
			memory[lefthand] = memoryRead(lefthand) * memoryRead(righthand)
		case token_divide:
			readPos++
			righthand := program[readPos]
			readPos++
			memory[lefthand] = memoryRead(lefthand) / memoryRead(righthand)

		case token_printChar:
			if consoleOutput {
				fmt.Print(string(byte(memoryRead(lefthand))))
			}
			if writeToFile {
				outputFile.Write([]byte{(byte(memoryRead(lefthand)))})
			}
		case token_printNumber:
			if consoleOutput {
				fmt.Print(memoryRead(lefthand))
			}
			if writeToFile {
				outputFile.WriteString(fmt.Sprint(memoryRead(lefthand)))
			}
		case token_readInput:
			//Read value
			val, err := getInput()
			if err != nil {
				return err
			}

			//Assign it to memory
			memory[lefthand] = val

		case token_equals:
			readPos++
			righthand := memoryRead(program[readPos])
			readPos++
			if memoryRead(lefthand) == righthand {
				readPos++
			} else {
				//Jump
				readPos = int(program[readPos])
			}
		case token_different:
			readPos++
			righthand := memoryRead(program[readPos])
			readPos++
			if memoryRead(lefthand) != righthand {
				readPos++
			} else {
				//Jump
				readPos = int(program[readPos])
			}
		case token_lessThan:
			readPos++
			righthand := memoryRead(program[readPos])
			readPos++
			if memoryRead(lefthand) < righthand {
				readPos++
			} else {
				//Jump
				readPos = int(program[readPos])
			}
		case token_lessEquals:
			readPos++
			righthand := memoryRead(program[readPos])
			readPos++
			if memoryRead(lefthand) <= righthand {
				readPos++
			} else {
				//Jump
				readPos = int(program[readPos])
			}
		case token_greaterThan:
			readPos++
			righthand := memoryRead(program[readPos])
			readPos++
			if memoryRead(lefthand) > righthand {
				readPos++
			} else {
				//Jump
				readPos = int(program[readPos])
			}
		case token_greaterEquals:
			readPos++
			righthand := memoryRead(program[readPos])
			readPos++
			if memoryRead(lefthand) >= righthand {
				readPos++
			} else {
				//Jump
				readPos = int(program[readPos])
			}

		default:
			fmt.Print("Uh oh, something went fucky :(")
		}
	}

	//Everything worked out
	return nil
}

//Prints program usage
func printUsage() {
	fmt.Println("Usage: numskull [-i file] [-t] [-o file] [-c] <program-file>")
	fmt.Println("Useful options:")
	fmt.Println("\t", usage_h)
	fmt.Println("\t", usage_v)
	fmt.Println("\t", usage_i)
	fmt.Println("\t", usage_t)
	fmt.Println("\t", usage_o)
	fmt.Println("\t", usage_c)
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
		return BytesliceToNumber([]byte(str))
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
