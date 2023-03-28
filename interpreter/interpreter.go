package interpreter

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sukus21/numskull/token"

	"github.com/sukus21/numskull/utils"
)

func New() *interpreter {
	return &interpreter{
		memory:        make(map[float64]float64),
		WriteToStdout: true,
		QuitDepth:     -1,
		WarnDepth:     32,
	}
}

type interpreter struct {
	WriteToStdout bool
	WriteToFile   bool
	ReadFromFile  bool
	InputText     bool
	WarnDepth     int
	QuitDepth     int
	outputFile    io.WriteCloser
	outputWriter  io.Writer
	inputReader   *bufio.Reader
	inputFile     io.ReadCloser
	memory        map[float64]float64
}

func (inter *interpreter) SetInputFile(file io.ReadCloser) {
	inter.inputFile = file
	inter.inputReader = bufio.NewReader(file)
}

func (inter *interpreter) HasInputFile() bool {
	return inter.inputFile != nil
}

func (inter *interpreter) SetOutputFile(file io.WriteCloser) {
	inter.outputFile = file
	inter.outputWriter = file
}

func (inter *interpreter) HasOutputFile() bool {
	return inter.outputFile != nil
}

// Read from memory
func (inter *interpreter) memoryRead(pos float64) float64 {

	//Does value exist in memory?
	if val, exists := inter.memory[pos]; exists {

		//It does, return that!
		return val
	}

	//It doesn't, return number itself
	return pos
}

// Get input from file or command line
func (inter *interpreter) getInput() (float64, error) {
	if inter.ReadFromFile {
		if !inter.InputText {
			b, err := inter.inputReader.ReadByte()
			if err != nil {
				if err == io.EOF {
					return -1, nil
				}
			}
			return float64(b), err
		} else {
			dat, err := inter.inputReader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					return -1, nil
				}
				return 0, err
			}
			return utils.BytesliceToNumber(dat)
		}
	} else {

		//Read input from stdin
		num := float64(0)
		_, err := fmt.Scan(&num)
		if err != nil {
			return 0, err
		}
		return num, err
	}
}

// Write to output file and/or stdout
func (inter *interpreter) writeOutput(data []byte) {
	if inter.WriteToStdout {
		fmt.Print(string(data))
	}
	if inter.WriteToFile {
		inter.outputWriter.Write(data)
	}
}

// Main program function
func (inter *interpreter) Execute(program []float64) error {
	inter.memory = make(map[float64]float64)
	callstack := make([]int, 0, 64)
	for readPos := 0; readPos < len(program); {

		tok := token.Token(program[readPos])
		readPos++
		switch tok {

		//End of function (return)
		case token.FunctionEnd:
			if len(callstack) == 0 {
				return printError("FATAL: cannot return from function, call stack empty")
			}

			readPos = callstack[len(callstack)-1]
			callstack = callstack[:len(callstack)-1]

		//Jump indicator
		case token.FunctionStart, token.SquareEnd:

			//Read jump point and jump
			readPos = int(program[readPos])
			continue

		//It's a number
		case token.Number:
			lefthand := program[readPos]
			readPos++

			//Start chaining lefthands
			for {

				//Make sure this is a chainer
				tok = token.Token(program[readPos])
				readPos++
				if tok != token.ChainMinus && tok != token.ChainPlus {
					break
				}
				readPos++

				//Chain lefthand
				if tok == token.ChainMinus {
					lefthand -= inter.memoryRead(program[readPos])
				} else {
					lefthand += inter.memoryRead(program[readPos])
				}
				readPos++
			}

			//We got us an operation
			switch tok {
			case token.Increment:
				inter.memory[lefthand] = inter.memoryRead(lefthand) + 1
			case token.Decrement:
				inter.memory[lefthand] = inter.memoryRead(lefthand) - 1

			case token.Assign:
				readPos++
				righthand := program[readPos]
				readPos++
				inter.memory[lefthand] = inter.memoryRead(righthand)
			case token.Add:
				readPos++
				righthand := program[readPos]
				readPos++
				inter.memory[lefthand] = inter.memoryRead(lefthand) + inter.memoryRead(righthand)
			case token.Sub:
				readPos++
				righthand := program[readPos]
				readPos++
				inter.memory[lefthand] = inter.memoryRead(lefthand) - inter.memoryRead(righthand)
			case token.Multiply:
				readPos++
				righthand := program[readPos]
				readPos++
				inter.memory[lefthand] = inter.memoryRead(lefthand) * inter.memoryRead(righthand)
			case token.Divide:
				readPos++
				righthand := program[readPos]
				readPos++
				inter.memory[lefthand] = inter.memoryRead(lefthand) / inter.memoryRead(righthand)

			case token.PrintChar:
				inter.writeOutput([]byte{byte(inter.memoryRead(lefthand))})
			case token.PrintNumber:
				inter.writeOutput([]byte(fmt.Sprint(inter.memoryRead(lefthand))))
			case token.ReadInput:
				//Read value
				val, err := inter.getInput()
				if err != nil {
					return printError("FATAL: %s", err)
				}

				//Assign it to memory
				inter.memory[lefthand] = val

			case token.Equals:
				readPos++
				righthand := inter.memoryRead(program[readPos])
				readPos++
				if inter.memoryRead(lefthand) == righthand {
					readPos++
				} else {
					//Jump
					readPos = int(program[readPos])
				}
			case token.Different:
				readPos++
				righthand := inter.memoryRead(program[readPos])
				readPos++
				if inter.memoryRead(lefthand) != righthand {
					readPos++
				} else {
					//Jump
					readPos = int(program[readPos])
				}
			case token.LessThan:
				readPos++
				righthand := inter.memoryRead(program[readPos])
				readPos++
				if inter.memoryRead(lefthand) < righthand {
					readPos++
				} else {
					//Jump
					readPos = int(program[readPos])
				}
			case token.LessEquals:
				readPos++
				righthand := inter.memoryRead(program[readPos])
				readPos++
				if inter.memoryRead(lefthand) <= righthand {
					readPos++
				} else {
					//Jump
					readPos = int(program[readPos])
				}
			case token.GreaterThan:
				readPos++
				righthand := inter.memoryRead(program[readPos])
				readPos++
				if inter.memoryRead(lefthand) > righthand {
					readPos++
				} else {
					//Jump
					readPos = int(program[readPos])
				}
			case token.GreaterEquals:
				readPos++
				righthand := inter.memoryRead(program[readPos])
				readPos++
				if inter.memoryRead(lefthand) >= righthand {
					readPos++
				} else {
					//Jump
					readPos = int(program[readPos])
				}

			case token.FunctionRun:

				//Push current position onto program stack
				callstack = append(callstack, readPos)
				if len(callstack) == inter.WarnDepth {
					printError(fmt.Sprintf("WARNING: call depth of %d reached", inter.WarnDepth))
				}
				if len(callstack) == inter.QuitDepth {
					return printError(fmt.Sprintf("FATAL: max call depth of %d reached", inter.QuitDepth))
				}

				//Move read position and verify function
				readPos = int(inter.memoryRead(lefthand))
				if program[readPos] != float64(token.FunctionStart) {
					return printError("FATAL: invalid function call")
				}
				readPos += 2

			default:
				return printError("FATAL: unknown operation '%s'", tok.GetTokenName())
			}
		}
	}

	//Everything worked out
	return nil
}

func (inter *interpreter) Close() {
	if inter.inputFile != nil {
		inter.inputFile.Close()
	}
	if inter.outputFile != nil {
		inter.outputFile.Close()
	}
}

// Pretty formatted error messages to stdout
func printError(raw string, args ...any) error {
	msg := fmt.Sprintf(raw, args...)
	fmt.Fprintf(os.Stderr, "%s\n  %s  \n%s\n",
		strings.Repeat("_", len(msg)+4),
		msg,
		strings.Repeat("‾", len(msg)+4))
	return fmt.Errorf(msg)
}
