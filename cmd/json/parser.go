package main

import (
	"errors"
	"strconv"
	"strings"

	pargo "parser-comb"
)

func intParser() pargo.Parser[JsonValue] {
	return pargo.Map(
		pargo.Some(digitParser()),
		func(s []string) (JsonValue, error) {
			value, err := strconv.Atoi(strings.Join(s, ""))
			if err != nil {
				return JsonInt{}, err
			}
			return JsonInt{Value: value}, nil
		},
	)
}

func floatParser() pargo.Parser[JsonValue] {
	return pargo.Sequence3(
		pargo.Some(digitParser()),
		pargo.Exactly("."),
		pargo.Some(digitParser()),
		func(a []string, _ string, b []string) JsonValue {
			float := strings.Join(a, "") + "." + strings.Join(b, "")
			value, err := strconv.ParseFloat(float, 64)
			if err != nil {
				return JsonFloat{}
			}
			return JsonFloat{Value: value}
		},
	)
}

func arrayParser() pargo.Parser[JsonValue] {
	return pargo.Sequence3(
		pargo.Exactly("["),
		pargo.ManySep(
			pargo.Lazy(jsonParser),
			pargo.Exactly(","),
		),
		pargo.Exactly("]"),
		func(_ string, v []JsonValue, _ string) JsonValue {
			return JsonArray{Value: v}
		},
	)
}

func objectParser() pargo.Parser[JsonValue] {
	return pargo.Sequence3(
		pargo.Exactly("{"),
		pargo.ManySep(
			pargo.Sequence3(
				stringParser(),
				pargo.Exactly(":"),
				pargo.Lazy(jsonParser),
				func(key JsonValue, _ string, value JsonValue) Pair {
					return Pair{Key: key.(JsonString).Value, Value: value}
				},
			),
			pargo.Exactly(","),
		),
		pargo.Exactly("}"),
		func(_ string, pairs []Pair, _ string) JsonValue {
			m := make(map[string]JsonValue, len(pairs))
			for _, pair := range pairs {
				m[pair.Key] = pair.Value
			}
			return JsonObject{Value: m}
		},
	)
}

func stringParser() pargo.Parser[JsonValue] {
	return pargo.Sequence3(
		pargo.Exactly(`"`),
		pargo.Many(pargo.OneOf(
			pargo.Exactly(`\"`), // Handle escaped quotes
			pargo.Except(`"`),   // Any character except unescaped quote
		)),
		pargo.Exactly(`"`),
		func(_ string, parts []string, _ string) JsonValue {
			value := strings.Join(parts, "")
			return JsonString{Value: value}
		},
	)
}

func boolParser() pargo.Parser[JsonValue] {
	return pargo.Map(
		pargo.Some(alphaParser()),
		func(parts []string) (JsonValue, error) {
			s := strings.Join(parts, "")
			if s != "true" && s != "false" {
				return JsonBool{}, errors.New("invalid boolean")
			}

			return JsonBool{Value: s == "true"}, nil
		},
	)
}

func nullParser() pargo.Parser[JsonValue] {
	return pargo.Map(
		pargo.Some(alphaParser()),
		func(parts []string) (JsonValue, error) {
			s := strings.Join(parts, "")
			if s != "null" {
				return JsonNull{}, errors.New("invalid null")
			}

			return JsonNull{}, nil
		},
	)
}

func jsonParser() pargo.Parser[JsonValue] {
	return pargo.OneOf[JsonValue](
		floatParser(),
		intParser(),
		stringParser(),
		boolParser(),
		nullParser(),
		arrayParser(),
		objectParser(),
	)
}
