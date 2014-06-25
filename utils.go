package jsontogo

import (
	"bytes"
	"reflect"
	"regexp"
	"strings"
)

const (
	camelCaseReplace = "_"
	replaceString = "_"
	replaceFirstLetter = "A"
	noNameStruct = "NameEmpty"
)

var firstLetterValid = regexp.MustCompile(`[A-Za-z_]`)
var nonValidIdentifiers = regexp.MustCompile(`[^A-Za-z0-9_]`)
var tagNonValidIdentifiers = regexp.MustCompile(`[^A-Za-z0-9_-]`)

// toGoStructCorrectName converts structName to a correct go struct name.
func toGoStructCorrectName(structName string) string {
	return toGoFieldCorrectName(structName)
}

// toGoFieldCorrectName converts fieldName to a correct go field name.
func toGoFieldCorrectName(fieldName string) string {
	correctFieldName := nonValidIdentifiers.ReplaceAllString(fieldName, replaceString) // Remove non valid identifiers
	correctFieldNameBuffer := bytes.Buffer{}
	correctFieldNameCamelCase := strings.Split(correctFieldName, camelCaseReplace)

	for _, correctFieldNamePiece := range correctFieldNameCamelCase {
		if len(correctFieldNamePiece) < 1 {
			continue
		}
		firstLetterFieldNamePiece := strings.ToUpper(correctFieldNamePiece[:1])
		glueFieldNamePiece := correctFieldNamePiece[1:]
		if correctFieldNameBuffer.Len() == 0 && !firstLetterValid.MatchString(firstLetterFieldNamePiece) {
			correctFieldNameBuffer.WriteString(replaceFirstLetter)
			correctFieldNameBuffer.WriteString(glueFieldNamePiece)
		} else {
			correctFieldNameBuffer.WriteString(firstLetterFieldNamePiece)
			correctFieldNameBuffer.WriteString(glueFieldNamePiece)
		}
	}
	if correctFieldNameBuffer.Len() == 0 {
		return noNameStruct
	}
	return correctFieldNameBuffer.String()
}

// toGoTagCorrectName converts tagName to a correct tag name.
func toGoTagCorrectName(tagName string) string {
	return tagNonValidIdentifiers.ReplaceAllString(tagName, replaceString) // Remove non valid identifiers
}

// isScalar returns if the specified kind is scalar
func isScalar(kind reflect.Kind) bool {
	switch kind {
	case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return true
	}
	return false
}
