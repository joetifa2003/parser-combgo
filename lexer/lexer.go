package lexer

type Lexer interface {
	Lex(input string) ([]Token, error)
}

type Token interface {
	String() string
	Type() int
}

type SimpleToken struct {
	ttype SimpleTokenType
	lit   string
}

func (t SimpleToken) String() string { return t.lit }

func (t SimpleToken) Type() int { return int(t.ttype) }

type SimpleTokenType int

const (
	TT_CHARACTER SimpleTokenType = iota
)

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

		res = append(res, SimpleToken{TT_CHARACTER, string(runes[i])})
		i++
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
