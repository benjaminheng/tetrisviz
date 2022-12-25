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
		Width  int64
		Height int64
	}
}

type Interpreter struct {
	statements    []any
	diagramConfig DiagramConfig
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

func (i *Interpreter) parseConfigStatement(stmt ConfigStmt) error {
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
	return nil
}
