package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func AssemblyFileToRom(iFile *os.File, oFile *os.File) error {
	labelLookupTable = buildLabelLookupTable(iFile)

	program := AssembleFromFile(iFile)

	writer := bufio.NewWriter(oFile)
	for i := 0; (i + 1) < len(program); i += 2 {
		writer.WriteByte(program[i])
		writer.WriteByte(program[i+1])
	}
	writer.Flush()

	return nil
}

var labelLookupTable *map[string]uint16 = nil

func buildLabelLookupTable(file *os.File) *map[string]uint16 {
	table := make(map[string]uint16)

	reader := bufio.NewReader(file)

	var address uint16 = 512
	for line, _, _ := reader.ReadLine(); line != nil; {
		tokens := TokenizeLine(string(line))
		if tokens.label != "" {
			table[tokens.label] = address
		}
		if tokens.instruction != "" {
			address += 2
		}
		line, _, _ = reader.ReadLine()
	}

	file.Seek(0, 0)

	return &table
}

func AssembleFromFile(file *os.File) []uint8 {
	var program []uint8

	reader := bufio.NewReader(file)

	for line, _, _ := reader.ReadLine(); line != nil; {
		tokens := TokenizeLine(string(line))

		if tokens.instruction != "" {
			high, low := AssembleFromLine(tokens)
			program = append(program, high, low)
		}

		line, _, _ = reader.ReadLine()
	}

	return program
}

func AssembleFromLine(tokens chip8AsmLineTokens) (high uint8, low uint8) {
	switch tokens.instruction {
	case NOP:
		return 0x00, 0x00
	case CLS:
		return 0x00, 0xe0
	case RET:
		return 0x00, 0xee
	case JMP:
		var opcat uint8
		var addr uint16
		if tokens.operand2 == "" {
			opcat = 0x10
			addr = getNNN(tokens.operand1)
		} else {
			opcat = 0xB0
			addr = getNNN(tokens.operand2)
		}
		high = opcat | uint8((addr>>8)&0x0f)
		low = uint8(addr & 0xff)
		return high, low
	case CALL:
		opcat := uint8(0x20)
		addr := getNNN(tokens.operand1)
		high = opcat | uint8((addr>>8)&0x0f)
		low = uint8(addr & 0xff)
		return high, low
	case SE:
		x := getV(tokens.operand1)
		if tokens.operand2[0:1] != V {
			opcat := uint8(0x30)
			nn := getNN(tokens.operand2)
			high = opcat | x
			low = nn
		} else {
			opcat := uint8(0x50)
			y := getV(tokens.operand2)
			high = opcat | x
			low = y << 4
		}
		return high, low
	case SNE:
		x := getV(tokens.operand1)
		if tokens.operand2[0:1] != V {
			opcat := uint8(0x40)
			nn := getNN(tokens.operand2)
			high = opcat | x
			low = nn
		} else {
			opcat := uint8(0x90)
			y := getV(tokens.operand2)
			high = opcat | x
			low = y << 4
		}
		return high, low
	case LD:
		if tokens.operand1[0:1] == V {
			var opcat uint8
			x := getV(tokens.operand1)
			if tokens.operand2[0:1] == V {
				opcat = 0x80
				y := getV(tokens.operand2)
				low = y << 4
			} else {
				switch tokens.operand2 {
				case DT:
					opcat = 0xf0
					low = 0x07
				case K:
					opcat = 0xf0
					low = 0x0a
				case R:
					opcat = 0xf0
					low = 0x65
				default:
					opcat = 0x60
					nn := getNN(tokens.operand2)
					low = nn
				}
			}
			high = opcat | x
		} else {
			switch tokens.operand1 {
			case I:
				opcat := uint8(0xa0)
				nnn := getNNN(tokens.operand2)
				high = opcat | uint8(nnn>>8)
				low = uint8(nnn & 0xff)
			case DT:
				opcat := uint8(0xf0)
				x := getV(tokens.operand2)
				high = opcat | x
				low = 0x15
			case ST:
				opcat := uint8(0xf0)
				x := getV(tokens.operand2)
				high = opcat | x
				low = 0x18
			case F:
				opcat := uint8(0xf0)
				x := getV(tokens.operand2)
				high = opcat | x
				low = 0x29
			case B:
				opcat := uint8(0xf0)
				x := getV(tokens.operand2)
				high = opcat | x
				low = 0x33
			case R:
				opcat := uint8(0xf0)
				x := getV(tokens.operand2)
				high = opcat | x
				low = 0x55
			default:
				high, low = 0, 0
			}
		}
		return high, low
	case ADD:
		var opcat, x uint8
		if tokens.operand1[:1] == V {
			x = getV(tokens.operand1)
			if tokens.operand2[:1] == V {
				opcat = uint8(0x80)
				y := getV(tokens.operand2)
				low = (y << 4) | 0x04
			} else {
				opcat = uint8(0x70)
				low = getNN(tokens.operand2)
			}
		} else {
			opcat = uint8(0xf0)
			x = getV(tokens.operand2)
			low = 0x1e
		}
		high = opcat | x
		return high, low
	case OR:
		opcat := uint8(0x80)
		x := getV(tokens.operand1)
		y := getV(tokens.operand2)
		high = opcat | x
		low = (y << 4) | 0x01
		return high, low
	case AND:
		opcat := uint8(0x80)
		x := getV(tokens.operand1)
		y := getV(tokens.operand2)
		high = opcat | x
		low = (y << 4) | 0x02
		return high, low
	case XOR:
		opcat := uint8(0x80)
		x := getV(tokens.operand1)
		y := getV(tokens.operand2)
		high = opcat | x
		low = (y << 4) | 0x03
		return high, low
	case SUB:
		opcat := uint8(0x80)
		x := getV(tokens.operand1)
		y := getV(tokens.operand2)
		high = opcat | x
		low = (y << 4) | 0x05
		return high, low
	case SHR:
		opcat := uint8(0x80)
		x := getV(tokens.operand1)
		high = opcat | x
		low = 0x06
		return high, low
	case SUBN:
		opcat := uint8(0x80)
		x := getV(tokens.operand1)
		y := getV(tokens.operand2)
		high = opcat | x
		low = (y << 4) | 0x07
		return high, low
	case SHL:
		opcat := uint8(0x80)
		x := getV(tokens.operand1)
		high = opcat | x
		low = 0x0e
		return high, low
	case RND:
		opcat := uint8(0xc0)
		x := getV(tokens.operand1)
		high = opcat | x
		low = getNN(tokens.operand2)
		return high, low
	case SKP:
		opcat := uint8(0xe0)
		x := getV(tokens.operand1)
		high = opcat | x
		low = 0x9e
		return high, low
	case SKNP:
		opcat := uint8(0xe0)
		x := getV(tokens.operand1)
		high = opcat | x
		low = 0xa1
		return high, low
	case DRAW:
		opcat := uint8(0xD0)
		x := getV(tokens.operand1)
		y := getV(tokens.operand2)
		n := getN(tokens.operand3)
		high = opcat | x
		low = (y << 4) | n
		return high, low
	default:
		return 0, 0
	}
}

func getV(s string) uint8 {
	v, _ := strconv.ParseUint(s[1:], 16, 4)
	return uint8(v)
}

func getN(s string) uint8 {
	n, _ := strconv.ParseUint(s, 10, 4)
	return uint8(n)
}

func getNN(s string) uint8 {
	nn, _ := strconv.ParseUint(s, 10, 8)
	return uint8(nn)
}

func getNNN(s string) uint16 {
	var nnn uint64
	if s[:1] == DOLLAR {
		nnn = uint64((*labelLookupTable)[s[1:]])
	} else {
		nnn, _ = strconv.ParseUint(s, 10, 12)
	}
	return uint16(nnn)
}

type chip8AsmLineTokens struct {
	label       string
	instruction string
	operand1    string
	operand2    string
	operand3    string
	comment     string
}

func TokenizeLine(line string) (tokens chip8AsmLineTokens) {
	line = strings.ToUpper(line)
	tokens.comment, line = getComment(line)
	tokens.label, line = getLabel(line)
	tokens.instruction, line = getInstruction(line)
	tokens.operand1, line = getOperand(line)
	tokens.operand2, line = getOperand(line)
	tokens.operand3, _ = getOperand(line)
	return tokens
}

func getLabel(line string) (label string, slice string) {
	i := strings.Index(line, LABEL)
	if i == -1 {
		return "", line
	}

	label = strings.TrimSpace(line[:i])
	slice = strings.TrimSpace(line[i+1:])
	return label, slice
}

func getInstruction(line string) (instruction string, slice string) {
	line = strings.TrimSpace(line)

	i := strings.Index(line, SPACE)
	if i == -1 {
		return line, ""
	}

	instruction = strings.TrimSpace(line[:i])
	slice = strings.TrimSpace(line[i:])
	return instruction, slice
}

func getOperand(line string) (operand string, slice string) {
	i := strings.Index(line, COMMA)
	if i == -1 {
		operand = strings.TrimSpace(line)
		return operand, ""
	}

	slice = strings.TrimSpace(line[i+1:])
	operand = strings.TrimSpace(line[:i])
	return operand, slice
}

func getComment(line string) (comment string, slice string) {
	i := strings.Index(line, COMMENT)
	if i == -1 {
		return "", line
	}

	comment = strings.TrimSpace(line[i+1:])
	slice = strings.TrimSpace(line[:i])
	return comment, slice
}
