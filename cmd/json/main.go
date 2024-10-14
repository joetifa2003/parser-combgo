package main

import (
	"github.com/davecgh/go-spew/spew"

	pargo "parser-comb"
	"parser-comb/lexer"
)

func alphaParser() pargo.Parser[string] {
	return pargo.OneOf(
		pargo.Exactly("a"),
		pargo.Exactly("b"),
		pargo.Exactly("c"),
		pargo.Exactly("d"),
		pargo.Exactly("e"),
		pargo.Exactly("f"),
		pargo.Exactly("g"),
		pargo.Exactly("h"),
		pargo.Exactly("i"),
		pargo.Exactly("j"),
		pargo.Exactly("k"),
		pargo.Exactly("l"),
		pargo.Exactly("m"),
		pargo.Exactly("n"),
		pargo.Exactly("o"),
		pargo.Exactly("p"),
		pargo.Exactly("q"),
		pargo.Exactly("r"),
		pargo.Exactly("s"),
		pargo.Exactly("t"),
		pargo.Exactly("u"),
		pargo.Exactly("v"),
		pargo.Exactly("w"),
		pargo.Exactly("x"),
		pargo.Exactly("y"),
		pargo.Exactly("z"),
	)
}

func digitParser() pargo.Parser[string] {
	return pargo.OneOf(
		pargo.Exactly("0"),
		pargo.Exactly("1"),
		pargo.Exactly("2"),
		pargo.Exactly("3"),
		pargo.Exactly("4"),
		pargo.Exactly("5"),
		pargo.Exactly("6"),
		pargo.Exactly("7"),
		pargo.Exactly("8"),
		pargo.Exactly("9"),
	)
}

type JsonValue interface{ jsonValue() }

type JsonInt struct{ Value int }

func (JsonInt) jsonValue() {}

type JsonFloat struct{ Value float64 }

func (JsonFloat) jsonValue() {}

type JsonArray struct{ Value []JsonValue }

func (JsonArray) jsonValue() {}

type JsonObject struct{ Value map[string]JsonValue }

func (JsonObject) jsonValue() {}

type JsonString struct{ Value string }

func (JsonString) jsonValue() {}

type JsonBool struct{ Value bool }

func (JsonBool) jsonValue() {}

type JsonNull struct{}

func (JsonNull) jsonValue() {}

type Pair struct {
	Key   string
	Value JsonValue
}

func main() {
	p, err := pargo.Parse(jsonParser(), lexer.NewSimple(), `{ "hi": true, "foo": null, "baz": [[[[123, "hi"]]]], "boo": "hello" }`)
	if err != nil {
		panic(err)
	}
	spew.Dump(p)
}
