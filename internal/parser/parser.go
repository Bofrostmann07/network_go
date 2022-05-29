package parser

import (
	"errors"
	"log"
	"network_go/internal/parser/ast"
	"network_go/internal/parser/lexer"
	"network_go/internal/parser/parser"
)

// ParseQuery parses a query string and returns the corresponding ast.Query.
func ParseQuery(input string) (*ast.Query, error) {
	lex := lexer.NewLexer([]byte(input))
	p := parser.NewParser()
	parsedQuery, err := p.Parse(lex)
	if err != nil {
		return nil, err
	}

	query, ok := parsedQuery.(*ast.Query)
	if !ok {
		log.Printf("Could not cast to query: %v\n", parsedQuery)
		return nil, errors.New("could not cast to query")
	}

	return query, nil
}
