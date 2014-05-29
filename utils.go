package jsontogo

import (
	"regexp"
	"bytes"
	"strings"
)

const (
	REPLACE_STRING = "_"
)

var firstLetterRegexp *regexp.Regexp = regexp.MustCompile("[^A-Za-z_]")
var nonValidIdentifiers *regexp.Regexp = regexp.MustCompile(`[^A-Za-z0-9_]`)

// toGoStructCorrectName converts structName to a correct go struct name.
func toGoStructCorrectName(structName string) string {
	return toGoFieldCorrectName(structName)
}

// toGoFieldCorrectName converts fieldName to a correct go field name.
func toGoFieldCorrectName(fieldName string) string {
	fieldNameBuffer := bytes.Buffer{}
	firstLetter := fieldName[:1]
	if firstLetterRegexp.MatchString(firstLetter) { // First name is not a letter
		fieldNameBuffer.WriteString(REPLACE_STRING + firstLetter)
	} else {
		fieldNameBuffer.WriteString(strings.ToUpper(firstLetter))
	}
	fieldNameBuffer.WriteString(fieldName[1:])
	return nonValidIdentifiers.ReplaceAllString(fieldNameBuffer.String(), REPLACE_STRING) // Remove non valid identifiers
}

// toGoFieldCorrectName converts tagName to a correct tag name.
func toGoTagCorrectName(tagName string) string {
	return nonValidIdentifiers.ReplaceAllString(tagName, REPLACE_STRING) // Remove non valid identifiers
}
