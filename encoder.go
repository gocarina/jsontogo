package jsontogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
)

// An Encoder writes a JSON element to its Golang struct representation.
type Encoder struct {
	w          io.Writer
	Name       string
	Tabulation string
	Tags       []string
}

// NewEncoderWithName returns a new encoder that writes to w with name.
func NewEncoderWithName(w io.Writer, name string) *Encoder {
	return &Encoder{w: w, Name: name, Tabulation: "\t"}

}

// NewEncoderWithNameAndTags returns a new encoder that writes to w with name and tags.
func NewEncoderWithNameAndTags(w io.Writer, name string, tags []string) *Encoder {
	return &Encoder{w: w, Name: name, Tabulation: "\t", Tags: tags}

}

// checkUp checks if the internal state of the Encoder is valid.
func (enc *Encoder) checkUp() error {
	if enc.Name == "" {
		return fmt.Errorf("encoder name cannot be empty")
	}
	return nil
}

// writeString writes the string to w with level of indentation.
func (enc *Encoder) writeString(s string, level int) {
	if level > 0 {
		s = strings.Repeat(enc.Tabulation, level) + s
	}
	enc.w.Write([]byte(s))
}

// writeScalar writes a scalar field with tags to w.
func (enc *Encoder) writeScalar(fieldName string, fieldType string, isSlice bool, level int) {
	if isSlice {
		enc.writeScalarSliceField(fieldName, fieldType, level)
	} else {
		enc.writeScalarField(fieldName, fieldType, level)
	}
	tagsLen := len(enc.Tags)
	if tagsLen > 0 {
		tagBuffer := bytes.Buffer{}
		tagBuffer.WriteString(" `")
		for i, tag := range enc.Tags {
			tagBuffer.WriteString(fmt.Sprintf(`%s:"%s"`, tag, toGoTagCorrectName(fieldName)))
			if i < tagsLen-1 {
				tagBuffer.WriteString(", ")
			}
		}
		tagBuffer.WriteString("`")
		enc.writeString(tagBuffer.String(), 0)
	}
	enc.writeString("\n", 0)
}

// writeField writes a scalar field to w.
func (enc *Encoder) writeScalarField(fieldName string, fieldType string, level int) {
	enc.writeString(fmt.Sprintf("%s %s", toGoFieldCorrectName(fieldName), fieldType), level)
}

// writeField writes a scalar slice field to w.
func (enc *Encoder) writeScalarSliceField(fieldName string, fieldType string, level int) {
	enc.writeString(fmt.Sprintf("%s []%s", toGoFieldCorrectName(fieldName), fieldType), level)
}

// writeStruct writes a struct to w.
func (enc *Encoder) writeStruct(structName string, level int) {
	enc.writeString(fmt.Sprintf("type %s struct {\n", toGoStructCorrectName(structName)), level)
}

// writeInnerStruct writes a nested struct to w.
func (enc *Encoder) writeInnerStruct(structName string, level int) {
	enc.writeString(fmt.Sprintf("%s struct {\n", toGoStructCorrectName(structName)), level)
}

// writeInnerSliceStruct writes a nested slice of struct to w.
func (enc *Encoder) writeInnerSliceStruct(structName string, level int) {
	enc.writeString(fmt.Sprintf("%s []*struct {\n", toGoStructCorrectName(structName)), level)
}

// writeCloseScope writes a closing bracket to w.
func (enc *Encoder) writeCloseScopeStruct(structName string, level int) {
	closeScopeStructBuffer := bytes.Buffer{}
	closeScopeStructBuffer.WriteString("}")
	tagsLen := len(enc.Tags)
	if tagsLen > 0 {
		closeTagsScopeStructBuffer := bytes.Buffer{}
		closeTagsScopeStructBuffer.WriteString(" `")
		for i, tag := range enc.Tags {
			closeTagsScopeStructBuffer.WriteString(fmt.Sprintf(`%s:"%s"`, tag, toGoTagCorrectName(structName)))
			if i < tagsLen-1 {
				closeTagsScopeStructBuffer.WriteString(", ")
			}
		}
		closeTagsScopeStructBuffer.WriteString("`")
		closeScopeStructBuffer.WriteString(closeTagsScopeStructBuffer.String())
	}
	closeScopeStructBuffer.WriteString("\n")
	enc.writeString(closeScopeStructBuffer.String(), level)
}

// writeCloseScope writes a closing bracket to w.
func (enc *Encoder) writeCloseScope(level int) {
	enc.writeString("}\n", level)
}

// encodeMapInterface reads the interface{} and writes its Go representation to w.
func (enc *Encoder) encodeInnerInterface(key string, value interface{}, isInnerSlice bool, level int) {
	valueType := reflect.TypeOf(value)
	valueKind := valueType.Kind()

	switch valueKind {
	case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		{
			if key != "" { // skip elements with empty names
				enc.writeScalar(key, valueKind.String(), false, level)
			}
		}
	case reflect.Slice:
		{
			if key != "" { // skip elements with empty names
				valueSlice := []interface{}(value.([]interface{}))
				for _, innerValueSlice := range valueSlice {
					innerValueSliceKind := reflect.TypeOf(innerValueSlice).Kind()
					if isScalar(innerValueSliceKind) {
						enc.writeScalar(key, innerValueSliceKind.String(), true, level)
					} else {
						enc.writeInnerSliceStruct(key, level)
						enc.encodeInnerInterface(key, innerValueSlice, true, level+1)
						enc.writeCloseScopeStruct(key, level)
					}
					break
				}
			}
		}
	case reflect.Map:
		if key != "" { // skip elements with empty names
			if isInnerSlice {
				enc.encodeMapInterface(value, level)
			} else {
				enc.writeInnerStruct(key, level)
				enc.encodeMapInterface(value, level+1)
				enc.writeCloseScopeStruct(key, level)
			}
		}
	}
}

// encodeMapInterface reads the map[string]interface{} and writes its Go representation to w.
func (enc *Encoder) encodeMapInterface(goJSONInterface interface{}, level int) {
	goMap := map[string]interface{}(goJSONInterface.(map[string]interface{}))
	for key, value := range goMap {
		enc.encodeInnerInterface(key, value, false, level)
	}

}

// encode reads the JSON element from data and write its Go representation to w.
func (enc *Encoder) encode(goJSONInterface interface{}) {
	enc.writeStruct(enc.Name, 0)
	enc.encodeMapInterface(goJSONInterface, 1)
	enc.writeCloseScope(0)
}

// Encode reads the JSON element from data and write its Go representation to w.
func (enc *Encoder) Encode(data []byte) error {
	if err := enc.checkUp(); err != nil {
		return err
	}
	var goJSONInterface interface{}
	if err := json.Unmarshal(data, &goJSONInterface); err != nil {
		return err
	}
	goJSONInterfaceKind := reflect.TypeOf(goJSONInterface).Kind()
	switch goJSONInterfaceKind {
	case reflect.Map:
		enc.encode(goJSONInterface)
	case reflect.Slice:
		return fmt.Errorf("expecting JSON Array, got JSON Element")
	default:
		return fmt.Errorf("unknown type")
	}
	return nil
}
