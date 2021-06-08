package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

func RomToAssemblyFile(inFile *os.File, outFile *os.File) error {
	reader := (*Chip8RomReader)(bufio.NewReader(inFile))
	writer := bufio.NewWriter(outFile)

	// Start counting line numbers at start of CHIP8 program memory
	for i := 0x200; ; i += 2 {
		o, err := reader.GetNextOpcode()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Format disassembly output to include line numbers and opcodes neatly aligned
		lineNumber := fmt.Sprintf("%4v", i) + LABEL + " "
		instruction := fmt.Sprintf("%-24v", GetChip8Disassembly(o))
		comment := COMMENT + " 0x" + fmt.Sprintf("%04v", strconv.FormatUint(uint64(o.opcode), 16))
		writer.WriteString(lineNumber + instruction + comment + "\n")
	}

	writer.Flush()
	return nil
}

func GetChip8Disassembly(o Chip8Opcode) string {
	switch o.c {
	case 0x0000:
		switch o.nn {
		case 0xe0:
			return CLS
		case 0xee:
			return RET
		default:
			return NOP
		}
	case 0x1000:
		return fmt.Sprintf(JMP+SPACE+"%v", o.nnn)
	case 0x2000:
		return fmt.Sprintf(CALL+SPACE+"%v", o.nnn)
	case 0x3000:
		return fmt.Sprintf(SE+SPACE+V+"%v"+COMMA+SPACE+"%v", strconv.FormatUint(uint64(o.x), 16), o.nn)
	case 0x4000:
		return fmt.Sprintf(SNE+SPACE+V+"%v"+COMMA+SPACE+"%v", strconv.FormatUint(uint64(o.x), 16), o.nn)
	case 0x5000:
		switch o.n {
		case 0x0:
			return fmt.Sprintf(SE+SPACE+V+"%v"+COMMA+SPACE+V+"%v", strconv.FormatUint(uint64(o.x), 16), strconv.FormatUint(uint64(o.y), 16))
		default:
			return NOP
		}
	case 0x6000:
		return fmt.Sprintf(LD+SPACE+V+"%v"+COMMA+SPACE+"%v", strconv.FormatUint(uint64(o.x), 16), o.nn)
	case 0x7000:
		return fmt.Sprintf(ADD+SPACE+V+"%v"+COMMA+SPACE+"%v", strconv.FormatUint(uint64(o.x), 16), o.nn)
	case 0x8000:
		switch o.n {
		case 0x0:
			return fmt.Sprintf(LD+SPACE+V+"%v"+COMMA+SPACE+V+"%v", strconv.FormatUint(uint64(o.x), 16), strconv.FormatUint(uint64(o.y), 16))
		case 0x1:
			return fmt.Sprintf(OR+SPACE+V+"%v"+COMMA+SPACE+V+"%v", strconv.FormatUint(uint64(o.x), 16), strconv.FormatUint(uint64(o.y), 16))
		case 0x2:
			return fmt.Sprintf(AND+SPACE+V+"%v"+COMMA+SPACE+V+"%v", strconv.FormatUint(uint64(o.x), 16), strconv.FormatUint(uint64(o.y), 16))
		case 0x3:
			return fmt.Sprintf(XOR+SPACE+V+"%v"+COMMA+SPACE+V+"%v", strconv.FormatUint(uint64(o.x), 16), strconv.FormatUint(uint64(o.y), 16))
		case 0x4:
			return fmt.Sprintf(ADD+SPACE+V+"%v"+COMMA+SPACE+V+"%v", strconv.FormatUint(uint64(o.x), 16), strconv.FormatUint(uint64(o.y), 16))
		case 0x5:
			return fmt.Sprintf(SUB+SPACE+V+"%v"+COMMA+SPACE+V+"%v", strconv.FormatUint(uint64(o.x), 16), strconv.FormatUint(uint64(o.y), 16))
		case 0x6:
			return fmt.Sprintf(SHR+SPACE+V+"%v", strconv.FormatUint(uint64(o.x), 16))
		case 0x7:
			return fmt.Sprintf(SUBN+SPACE+V+"%v"+COMMA+SPACE+V+"%v", strconv.FormatUint(uint64(o.x), 16), strconv.FormatUint(uint64(o.y), 16))
		case 0x8:
			return fmt.Sprintf(SHL+SPACE+V+"%v", strconv.FormatUint(uint64(o.x), 16))
		default:
			return NOP
		}
	case 0x9000:
		switch o.n {
		case 0x0:
			return fmt.Sprintf(SNE+SPACE+V+"%v"+COMMA+SPACE+V+"%v", strconv.FormatUint(uint64(o.x), 16), strconv.FormatUint(uint64(o.y), 16))
		default:
			return NOP
		}
	case 0xa000:
		return fmt.Sprintf(LD+SPACE+I+COMMA+SPACE+"%v", o.nnn)
	case 0xb000:
		return fmt.Sprintf(JMP+SPACE+V+"0"+COMMA+SPACE+"%v", o.nnn)
	case 0xc000:
		return fmt.Sprintf(RND+SPACE+V+"%v"+COMMA+SPACE+"%v", strconv.FormatUint(uint64(o.x), 16), o.nn)
	case 0xd000:
		return fmt.Sprintf(DRAW+SPACE+V+"%v"+COMMA+SPACE+V+"%v"+COMMA+SPACE+"%v", strconv.FormatUint(uint64(o.x), 16), strconv.FormatUint(uint64(o.y), 16), o.n)
	case 0xe000:
		switch o.nn {
		case 0x9e:
			return fmt.Sprintf(SKP+SPACE+V+"%v", strconv.FormatUint(uint64(o.x), 16))
		case 0xa1:
			return fmt.Sprintf(SKNP+SPACE+V+"%v", strconv.FormatUint(uint64(o.x), 16))
		default:
			return NOP
		}
	case 0xf000:
		switch o.nn {
		case 0x07:
			return fmt.Sprintf(LD+SPACE+V+"%v"+COMMA+SPACE+DT, strconv.FormatUint(uint64(o.x), 16))
		case 0x0a:
			return fmt.Sprintf(LD+SPACE+V+"%v"+COMMA+SPACE+K, strconv.FormatUint(uint64(o.x), 16))
		case 0x15:
			return fmt.Sprintf(LD+SPACE+DT+COMMA+SPACE+V+"%v", strconv.FormatUint(uint64(o.x), 16))
		case 0x18:
			return fmt.Sprintf(LD+SPACE+ST+COMMA+SPACE+V+"%v", strconv.FormatUint(uint64(o.x), 16))
		case 0x1e:
			return fmt.Sprintf(ADD+SPACE+I+COMMA+SPACE+V+"%v", strconv.FormatUint(uint64(o.x), 16))
		case 0x29:
			return fmt.Sprintf(LD+SPACE+F+COMMA+SPACE+V+"%v", strconv.FormatUint(uint64(o.x), 16))
		case 0x33:
			return fmt.Sprintf(LD+SPACE+B+COMMA+SPACE+V+"%v", strconv.FormatUint(uint64(o.x), 16))
		case 0x55:
			return fmt.Sprintf(LD+SPACE+R+COMMA+SPACE+V+"%v", strconv.FormatUint(uint64(o.x), 16))
		case 0x65:
			return fmt.Sprintf(LD+SPACE+V+"%v"+COMMA+SPACE+R, strconv.FormatUint(uint64(o.x), 16))
		default:
			return NOP
		}
	default:
		return NOP
	}
}
