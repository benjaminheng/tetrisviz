package main

// #cgo CFLAGS: -g -Wall
// #include <stdlib.h>
// #include "pikchr.h"
import "C"
import (
	"errors"
	"flag"
	"fmt"
	"os"
	"unsafe"
)

var (
	flagFormat = flag.String("format", "svg", "Output format. (Options: pikchr, svg)")
)

type OutputFormat int

const (
	OutputFormatPikchr OutputFormat = iota
	OutputFormatSVG
)

type Config struct {
	InputFile    string
	OutputFormat OutputFormat
}

func main() {
	err := execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func parseFlags() (Config, error) {
	c := Config{
		OutputFormat: OutputFormatSVG,
	}
	flag.Parse()

	switch *flagFormat {
	case "pikchr":
		c.OutputFormat = OutputFormatPikchr
	case "svg":
		c.OutputFormat = OutputFormatSVG
	default:
		return c, errors.New("unsupported format")
	}

	if flag.NArg() == 0 {
		return c, errors.New("no input files specified")
	} else if flag.NArg() > 1 {
		return c, errors.New("at most one input can be specified")
	}
	c.InputFile = flag.Arg(0)
	return c, nil
}

func execute() error {
	config, err := parseFlags()
	if err != nil {
		return err
	}

	f, err := os.Open(config.InputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	lexer := NewLexer(f)
	parser := NewParser(lexer)
	statements, err := parser.Parse()
	if err != nil {
		return err
	}

	interpreter := NewInterpreter(statements)
	err = interpreter.Eval()
	if err != nil {
		return err
	}

	switch config.OutputFormat {
	case OutputFormatPikchr:
		fmt.Println(interpreter.OutputPikchr())
	case OutputFormatSVG:
		fmt.Println(interpreter.OutputSVG())
	}

	return nil
}

func pikchr(source string) string {
	zText := C.CString(source)
	defer C.free(unsafe.Pointer(zText))
	mFlags := C.uint(0)
	output := C.pikchr(zText, nil, mFlags, nil, nil)
	return C.GoString(output)
}
