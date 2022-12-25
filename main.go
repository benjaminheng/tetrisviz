package main

// #cgo CFLAGS: -g -Wall
// #include <stdlib.h>
// #include "pikchr.h"
import "C"
import (
	"log"
	"os"
	"unsafe"
)

func main() {
	err := execute()
	if err != nil {
		log.Fatal(err)
	}
}

func execute() error {
	filename := "examples/board.tetrisviz"
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	lexer := NewLexer(f)
	parser := NewParser(lexer)
	err = parser.Parse()
	if err != nil {
		return err
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
