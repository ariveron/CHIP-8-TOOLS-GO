package main

const (
	// Special
	COMMA   = ","
	COMMENT = ";"
	LABEL   = ":"
	SPACE   = " "
	DOLLAR  = "$"
	// Identifiers
	V  = "V"  // General register prefix
	I  = "I"  // Address register
	K  = "K"  // Next key press
	F  = "F"  // Used to load font into I
	B  = "B"  // Used to load BCD into I
	R  = "R"  // Used to load registers to/from I
	DT = "DT" // Delay timer
	ST = "ST" // Sound timer
	// Instructions
	NOP  = "NOP"  // No operation
	CLS  = "CLS"  // Clear screen
	RET  = "RET"  // Return from subroutine
	JMP  = "JMP"  // Jump to address
	CALL = "CALL" // Call subroutine
	SE   = "SE"   // Skip next instruction if equal
	SNE  = "SNE"  // Skip next instruction if not equal
	SKP  = "SKP"  // Skip if key pressed
	SKNP = "SKNP" // Skip if key not pressed
	LD   = "LD"   // Load right into left
	ADD  = "ADD"  // Add right into left
	SUB  = "SUB"  // Subtract right from left into left
	SUBN = "SUBN" // Subtract left from right into left
	OR   = "OR"   // Or right into left
	AND  = "AND"  // And right into left
	XOR  = "XOR"  // Xor right into left
	SHR  = "SHR"  // Shift right
	SHL  = "SHL"  // Shift left
	RND  = "RND"  // And random number into left
	DRAW = "DRAW" // Draw sprite from I
)
