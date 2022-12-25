package main

// #cgo CFLAGS: -g -Wall
// #include <stdlib.h>
// #include "pikchr.h"
import "C"
import (
	"fmt"
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
	var tokens []Token
	for {
		token := lexer.Scan()
		if token.Type == TokenTypeEOF {
			break
		}
		tokens = append(tokens, token)
	}

	for _, token := range tokens {
		fmt.Printf("%+v\n", token)
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
