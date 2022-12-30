package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// DiagramConfig contains configuration values that are defined from the
// .tetrisviz file.
type DiagramConfig struct {
	Board struct {
		Width  int64
		Height int64
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

	// TODO: remove, for debugging only
	for _, v := range i.diagram {
		for _, v := range v {
			fmt.Printf("%+v ", string(v))
		}
		fmt.Println()
	}
	return nil
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
		width, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return err
		}
		height, err := strconv.ParseInt(parts[1], 10, 64)
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
	if i.diagramConfig.Board.Width > 0 && int64(len(runes)) > i.diagramConfig.Board.Width {
		return errors.New("diagram exceeds board width")
	} else if i.diagramConfig.Board.Height > 0 && int64(len(i.diagram)) >= i.diagramConfig.Board.Height {
		return errors.New("diagram exceeds board height")
	}
	i.diagram = append(i.diagram, runes)
	return nil
}
