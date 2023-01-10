package main

import (
	"errors"
	"strconv"
	"strings"
)

// DiagramConfig contains configuration values that are defined from the
// .tetrisviz file.
type DiagramConfig struct {
	Board struct {
		Width  int
		Height int
	}
}

type Interpreter struct {
	isConfigDisallowed bool
	statements         []any
	diagramConfig      DiagramConfig
	diagram            [][]rune
}

func NewInterpreter(statements []any) *Interpreter {
	i := &Interpreter{
		statements: statements,
	}
	return i
}

func (i *Interpreter) Eval() error {
	for _, stmt := range i.statements {
		switch s := stmt.(type) {
		case ConfigStmt:
			if err := i.parseConfigStatement(s); err != nil {
				return err
			}
		case DiagramStmt:
			if err := i.parseDiagramStatement(s); err != nil {
				return err
			}
		default:
			return errors.New("unknown statement type")
		}
	}

	return nil
}

func (i *Interpreter) OutputSVG() (string, error) {
	return "", nil
}

func (i *Interpreter) OutputPikchr() string {
	output := &PikchrTemplate{}

	// If board height exceeds the number of diagram lines, fill in the
	// empty space above the diagram.
	if i.diagramConfig.Board.Height > 0 && i.diagramConfig.Board.Height > len(i.diagram) {
		numMissingLines := i.diagramConfig.Board.Height - len(i.diagram)
		for j := 0; j < numMissingLines; j++ {
			for k := 0; k < i.diagramConfig.Board.Width; k++ {
				output.Draw('-')
			}
			output.Draw('\n')
		}
	}

	for _, lines := range i.diagram {
		for _, block := range lines {
			output.Draw(block)
		}
		output.Draw('\n')
	}
	return output.String()
}

func (i *Interpreter) parseConfigStatement(stmt ConfigStmt) error {
	if i.isConfigDisallowed {
		return errors.New("config statements can only be defined before the diagram")
	}
	switch stmt.Key {
	case "board":
		parts := strings.Split(stmt.Value, "x")
		if len(parts) != 2 {
			return errors.New("invalid value for board config")
		}
		width, err := strconv.Atoi(parts[0])
		if err != nil {
			return err
		}
		height, err := strconv.Atoi(parts[1])
		if err != nil {
			return err
		}
		i.diagramConfig.Board.Width = width
		i.diagramConfig.Board.Height = height
	default:
		return errors.New("unsupported config key")
	}
	return nil
}

func (i *Interpreter) parseDiagramStatement(stmt DiagramStmt) error {
	// When configuration is still allowed and we see an empty line, don't
	// treat the line as part of a diagram.
	if !i.isConfigDisallowed && stmt.IsEmpty() {
		return nil
	}

	// When we see a non-empty diagram statement for the first time,
	// disallow any further configuration statements.
	if !i.isConfigDisallowed && !stmt.IsEmpty() {
		i.isConfigDisallowed = true
	}
	runes := []rune(stmt.Value)
	if i.diagramConfig.Board.Width > 0 && len(runes) > i.diagramConfig.Board.Width {
		return errors.New("diagram exceeds board width")
	} else if i.diagramConfig.Board.Height > 0 && len(i.diagram) >= i.diagramConfig.Board.Height {
		return errors.New("diagram exceeds board height")
	}
	i.diagram = append(i.diagram, runes)
	return nil
}

type PikchrTemplate struct {
	blockMacros []string
	seenBlocks  map[rune]bool
	statements  []string
}

func (t *PikchrTemplate) addBlockMacro(block rune) {
	if _, ok := t.seenBlocks[block]; ok {
		return
	}

	switch block {
	case 'b': // blue
		t.blockMacros = append(t.blockMacros, `define $b { box "" fill skyblue }`)
	case 'o': // orange
		t.blockMacros = append(t.blockMacros, `define $o { box "" fill 0xfbbb11 }`)
	case 'y': // yellow
		t.blockMacros = append(t.blockMacros, `define $y { box "" fill 0xfff223 }`)
	case 'r': // red
		t.blockMacros = append(t.blockMacros, `define $r { box "" fill 0xf13636 }`)
	case 'g': // green
		t.blockMacros = append(t.blockMacros, `define $g { box "" fill 0x39e572 }`)
	case 't': // teal
		t.blockMacros = append(t.blockMacros, `define $t { box "" fill 0x67edf5 }`)
	case 'p': // purple
		t.blockMacros = append(t.blockMacros, `define $p { box "" fill 0xc936f1 }`)
	case '-': // empty
		t.blockMacros = append(t.blockMacros, `define $e { box "" fill 0xc1c1c1 }`)
	}
}

func (t *PikchrTemplate) Draw(block rune) error {
	t.addBlockMacro(block)
	switch block {
	case 'b', 'o', 'y', 'r', 'g', 't', 'p':
		t.statements = append(t.statements, "$"+string(block))
	case '-':
		t.statements = append(t.statements, "$e")
	case '\n':
		t.statements = append(t.statements, "next")
	}

	// Mark block as seen before
	if t.seenBlocks == nil {
		t.seenBlocks = make(map[rune]bool)
	}
	if _, ok := t.seenBlocks[block]; !ok {
		t.seenBlocks[block] = true
	}
	return nil
}

func (t *PikchrTemplate) String() string {
	template := `boxwid = 0.2
boxht = boxwid
$currLine = 1
define next {
  box invis at (-boxwid, -boxwid*$currLine)
  $currLine = $currLine + 1
}
`
	var b strings.Builder
	b.WriteString(template)
	for _, v := range t.blockMacros {
		b.WriteString(v + "\n")
	}

	b.WriteString("\n")
	for _, v := range t.statements {
		b.WriteString(v + ";")
		if v == "next" {
			b.WriteString("\n")
		}
	}
	return strings.TrimSpace(b.String())
}
