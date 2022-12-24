package main

// #cgo CFLAGS: -g -Wall
// #include <stdlib.h>
// #include "pikchr.h"
import "C"
import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unsafe"
)

func main() {
	err := execute()
	if err != nil {
		log.Fatal(err)
	}
}

func execute() error {
	filename := "examples/board.tetris-svg"
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	source, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	// Variables for the `board` config
	var boardModeEnabled bool
	var boardWidth, boardHeight int64

	configStatementsAllowed := true

	lines := strings.Split(string(source), "\n")
	for _, line := range lines {
		// Configuration options start with +
		if strings.HasPrefix(line, "+") {
			if configStatementsAllowed {
				configKey, configValue, err := getConfigValue(line)
				if err != nil {
					return err
				}
				switch configKey {
				case "board":
					x, y, err := validateBoardConfig(configValue)
					if err != nil {
						return err
					}
					boardWidth = x
					boardHeight = y
					boardModeEnabled = true
				}
			} else {
				return errors.New("config statements must only appear at the beginning of the file")
			}
		} else if line != "\n" {
			// Config statements are only allowed at the beginning
			// of the file, excluding newlines. Once we start
			// seeing other statements, we stop reading configs.
			configStatementsAllowed = false

			// TODO: parse our DSL here
			_ = boardModeEnabled
			_ = boardWidth
			_ = boardHeight
		}
	}
	return nil
}

func validateBoardConfig(value string) (x, y int64, err error) {
	xy := strings.SplitN(value, "x", 2)
	if len(xy) != 2 {
		return 0, 0, errors.New("invalid board config value")
	}
	sx := xy[0]
	sy := xy[1]
	x, err = strconv.ParseInt(sx, 10, 64)
	if err != nil {
		return 0, 0, err
	}
	y, err = strconv.ParseInt(sy, 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return x, y, nil
}

func getConfigValue(line string) (key string, value string, err error) {
	if !strings.HasPrefix(line, "+") {
		return "", "", errors.New("config statement not found")
	}
	re := regexp.MustCompile(`\+([a-zA-Z][a-zA-Z0-9_-]+) (.+)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 3 {
		return "", "", errors.New("cannot parse config statement")
	}
	return matches[1], matches[2], nil
}

func pikchr(source string) string {
	zText := C.CString(source)
	defer C.free(unsafe.Pointer(zText))
	mFlags := C.uint(0)
	output := C.pikchr(zText, nil, mFlags, nil, nil)
	return C.GoString(output)
}
