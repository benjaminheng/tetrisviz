package main

import "errors"

// DiagramConfig contains configuration values that are defined from the
// .tetrisviz file.
type DiagramConfig struct {
	Board struct {
		Width  int
		Height int
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
			if err := i.parseBoardStatement(s); err != nil {
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

func (i *Interpreter) parseBoardStatement(stmt ConfigStmt) error {
	switch stmt.Key {
	case "board":
		// TODO: parse config
	default:
		return errors.New("unsupported config key")
	}
	return nil
}

func (i *Interpreter) parseDiagramStatement(stmt DiagramStmt) error {
	return nil
}
