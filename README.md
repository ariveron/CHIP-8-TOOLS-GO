# CHIP-8-TOOLS-GO
CHIP-8 Assembler and Disassembler written in GO

This is a fully working assembler and disassembler, however, it is still a work in progress and more error handling and refactoring is needed.

Run the executable with the `-i` and `-o` flags, input and output files respectively, to assemble a CHIP-8 assembly file into a ROM output file compatible with any CHIP-8 emulator.

Use the `-d` flag to disassemble a rom input file into an assembly language output file.

The `assembler_test.go` file has a test for every instruction available that can be used as a reference.

The `sample` folder contains a small sample program written in this CHIP-8 assembly language; you'll find the original assembly file, ROM, and disassembly of the ROM.
