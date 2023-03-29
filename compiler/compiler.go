package compiler

type memSize int

const (
	memSize8  memSize = 8
	memSize16 memSize = 16
	memSize32 memSize = 32
	memSize64 memSize = 64
)

type Symbol struct {
	Resolved bool
	Address  uint64
}

type Architecture interface {
	EncodeReadMemory(address *Symbol, size memSize)
	EncodeWriteMemory(address *Symbol, size memSize, value uint64)
}

type architectureARM64 int

func (architectureARM64) ReadMemory(address *Symbol, size memSize)

var ArchitectureARM64 architectureARM64

type Compiler struct {
	Architecture
}
