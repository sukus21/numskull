package token

//Token struct
type Token float64

const (
	Invalid Token = iota
	Number
	ChainPlus
	ChainMinus
	Assign
	Add
	Sub
	Multiply
	Divide
	Increment
	Decrement
	PrintNumber
	PrintChar
	ReadInput
	Equals
	Different
	LessThan
	GreaterThan
	LessEquals
	GreaterEquals
	CurlyStart
	CurlyEnd
	SquareStart
	SquareEnd
	Newline
	FunctionStart
	FunctionEnd
	FunctionRun
)

//Returns the name of the given token as a string
func (tok Token) GetTokenName() string {
	switch tok {
	case Add:
		return "+="
	case Sub:
		return "-="
	case Assign:
		return "="
	case Multiply:
		return "*="
	case Divide:
		return "/="
	case Increment:
		return "++"
	case Decrement:
		return "--"
	case ChainPlus:
		return "+"
	case ChainMinus:
		return "-"

	case Equals:
		return "?="
	case Different:
		return "?!"
	case LessEquals:
		return "?<="
	case LessThan:
		return "?<"
	case GreaterEquals:
		return "?>="
	case GreaterThan:
		return "?>"

	case Number:
		return "number"
	case PrintChar:
		return "#"
	case PrintNumber:
		return "!"
	case ReadInput:
		return "\""

	case CurlyStart:
		return "{"
	case CurlyEnd:
		return "}"
	case SquareStart:
		return "["
	case SquareEnd:
		return "]"

	case FunctionStart:
		return "<"
	case FunctionEnd:
		return ">"
	case FunctionRun:
		return "()"

	case Newline:
		return "newline"
	case Invalid:
		return "invalid"

	default:
		return "unknown"
	}
}
