package main

// #cgo CFLAGS: -g -Wall
// #include <stdlib.h>
// #include "pikchr.h"
import "C"
import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"unsafe"

	"golang.org/x/net/html"
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

func getOutputFilename(inputFilename string, format OutputFormat) string {
	basename := strings.TrimSuffix(inputFilename, ".tetrisviz")
	switch format {
	case OutputFormatPikchr:
		return basename + ".pikchr"
	default:
		return basename + ".svg"
	}
}

func execute() error {
	config, err := parseFlags()
	if err != nil {
		return err
	}

	// Read from stdin or from file
	var r io.Reader
	if config.InputFile == "-" {
		r = bufio.NewReader(os.Stdin)
	} else {
		f, err := os.Open(config.InputFile)
		if err != nil {
			return err
		}
		defer f.Close()
		r = f
	}

	// Interpret .tetrisviz data
	lexer := NewLexer(r)
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

	// generate output
	var output string
	switch config.OutputFormat {
	case OutputFormatPikchr:
		output = interpreter.OutputPikchr()
	case OutputFormatSVG:
		output = interpreter.OutputSVG()
		// TODO: put behind a flag
		output, err = applyFixedDimensions(output)
		if err != nil {
			return err
		}
	}

	// write to output file
	if config.InputFile == "-" {
		fmt.Println(output)
	} else {
		outputFilename := getOutputFilename(config.InputFile, config.OutputFormat)
		outputFile, err := os.OpenFile(outputFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return err
		}
		defer outputFile.Close()
		_, err = outputFile.WriteString(output)
		if err != nil {
			return err
		}
		fmt.Printf("success: compiled %s to %s\n", config.InputFile, outputFilename)
	}

	return nil
}

func applyFixedDimensions(svg string) (string, error) {
	doc, err := html.Parse(strings.NewReader(svg))
	if err != nil {
		return "", err
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "svg" {
			for _, a := range n.Attr {
				if a.Key == "viewBox" {
					components := strings.Split(a.Val, " ")
					width, err := strconv.ParseFloat(components[2], 64)
					if err != nil {
						break
					}
					intWidth := int64(math.Ceil(width))
					height, err := strconv.ParseFloat(components[3], 64)
					if err != nil {
						break
					}
					intHeight := int64(math.Ceil(height))
					style := fmt.Sprintf("width: %dpx; height: %dpx;", intWidth, intHeight)
					n.Attr = append(n.Attr, html.Attribute{Key: "style", Val: style})
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	w := &strings.Builder{}
	err = html.Render(w, doc)
	if err != nil {
		return "", err
	}
	svg = strings.TrimPrefix(w.String(), "<html><head></head><body>")
	svg = strings.TrimSuffix(svg, "</body></html>")

	return svg, nil
}

func pikchr(source string) string {
	zText := C.CString(source)
	defer C.free(unsafe.Pointer(zText))
	mFlags := C.uint(0)
	output := C.pikchr(zText, nil, mFlags, nil, nil)
	return C.GoString(output)
}
