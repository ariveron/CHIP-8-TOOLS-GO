package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	iFPath := flag.String("i", "rom", "Input file")
	oFPath := flag.String("o", "rom.c8asm", "Output file")
	disAsm := flag.Bool("d", false, "Set to run disassembler, assembler is run by default")
	requiredFlags := []string{"i", "o"}
	flag.Parse()
	err := ValidateFlags(requiredFlags)
	if err != nil {
		log.Fatal("Error with command line flag(s), ", err)
	}

	iFile, err := os.Open(*iFPath)
	if err != nil {
		log.Fatal("Unable to open input file: ", err)
	}
	defer func() {
		iFile.Close()
	}()

	oFile, err := os.Create(*oFPath)
	if err != nil {
		log.Fatal("Unable to open output file: ", err)
	}
	defer func() {
		oFile.Close()
	}()

	if *disAsm {
		err = RomToAssemblyFile(iFile, oFile)
		if err != nil {
			log.Fatal("Error running disassembler: ", err)
		}
	} else {
		err = AssemblyFileToRom(iFile, oFile)
		if err != nil {
			log.Fatal("Error running assembler: ", err)
		}
	}
}

func ValidateFlags(required []string) error {
	args := make(map[string](bool))

	flag.Visit(func(f *flag.Flag) { args[f.Name] = true })

	var notFound []string
	for _, req := range required {
		if !args[req] {
			notFound = append(notFound, req)
		}
	}

	var err error = nil
	if len(notFound) > 0 {
		err = errors.New("required flag(s) not set: " + fmt.Sprint(notFound))
	}

	return err
}
