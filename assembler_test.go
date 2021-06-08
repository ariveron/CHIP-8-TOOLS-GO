package main

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestLineAssembly(t *testing.T) {
	ls := [](struct {
		inst string
		exHi uint8
		exLo uint8
	}){
		{"NOP", 0x00, 0x00},
		{"CLS", 0x00, 0xe0},
		{"RET", 0x00, 0xee},
		{"JMP 1234", 0x14, 0xd2},
		{"JMP V0, 1234", 0xB4, 0xd2},
		{"CALL 1234", 0x24, 0xd2},
		{"SE Vb, 201", 0x3b, 0xc9},
		{"SE Vb, V5", 0x5b, 0x50},
		{"SNE Vb, 201", 0x4b, 0xc9},
		{"SNE Vb, V5", 0x9b, 0x50},
		{"LD Vb, V5", 0x8b, 0x50},
		{"LD Vc, 10", 0x6c, 0x0a},
		{"LD V1, DT", 0xf1, 0x07},
		{"LD Vf, K", 0xff, 0x0a},
		{"LD V9, R", 0xf9, 0x65},
		{"LD I, 4095", 0xaf, 0xff},
		{"LD DT, V0", 0xf0, 0x15},
		{"LD ST, V1", 0xf1, 0x18},
		{"LD F, V2", 0xf2, 0x29},
		{"LD B, V3", 0xf3, 0x33},
		{"LD R, V4", 0xf4, 0x55},
		{"ADD V8, 180", 0x78, 0xb4},
		{"ADD V3, Ve", 0x83, 0xe4},
		{"ADD I, Vc", 0xfc, 0x1e},
		{"OR V5, V5", 0x85, 0x51},
		{"AND V5, V5", 0x85, 0x52},
		{"XOR V5, V5", 0x85, 0x53},
		{"SUB V5, V5", 0x85, 0x55},
		{"SHR V7", 0x87, 06},
		{"SUBN V3, V4", 0x83, 0x47},
		{"SHL V7", 0x87, 0x0e},
		{"RND Vb, 255", 0xcb, 0xff},
		{"SKP V6", 0xe6, 0x9e},
		{"SKNP V4", 0xe4, 0xa1},
		{"DRAW V0, vf, 10", 0xd0, 0xfa},
	}

	for _, l := range ls {
		msg := testLineAssembly(l.inst, l.exHi, l.exLo)
		if msg != "" {
			t.Log(msg)
			t.Fail()
		}
	}
}

func testLineAssembly(inst string, exHi uint8, exLo uint8) string {
	line := fmt.Sprintf(" LABEL: %v ; comment", inst)

	tokens := TokenizeLine(line)
	high, low := AssembleFromLine(tokens)

	if high != exHi || low != exLo {
		sHi, sExHi := hex8String(high), hex8String(exHi)
		sLo, sExLo := hex8String(low), hex8String(exLo)

		var b strings.Builder
		b.WriteString(fmt.Sprintf("\n[+] \"%v\"", line))
		b.WriteString(fmt.Sprintf("\n    - high byte: expected %v, found %v", sExHi, sHi))
		b.WriteString(fmt.Sprintf("\n    -  low byte: expected %v, found %v", sExLo, sLo))
		return b.String()
	} else {
		return ""
	}
}

func hex8String(x uint8) string {
	return "0x" + fmt.Sprintf("%02v", strconv.FormatUint(uint64(x), 16))
}
