package jsontogo

import (
	"bytes"
	"reflect"
	"regexp"
	"strings"
)

const (
	replaceString = "_"
)

var firstLetterRegexp = regexp.MustCompile("[^A-Za-z_]")
var nonValidIdentifiers = regexp.MustCompile(`[^A-Za-z0-9_]`)

// toGoStructCorrectName converts structName to a correct go struct name.
func toGoStructCorrectName(structName string) string {
	return toGoFieldCorrectName(structName)
}

// toGoFieldCorrectName converts fieldName to a correct go field name.
func toGoFieldCorrectName(fieldName string) string {
	fieldNameBuffer := bytes.Buffer{}
	firstLetter := fieldName[:1]
	if firstLetterRegexp.MatchString(firstLetter) { // First name is not a letter
		fieldNameBuffer.WriteString(replaceString + firstLetter)
	} else {
		fieldNameBuffer.WriteString(strings.ToUpper(firstLetter))
	}
	fieldNameBuffer.WriteString(fieldName[1:])
	return nonValidIdentifiers.ReplaceAllString(fieldNameBuffer.String(), replaceString) // Remove non valid identifiers
}

// toGoFieldCorrectName converts tagName to a correct tag name.
func toGoTagCorrectName(tagName string) string {
	return nonValidIdentifiers.ReplaceAllString(tagName, replaceString) // Remove non valid identifiers
}

// isScalar returns if the specified kind is scalar
func isScalar(kind reflect.Kind) bool {
	switch kind {
	case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return true
	}
	return false
}
