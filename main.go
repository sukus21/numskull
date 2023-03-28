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
	usage_main string = "Usage: %s [arguments] <program-file>"
	usage_v    string = "-v, --version                Print program version number"
	usage_h    string = "-h, --help <argument>        How to use this program"
	usage_i    string = "-i, --input <path>           Read input from file"
	usage_o    string = "-o, --output <path>          Write output to file"
	usage_t    string = "-t, --text                   Read input file as text"
	usage_c    string = "-c, --console                Force output write to stdout"
	usage_d    string = "-d, --depth-warn <int>       Warn at call depth"
	usage_q    string = "-q, --depth-quit <int>       Quit at call depth"
)

// Version numbers
const (
	version_interpreter string = "v0.4"
	version_language    string = "v1.2"
)

// Entrypoint, reads command line arguments
func main() {

	//No arguments provided
	if len(os.Args) == 1 {
		printUsage()
		return
	}

	//Create runner struct
	runner := interpreter.New()
	forceConsole := false

	//Read and process arguments
	somethingDone := false
	argPos := 0
	for {

		//Loop condition
		argPos++
		if argPos >= len(os.Args) {
			break
		}
		arg := os.Args[argPos]

		//Is this not a flag?
		if arg[0] != '-' {
			break
		}

		//What kind of flag?
		switch arg {

		//Help :)
		case "-h", "-H", "--help":
			if argPos != 1 {
				fmt.Printf("Error: '%s' should be passed as the first argument\n", arg)
				return
			}
			argPos++
			for argPos < len(os.Args) {
				if argPos != 2 {
					fmt.Println()
					fmt.Println()
				}
				harg := os.Args[argPos]
				switch harg {

				//Help for the version tag
				case "-v", "-V", "--version":
					fmt.Println(usage_v)
					fmt.Println()
					fmt.Println("That is really all it does.")
					fmt.Println("It can be run using", os.Args[0], harg, "to verify the install path.")

				//Help for the help tag
				case "-h", "-H", "--help":
					fmt.Println(usage_h)
					fmt.Println()
					fmt.Println("You made it here!")
					fmt.Println("If you want help with a specific argument, pass in one (or more) arguments after", arg)
					fmt.Println("If you want a list of all possible arguments, simply don't pass an argument after", arg)
					fmt.Println("The program stops once help has been printed.")
					fmt.Println()
					fmt.Println("Example:", os.Args[0], harg, "-v")
					fmt.Println("Shows the help page for the parameter '-v'.")

				//Help for the input tag
				case "-i", "-I", "--input":
					fmt.Println(usage_i)
					fmt.Println()
					fmt.Println("When reading after the end of file, the result is always -1.")
					fmt.Println("If this argument isn't present, input is given through the console.")
					fmt.Println("Input file will be read as binary by default.")
					fmt.Println("To read file as text, look up help for '-t'.")
					fmt.Println()
					fmt.Println("Example:", os.Args[0], harg, "numbers.bin program.nms")
					fmt.Println("Opens program.nms, and reads from numbers.bin when reading input.")

				//Help for the output tag
				case "-o", "-O", "--output":
					fmt.Println(usage_o)
					fmt.Println()
					fmt.Println("When outputting, you can choose to also write that output to a file.")
					fmt.Println("The file is treated as a byte array, so any value written will be rounded and truncated.")
					fmt.Println("If this argument isn't present, the output of the program is displayed in the console.")
					fmt.Println("If you still want console output AND saving to a file, use the \"-c\" argument.")
					fmt.Println("If the program stops due to an error, the file is still saved.")
					fmt.Println()
					fmt.Println("Example:", os.Args[0], harg, "numbers.bin program.nms")
					fmt.Println("Opens program.nms, and saves output to numbers.bin, once the program stops running.")

				//Deprecated type tag
				case "--type":
					fmt.Println("Warning: '--type' is deprecated since interpreter v0.4. Use '--text' instead.")
					harg = "--text"
					fallthrough

				//Help for the text tag
				case "-t", "-T", "--text":
					fmt.Println(usage_t)
					fmt.Println()
					fmt.Println("Normally when reading input from a file, the input will be read as binary.")
					fmt.Println("That is, one byte per input, each consisting of a number from 0 to 255.")
					fmt.Println("Some may not want this behaviour, so passing in '" + harg + "' will make input read as text.")
					fmt.Println("Entries are read as numbers, seperated by whitespace (tabs, spaces, or newlines).")
					fmt.Println("An incorrectly formatted entry will cause an error when trying to read it.")
					fmt.Println("If the '-i' argument isn't present, this argument does nothing.")

				//Help for the console tag
				case "-c", "-C", "--console":
					fmt.Println(usage_c)
					fmt.Println()
					fmt.Println("When outputting to a file, console output is turned off by default.")
					fmt.Println("Use this argument to reenable it, while also writing the output to a file, using '-o'.")
					fmt.Println("If the '-o' argument isn't present, this argument does nothing.")

				//Help for call depth warning tag
				case "-d", "-D", "--depth-warn":
					fmt.Println(usage_d)
					fmt.Println()
					fmt.Println("Recursion is fun and all, but can get out of hand quickly.")
					fmt.Println("With this argument, you can make the interpreter spit warnings when the given depth is reached.")
					fmt.Println("The default value for this argument is 32.")
					fmt.Println("If you don't want any warnings, set this to -1.")
					fmt.Println("If you want to stop execution entirely when reaching a certain depth, check out the flag '-q'.")

				case "-q", "-Q", "--depth-quit":
					fmt.Println(usage_q)
					fmt.Println()
					fmt.Println("If you want to make sure you don't destroy your stack, set this flag.")
					fmt.Println("This will make the application quit completely when the given call depth is reached.")
					fmt.Println("The default value for this argument is to never quit.")
					fmt.Println("If you just want warnings when reaching high call depth, check out the flag '-d'.")

				default:
					fmt.Printf("Error: unknown argument '%s'\n", harg)
					fmt.Println("To get a list of all possible arguments:")
					fmt.Println(os.Args[0], arg)
					return
				}

				argPos++
			}

			//Show all options
			if len(os.Args) == 2 {
				fmt.Println("All options:")
				fmt.Println("\t", usage_h)
				fmt.Println("\t", usage_v)
				fmt.Println("\t", usage_i)
				fmt.Println("\t", usage_t)
				fmt.Println("\t", usage_o)
				fmt.Println("\t", usage_c)
				fmt.Println("\t", usage_d)
				fmt.Println("\t", usage_q)
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

		//Deprecated type tag
		case "--type":
			fmt.Println("Warning: '--type' is deprecated since interpreter v0.4. Use '--text' instead.")
			arg = "--text"
			fallthrough
		//Text reading mode
		case "-t", "-T", "--text":
			runner.InputText = true

		//Specify output file
		case "-o", "-O", "--output":
			argPos++

			//Output file already specified
			if runner.HasOutputFile() {
				fmt.Println("Error: only one output file can be specified")
				fmt.Println()
				fmt.Println("Usage:", usage_o)
				return
			}

			//No file specified
			if argPos >= len(os.Args)-1 {
				fmt.Println("Error: no output file specified")
				fmt.Println()
				fmt.Println("Usage:", usage_o)
				return
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
				fmt.Println()
				fmt.Println("Usage:", usage_i)
				return
			}

			//Open file
			f, err := os.Open(os.Args[argPos])
			if err != nil {
				fmt.Println("Error: could not open input file")
				fmt.Println(err)
				fmt.Println()
				fmt.Println("Usage:", usage_i)
				return
			}

			runner.SetInputFile(f)
			runner.ReadFromFile = true

		//Call depth warning
		case "-d", "-D", "--depth-warn":
			//No depth specified
			if argPos >= len(os.Args)-1 {
				fmt.Println("Error: no depth was specified")
				fmt.Println()
				fmt.Println("Usage:", usage_d)
				return
			}

			argPos++
			depth, err := strconv.Atoi(os.Args[argPos])
			if err != nil {
				fmt.Println("Error: value for argument '" + os.Args[argPos-1] + "' should be integer")
				fmt.Println(err)
				fmt.Println()
				fmt.Println("Usage:", usage_d)
				return
			}

			//Apply
			runner.WarnDepth = depth

		//Call depth force quit
		case "-q", "-Q", "--depth-quit":
			//No depth specified
			if argPos >= len(os.Args)-1 {
				fmt.Println("Error: no depth was specified")
				fmt.Println()
				fmt.Println("Usage:", usage_q)
				return
			}

			argPos++
			depth, err := strconv.Atoi(os.Args[argPos])
			if err != nil {
				fmt.Println("Error: value for argument'" + os.Args[argPos-1] + "' should be integer")
				fmt.Println(err)
				fmt.Println()
				fmt.Println("Usage:", usage_q)
				return
			}

			//Apply
			runner.QuitDepth = depth

		//Unknown parameter
		default:
			fmt.Println("Error: Unknown argument: " + os.Args[argPos])
			printUsage()
			return
		}
	}

	//No more arguments? Was anything achieved?
	trail := os.Args[argPos:]
	if len(trail) == 0 {
		if !somethingDone {
			fmt.Println("Error: no program file specified")
			fmt.Println()
			fmt.Printf(usage_main+"\n", os.Args[0])
		}
		return
	}

	//More than one program file?
	if len(trail) > 1 {
		fmt.Println("Error: only one program file can be used at once")
		fmt.Println()
		fmt.Printf(usage_main+"\n", os.Args[0])
		return
	}

	//Try to read program file
	file, err := os.ReadFile(os.Args[len(os.Args)-1])
	if err != nil {
		fmt.Println("Error: could not open program file")
		fmt.Println(err)
		fmt.Println()
		fmt.Printf(usage_main+"\n", os.Args[0])
		return
	}

	//Console?
	if forceConsole {
		runner.WriteToStdout = true
		if !runner.HasOutputFile() {
			fmt.Println("Warning: '-c' passed, but no output file specified")
		}
	}
	if runner.InputText && !runner.HasInputFile() {
		fmt.Println("Warning: '-t' passed, but no input file specified")
	}

	//Start executing it
	program, success := parser.ParseProgram(string(file))
	if success {
		err = runner.Execute(program)
	}

	//Close interpreter
	runner.Close()
	if err != nil {
		os.Exit(1)
	}
}

// Prints program usage
func printUsage() {
	fmt.Println()
	fmt.Printf(usage_main+"\n", os.Args[0])
	fmt.Println("Useful options:")
	fmt.Println("\t", usage_h)
	fmt.Println("\t", usage_v)
	fmt.Println("\t", usage_i)
	fmt.Println("\t", usage_o)
	fmt.Println("\t", usage_d)
}
