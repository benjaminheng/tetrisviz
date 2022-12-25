package main

import (
	"errors"
)

type Config struct {
	Board struct {
		Width  int
		Height int
	}
}

type ConfigStmt struct {
	Key   string
	Value string
}

type DiagramStmt struct {
	Value string
}

type Parser struct {
	lexer        *Lexer
	statements   []any
	buffer       []Token
	numUnscanned int
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{
		lexer: lexer,
	}
	return p
}

func (p *Parser) scan() Token {
	// If we have unscanned before, read from the buffer instead.
	if p.numUnscanned > 0 && len(p.buffer) > 0 {
		token := p.buffer[len(p.buffer)-p.numUnscanned]
		p.numUnscanned--
		return token
	}

	token := p.lexer.Scan()

	// Reset buffer when we reach the end of a statement
	if token.Type == TokenTypeEOF || token.Type == TokenTypeEndStmt {
		p.buffer = []Token{}
		p.numUnscanned = 0
	}

	// Append to the buffer so we can unscan later
	p.buffer = append(p.buffer, token)
	return token
}

func (p *Parser) unscan() {
	// Unscanning is implemented by moving the buffer cursor backwards. The
	// scan function will read from the buffer.
	p.numUnscanned++
}

func (p *Parser) Parse() error {
	var statements []any
	for {
		token := p.scan()
		if token.Type == TokenTypeEOF {
			break
		} else if token.Type == TokenTypeConfig {
			p.unscan()
			configStmt, err := p.parseConfigStmt()
			if err != nil {
				return err
			}
			statements = append(statements, configStmt)
			continue
		} else {
			p.unscan()
			diagramStmt, err := p.parseDiagramStmt()
			if err != nil {
				return err
			}
			statements = append(statements, diagramStmt)
		}
	}
	return nil
}

func (p *Parser) parseDiagramStmt() (DiagramStmt, error) {
	var diagramStmt DiagramStmt
	for {
		token := p.scan()
		if token.Type == TokenTypeEOF || token.Type == TokenTypeEndStmt {
			break
		}

		switch token.Type {
		case TokenTypeHyphen, TokenTypeIdent, TokenTypeSpace:
		default:
			return DiagramStmt{}, errors.New("invalid syntax")
		}

		for _, v := range token.Value {
			switch v {
			// These are the valid characters to describe a diagram
			case '-', ' ', 'r', 'g', 'b', 'o', 'y', 'p', 't':
			default:
				return DiagramStmt{}, errors.New("invalid identifier")
			}
			diagramStmt.Value += string(v)
		}
	}
	return diagramStmt, nil
}

func (p *Parser) parseConfigStmt() (ConfigStmt, error) {
	token := p.scan()
	if token.Type != TokenTypeConfig {
		return ConfigStmt{}, errors.New("config statements must start with +")
	}

	token = p.scan()
	if token.Type != TokenTypeIdent {
		return ConfigStmt{}, errors.New("invalid config key")
	}
	config := ConfigStmt{Key: token.Value}

	token = p.scan()
	if token.Type != TokenTypeSpace {
		return ConfigStmt{}, errors.New("invalid syntax")
	}

	token = p.scan()
	if token.Type != TokenTypeIdent {
		return ConfigStmt{}, errors.New("invalid identifier")
	}
	config.Value = token.Value

	token = p.scan()
	if token.Type != TokenTypeEndStmt {
		return ConfigStmt{}, errors.New("invalid syntax")
	}
	return config, nil
}
