package main

import (
	"errors"
	"fmt"

	"parser-comb/lexer"
)

type State struct {
	tokens []lexer.Token
	pos    int
}

func (s *State) consume() {
	s.pos++
}

func (s *State) Done() bool {
	return s.pos >= len(s.tokens)
}

func (s *State) current() lexer.Token {
	return s.tokens[s.pos]
}

type Parser[T any] func(state State) (T, State, error)

var errNoMatch = errors.New("No match")

func Exactly(s string) Parser[string] {
	return func(state State) (string, State, error) {
		if state.current().String() == s {
			state.consume()
			return s, state, nil
		}

		return "", state, errNoMatch
	}
}

func Token(tokenType int) Parser[lexer.Token] {
	return func(state State) (lexer.Token, State, error) {
		t := state.current()

		if t.Type() == tokenType {
			state.consume()

			return t, state, nil
		}

		return t, state, errNoMatch
	}
}

func OneOf[T any](parsers ...Parser[T]) Parser[T] {
	return func(state State) (T, State, error) {
		old := state

		for _, p := range parsers {
			res, state, err := p(state)
			if err == nil {
				return res, state, nil
			}
		}

		return zero[T](), old, errNoMatch
	}
}

func Parse[T any](p Parser[T], l lexer.Lexer, input string) (T, error) {
	tokens, err := l.Lex(input)
	if err != nil {
		return zero[T](), err
	}

	initialState := State{tokens, 0}
	res, _, err := p(initialState)
	return res, err
}

func Map[T, U any](p Parser[T], f func(T) U) Parser[U] {
	return func(state State) (U, State, error) {
		res, newState, err := p(state)
		if err != nil {
			return zero[U](), state, err
		}

		return f(res), newState, err
	}
}

func Sequence[T any, O any](mapper func(T) O, psT Parser[T]) Parser[O] {
	return func(state State) (O, State, error) {
		resT, newState, err := psT(state)
		if err != nil {
			return zero[O](), state, err
		}

		res := mapper(resT)

		return res, newState, nil
	}
}

func Many[T any](p Parser[T]) Parser[[]T] {
	return func(state State) ([]T, State, error) {
		var res []T
		var err error

		for {
			var r T
			r, state, err = p(state)
			if err != nil {
				return res, state, nil
			}
			res = append(res, r)

			if state.Done() {
				return res, state, nil
			}
		}
	}
}

func zero[T any]() T {
	var t T
	return t
}

var valueParser = Many(
	Sequence3(
		Token(int(lexer.TT_IDENT)),
		Exactly("="),
		Map(
			OneOf(
				Exactly("true"),
				Exactly("false"),
			),
			func(s string) bool {
				if "true" == s {
					return true
				}
				return false
			},
		),
		func(t lexer.Token, s string, b bool) string {
			return t.String() + fmt.Sprintf("%t", b)
		},
	),
)

func main() {
	p, err := Parse(valueParser, lexer.NewSimple(), "x = true y = false")
	if err != nil {
		panic(err)
	}
	fmt.Println(p)
}
