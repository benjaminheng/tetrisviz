package main

import (
	"errors"
	"fmt"
)

type ConfigStmt struct {
	Key   string
	Value string
}

type DiagramStmt struct{}

type Parser struct {
	Lexer        *Lexer
	statements   []any
	buffer       []Token
	numUnscanned int
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{
		Lexer: lexer,
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

	token := p.Lexer.Scan()

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
	// Unscanning is implemented by moving the cursor backwards. The scan
	// function will read from the buffer.
	p.numUnscanned++
}

func (p *Parser) Parse() error {
	configurationAllowed := true
	var statements []any
	for {
		token := p.scan()
		if token.Type == TokenTypeEOF {
			break
		}

		// Configuration statements are only allowed at the beginning
		// of the file before any other statements.
		if configurationAllowed && token.Type == TokenTypeConfig {
			p.unscan()
			configStmt, err := p.parseConfig()
			if err != nil {
				return err
			}
			statements = append(statements, configStmt)
		}
	}

	// TODO: remove, debugging
	for _, statement := range statements {
		fmt.Printf("statement = %+v\n", statement)
	}
	return nil
}

func (p *Parser) parseConfig() (ConfigStmt, error) {
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
