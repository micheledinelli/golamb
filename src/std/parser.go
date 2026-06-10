package std

import (
	"fmt"

	c "github.com/micheledinelli/golamb/common"
)

type Parser struct {
	tokens []Token
	pos    int
}

func (p *Parser) current() Token { return p.tokens[p.pos] }
func (p *Parser) advance()       { p.pos++ }

// Parse converts a string into a lambda expression.
// It builds an abstract syntax tree (AST) from the input string,
// which can then be evaluated using the evaluation functions see [Normalize].
func Parse(input string) (c.Expr, error) {
	tokens, err := Lex(input)
	if err != nil {
		return nil, err
	}
	p := &Parser{tokens: tokens}
	return p.parseExpr()
}

func (p *Parser) parseExpr() (c.Expr, error) {
	if p.current().Type == LAMBDA {
		return p.parseLambda()
	}

	left, err := p.parseAtom()
	if err != nil {
		return nil, err
	}

	for {
		switch p.current().Type {
		case IDENT, LPAREN:
			right, err := p.parseAtom()
			if err != nil {
				return nil, err
			}
			left = &c.App{
				Fn:  left,
				Arg: right,
			}
		default:
			return left, nil
		}
	}
}

func (p *Parser) parseAtom() (c.Expr, error) {
	tok := p.current()

	switch tok.Type {
	case IDENT:
		p.advance()
		return &c.Var{Name: tok.Value}, nil

	case LPAREN:
		p.advance()
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if p.current().Type != RPAREN {
			return nil, fmt.Errorf("expected )")
		}
		p.advance()
		return expr, nil

	default:
		return nil, fmt.Errorf("unexpected atom token: %+v", tok)
	}
}

func (p *Parser) parseLambda() (c.Expr, error) {
	p.advance()

	if p.current().Type != IDENT {
		return nil, fmt.Errorf("expected parameter after lambda")
	}
	param := p.current().Value
	p.advance()

	if p.current().Type != DOT {
		return nil, fmt.Errorf("expected . after parameter")
	}
	p.advance()

	body, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	return &c.Abs{
		Param: param,
		Body:  body,
	}, nil
}
