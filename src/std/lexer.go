package std

import "fmt"

type TokenType int

const (
	IDENT TokenType = iota
	LAMBDA
	DOT
	LPAREN
	RPAREN
	ASSIGN
	EOF
)

type Token struct {
	Type  TokenType
	Value string
}

// Lex takes an input string and produces a slice of tokens for the parser.
func Lex(input string) ([]Token, error) {
	var tokens []Token
	for i := 0; i < len(input); {
		switch c := input[i]; {
		case c == '\\':
			tokens = append(tokens, Token{Type: LAMBDA})
			i++
		case c == '.':
			tokens = append(tokens, Token{Type: DOT})
			i++
		case c == '(':
			tokens = append(tokens, Token{Type: LPAREN})
			i++
		case c == ')':
			tokens = append(tokens, Token{Type: RPAREN})
			i++
		case c == '=':
			tokens = append(tokens, Token{Type: ASSIGN})
			i++
		case c == ' ' || c == '\t' || c == '\n':
			i++

		case isLetter(c) || isDigit(c) || c == '_':
			start := i
			for i < len(input) && (isLetter(input[i]) || isDigit(input[i]) || input[i] == '_') {
				i++
			}
			tokens = append(tokens, Token{
				Type:  IDENT,
				Value: input[start:i],
			})

		default:
			return nil, fmt.Errorf("unexpected char %q", c)
		}
	}

	tokens = append(tokens, Token{Type: EOF})
	return tokens, nil
}

func isLetter(c byte) bool {
	return ('a' <= c && c <= 'z') ||
		('A' <= c && c <= 'Z') ||
		c == '_'
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
