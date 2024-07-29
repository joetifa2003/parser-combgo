package lexer

import "errors"

type Lexer interface {
	Lex(input string) ([]Token, error)
}

type Token interface {
	String() string
	Type() int
}

type SimpleTokenType int

const (
	TT_KEYWORD SimpleTokenType = iota
	TT_IDENT
	TT_SYMBOL
)

type SimpleToken struct {
	ttype SimpleTokenType
	lit   string
}

func (t SimpleToken) String() string { return t.lit }

func (t SimpleToken) Type() int { return int(t.ttype) }

type SimpleLexer struct {
}

func NewSimple() *SimpleLexer {
	return &SimpleLexer{}
}

func (l *SimpleLexer) Lex(input string) ([]Token, error) {
	runes := []rune(input)

	res := []Token{}

	i := 0
	for i < len(runes) {
		for isWhiteSpace(runes[i]) {
			i++
		}

		if isCharacter(runes[i]) {
			t := SimpleToken{ttype: TT_IDENT}

			for i := i; i < len(runes) && isCharacter(runes[i]); i++ {
				t.lit += string(runes[i])
			}

			i += len(t.lit)

			res = append(res, t)

			continue
		}

		if isSymbol(runes[i]) {
			s := runes[i]

			if s == '=' {
				t := SimpleToken{ttype: TT_SYMBOL}
				t.lit = "="
				i++

				res = append(res, t)

				continue
			}
		}

		return nil, errors.New("cannot lex")
	}

	return res, nil
}

func isCharacter(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z'
}

func isWhiteSpace(r rune) bool {
	return r == ' ' || r == '\r' || r == '\n'
}

func isSymbol(r rune) bool {
	return r == '='
}
