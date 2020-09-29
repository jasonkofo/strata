package strata

import (
	"strings"
)

func delimit(char string, elems ...string) string {
	sql := ""
	for i, elem := range elems {
		if i > 0 {
			sql += char
		}
		sql += elem
	}
	return sql
}

func delimitQuoted(char string, elems ...string) string {
	sql := ""
	for i, elem := range elems {
		if i > 0 {
			sql += char
		}
		sql += insertDoubleQuotes(elem)
	}
	return sql
}

func delimitDot(elems ...string) string {
	return delimit(".", elems...)
}

func chainSelector(elems ...string) string {
	return delimitQuoted(".", elems...)
}

func delimitSpace(elems ...string) string {
	return delimit(" ", elems...)
}

func surroundWithSpaces(word string) string {
	return " " + word + " "
}

func insertDoubleQuotes(selector string) string {
	return "\"" + escapeLiterals(selector, "\"") + "\""
}

func insertSingleQuotes(selector string) string {
	return "'" + escapeLiterals(selector, "'") + "'"
}

func escapeLiterals(selector string, args ...string) string {
	out := ""
	for _, arg := range args {
		out = strings.ReplaceAll(selector, arg, "")
	}
	return out
}

func cleanString(name string) string {
	str := strings.ToLower(name)
	str = strings.ReplaceAll(str, " ", "")
	return str
}

func containsQuoted(outerString, innerString string) bool {
	if !strings.Contains(outerString, innerString) {
		return false
	}

	lowerBound := strings.Index(outerString, innerString) - 1
	upperBound := lowerBound + len(innerString) + 1

	if lowerBound == -1 || (lowerBound <= 1 && lowerBound > len(outerString)-2) {
		return false
	}

	return outerString[lowerBound] == '"' && outerString[upperBound] == '"'
}
