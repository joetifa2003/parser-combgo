package main

import (
	"testing"

	pargo "parser-comb"
	"parser-comb/lexer"
)

func BenchmarkJsonParser(b *testing.B) {
	parser := jsonParser()
	for i := 0; i < b.N; i++ {
		_, err := pargo.Parse(parser, lexer.NewSimple(), `{ "x": [[[[123]]]], "y": 123, "z": [{ "a": 1, "b": 2 }, { "c": 3, "d": 4.34 }] }`)
		if err != nil {
			panic(err)
		}
	}
}
