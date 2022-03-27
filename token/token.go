package token

//Token struct
type Token float64

const (
	Invalid       Token = iota
	Number        Token = iota
	ChainPlus     Token = iota
	ChainMinus    Token = iota
	Assign        Token = iota
	Add           Token = iota
	Sub           Token = iota
	Multiply      Token = iota
	Divide        Token = iota
	Increment     Token = iota
	Decrement     Token = iota
	PrintNumber   Token = iota
	PrintChar     Token = iota
	ReadInput     Token = iota
	Equals        Token = iota
	Different     Token = iota
	LessThan      Token = iota
	GreaterThan   Token = iota
	LessEquals    Token = iota
	GreaterEquals Token = iota
	CurlyStart    Token = iota
	CurlyEnd      Token = iota
	SquareStart   Token = iota
	SquareEnd     Token = iota
	Newline       Token = iota
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
	case Newline:
		return "newline"
	case Invalid:
		return "invalid"

	default:
		return "unknown"
	}
}
