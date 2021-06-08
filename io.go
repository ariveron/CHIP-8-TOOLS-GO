package main

import (
	"bufio"
	"errors"
	"io"
)

type Chip8RomReader bufio.Reader

func (r *Chip8RomReader) GetNextOpcode() (Chip8Opcode, error) {
	b := make([]uint8, 2)

	n, err := (*bufio.Reader)(r).Read(b)
	if n < 2 || err != nil {
		if err != io.EOF {
			err = errors.New("unable to read opcode")
		}
		return Chip8Opcode{}, err
	}

	return NewChip8Opcode(uint16(b[0])<<8 | uint16(b[1])), nil
}

type Chip8RomWriter bufio.Writer

func (w *Chip8RomWriter) WriteOpcode(o Chip8Opcode) error {
	b := make([]uint8, 2)
	b[0], b[1] = uint8(o.opcode>>8), uint8(o.opcode&0xff)

	n, err := (*bufio.Writer)(w).Write(b)
	if n < 2 || err != nil {
		return errors.New("unable to write opcode")
	}

	return nil
}

func (w *Chip8RomWriter) Flush() {
	(*bufio.Writer)(w).Flush()
}
