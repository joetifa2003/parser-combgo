package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"

	pargo "parser-comb"
	"parser-comb/lexer"
)

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

type Expr interface{ eval() int }

type ExprNumber struct{ Value int }

func (e ExprNumber) eval() int { return e.Value }

type BinaryExpression struct {
	Operands []Expr
	Operator string
}

func (b BinaryExpression) eval() int {
	switch b.Operator {
	case "+":
		return reduce(b.Operands, 0, func(a Expr, b int) int { return a.eval() + b })
	case "*":
		return reduce(b.Operands, 1, func(a Expr, b int) int { return a.eval() * b })
	default:
		panic("invalid operator")
	}
}

func numberParser() pargo.Parser[Expr] {
	return pargo.Map(pargo.Some(digitParser()),
		func(digits []string) (Expr, error) {
			val, err := strconv.Atoi(strings.Join(digits, ""))
			if err != nil {
				return ExprNumber{}, err
			}

			return ExprNumber{Value: val}, nil
		},
	)
}

func binaryOperatorParser(left pargo.Parser[Expr], op string) pargo.Parser[Expr] {
	return pargo.Transform(
		pargo.ManySep(
			left,
			pargo.Exactly(op),
		),
		func(exprs []Expr) Expr {
			if len(exprs) == 1 {
				return exprs[0]
			}
			return BinaryExpression{Operands: exprs, Operator: op}
		},
	)
}

func termParser() pargo.Parser[Expr] {
	return binaryOperatorParser(factorParser(), "+")
}

func factorParser() pargo.Parser[Expr] {
	return binaryOperatorParser(parenParser(), "*")
}

func parenParser() pargo.Parser[Expr] {
	return pargo.OneOf(
		pargo.Sequence3(
			pargo.Exactly("("),
			pargo.Lazy(termParser),
			pargo.Exactly(")"),
			func(_ string, Expr Expr, _ string) Expr {
				return Expr
			},
		),
		numberParser(),
	)
}

func reduce[T, U any](arr []T, initial U, f func(T, U) U) U {
	var res U = initial
	for _, v := range arr {
		res = f(v, res)
	}
	return res
}

func main() {
	p, err := pargo.Parse(termParser(), lexer.NewSimple(), "(1 + 1) * 5")
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("output.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		panic(err)
	}

	spew.Dump(p)
	fmt.Println(p.eval())
}
