package main

type Chip8Opcode struct {
	opcode uint16
	c      uint16
	x      uint8
	y      uint8
	n      uint8
	nn     uint8
	nnn    uint16
}

func NewChip8Opcode(o uint16) Chip8Opcode {
	return Chip8Opcode{
		opcode: o,
		c:      (o & 0xf000),
		x:      uint8((o & 0x0f00) >> 8),
		y:      uint8((o & 0x00f0) >> 4),
		n:      uint8(o & 0x000f),
		nn:     uint8(o & 0x00ff),
		nnn:    o & 0x0fff,
	}
}
