package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/sukus21/numskull/interpreter"
	"github.com/sukus21/numskull/parser"
)

// Usage strings
const (
	usage_v string = "-v, --version                Prints program version number"
	usage_h string = "-h, --help <argument>        Prints usage for given argument"
	usage_i string = "-i, --input <path>           File to read input from"
	usage_o string = "-o, --output <path>          File to print output to"
	usage_t string = "-t, --type                   Tells program to read input file as text"
	usage_c string = "-c, --console                Force program output to console"
	usage_d string = "-d, --depth-warn <number>    Warn at call depth"
	usage_q string = "-q, --depth-quit <number>    Quit at call depth"
)

// Version numbers
const (
	version_interpreter string = "v0.4"
	version_language    string = "v1.2"
)

// Entrypoint, reads command line arguments
func main() {

	//No arguments provided
	fmt.Println()
	if len(os.Args) == 1 {
		printUsage()
		return
	}

	//Create runner struct
	runner := interpreter.New()

	//Need to keep track of this separately
	forceConsole := false

	//Read and process arguments
	var somethingDone bool = false
	argPos := 1
	for ; argPos < len(os.Args); argPos++ {
		switch os.Args[argPos] {

		//Help :)
		case "-h", "-H", "--help":
			argPos++

			//No help thing provided
			if argPos == len(os.Args) {
				fmt.Println("Error: specify an argument to explain.")
				fmt.Println("Example:", os.Args[0], os.Args[argPos-1], "v")
				fmt.Println("Pass in no arguments to see list of available arguments.")
				return
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
				fmt.Println("If the '-i' argument isn't present, this argument does nothing.")

			//Help for the console tag
			case "c", "C", "console":
				fmt.Println(usage_c)
				fmt.Println()
				fmt.Println("When outputting to a file, console output is turned off by default.")
				fmt.Println("Use this argument to reenable it, while also writing the output to a file, using -o.")
				fmt.Println("If the '-o' argument isn't present, this argument does nothing.")

			//Help for call depth warning tag
			case "d", "D", "depth-warn":
				fmt.Println(usage_d)
				fmt.Println()
				fmt.Println("Recursion is fun and all, but can get out of hand quickly.")
				fmt.Println("With this argument, you can make the interpreter spit warnings when the given depth is reached.")
				fmt.Println("The default value for this argument is 32.")
				fmt.Println("If you don't want any warnings, set this to -1.")
				fmt.Println("If you want to stop execution entirely when reaching a certain depth, check out the flag '-q'.")

			case "q", "Q", "depth-quit":
				fmt.Println(usage_q)
				fmt.Println()
				fmt.Println("If you want to make sure you don't destroy your stack, set this flag.")
				fmt.Println("This will make the application quit completely when the given call depth is reached.")
				fmt.Println("The default value for this argument is to never quit.")
				fmt.Println("If you just want warnings when reaching high call depth, check out the flag '-d'.")

			default:
				fmt.Println("Error: unknown argument.")
				fmt.Println()
				printUsage()
				return
			}

			//Exit
			return

		//Print numskull version
		case "-v", "-V", "--version":
			fmt.Println("Numskull interpreter", version_interpreter)
			fmt.Println("runs Numskull version", version_language)
			somethingDone = true

		//Force console output
		case "-c", "-C", "--console":
			forceConsole = true

		//Text reading mode
		case "-t", "-T", "--type":
			runner.InputText = true

		//Specify output file
		case "-o", "-O", "--output":
			argPos++

			//Output file already specified
			if runner.HasOutputFile() {
				fmt.Println("Error: only one output file can be specified")
				return
			}

			//No file specified
			if argPos >= len(os.Args)-1 {
				fmt.Println("Error: no output file specified")
				fmt.Println(usage_o)
				break
			}

			//Update stuff
			f, err := os.Create(os.Args[argPos])
			if err != nil {
				fmt.Println("Error: could not create output file")
				fmt.Println(err)
				return
			}
			runner.SetOutputFile(f)
			runner.WriteToStdout = false
			runner.WriteToFile = true

		//Read input file
		case "-i", "-I", "--input":
			argPos++

			//Input file already specified
			if runner.HasInputFile() {
				fmt.Println("Error: only one input file can be specified")
				return
			}

			//No file specified
			if argPos >= len(os.Args)-1 {
				fmt.Println("Error: no input file specified")
				fmt.Println(usage_i)
				break
			}

			//Open file
			f, err := os.Open(os.Args[argPos])
			if err != nil {
				fmt.Println("Error: could not open input file")
				fmt.Println(err)
				fmt.Println(usage_i)
				return
			}

			runner.SetInputFile(f)
			runner.ReadFromFile = true

		case "-d", "-D", "--depth-warn":
			//No depth specified
			if argPos >= len(os.Args)-1 {
				fmt.Println("Error: no depth was specified")
				fmt.Println(usage_d)
				break
			}

			argPos++
			depth, err := strconv.Atoi(os.Args[argPos])
			if err != nil {
				fmt.Println("Error: value for argument'" + os.Args[argPos-1] + "' should be integer")
				fmt.Println(err)
				fmt.Println(usage_d)
				break
			}

			//Apply
			runner.WarnDepth = depth

		case "-q", "-Q", "--depth-quit":
			//No depth specified
			if argPos >= len(os.Args)-1 {
				fmt.Println("Error: no depth was specified")
				fmt.Println(usage_q)
				break
			}

			argPos++
			depth, err := strconv.Atoi(os.Args[argPos])
			if err != nil {
				fmt.Println("Error: value for argument'" + os.Args[argPos-1] + "' should be integer")
				fmt.Println(err)
				fmt.Println(usage_q)
				break
			}

			//Apply
			runner.QuitDepth = depth

		//Unknown parameter
		default:
			if argPos == len(os.Args)-1 {
				break
			}

			fmt.Println("Error: Unknown argument: " + os.Args[argPos])
			printUsage()
			return
		}
	}

	//No more arguments? Was anything achieved?
	finArg := os.Args[len(os.Args)-1]
	if argPos >= len(os.Args)-1 {
		if !somethingDone {
			fmt.Println("Error: no program file specified")
		}
		return
	}

	//Try to read program file
	file, err := os.ReadFile(finArg)
	if err != nil {
		fmt.Println("Error: could not open program file")
		fmt.Println(err)
		return
	}

	//Console?
	if forceConsole {
		runner.WriteToStdout = true
	}

	//Start executing it
	program, success := parser.ParseProgram(string(file))
	if success {

		//Run program
		err = runner.Execute(program)

		//Handle program output
		fmt.Println()
		fmt.Println()
		if err != nil {
			//Print error
			fmt.Println(err)
		} else {
			//Program finished :)
			fmt.Println("Program finished :)")
		}
	}

	//Close interpreter
	runner.Close()
}

// Prints program usage
func printUsage() {
	fmt.Println("Usage: numskull [-i file] [-t] [-o file] [-c] <program file>")
	fmt.Println("Useful options:")
	fmt.Println("\t", usage_h)
	fmt.Println("\t", usage_v)
	fmt.Println("\t", usage_i)
	fmt.Println("\t", usage_t)
	fmt.Println("\t", usage_o)
	fmt.Println("\t", usage_c)
	fmt.Println("\t", usage_d)
	fmt.Println("\t", usage_q)
}
