package jsontogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
)

// An Encoder writes a JSON element to its Golang struct representation
type Encoder struct {
	w          io.Writer
	Name       string
	Tabulation string
	Tags       []string
}

// NewEncoderWithName returns a new encoder that writes to w with name.
func NewEncoderWithName(w io.Writer, name string) *Encoder {
	return &Encoder{
		w:          w,
		Name:       name,
		Tabulation: "\t",
	}
}

// NewEncoderWithNameAndTags returns a new encoder that writes to w with name and tags.
func NewEncoderWithNameAndTags(w io.Writer, name string, tags []string) *Encoder {
	return &Encoder{
		w:          w,
		Name:       name,
		Tabulation: "\t",
		Tags:       tags,
	}
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

// writeField writes a field to w.
func (enc *Encoder) writeField(fieldName string, fieldType string, level int) {
	enc.writeString(fmt.Sprintf("%s %s", toGoFieldCorrectName(fieldName), fieldType), level)
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
func (enc *Encoder) writeCloseScope(level int) {
	enc.writeString("}\n", level)
}

// browseInterface browses the interface and writes the Go representation to w.
func (enc *Encoder) browseInterface(goTypeMap map[string]interface{}, level int) {
	for k, v := range goTypeMap {
		vType := reflect.TypeOf(v)
		vKind := vType.Kind()

		switch vKind {
		case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			{
				if k != "" { // skip elements with empty names
					enc.writeField(k, vKind.String(), level)
					tagsLen := len(enc.Tags)
					if tagsLen > 0 {
						tagBuffer := bytes.Buffer{}
						tagBuffer.WriteString(" `")
						for i, tag := range enc.Tags {
							tagBuffer.WriteString(fmt.Sprintf(`%s:"%s"`, tag, toGoTagCorrectName(k)))
							if i < tagsLen-1 {
								tagBuffer.WriteString(", ")
							}
						}
						tagBuffer.WriteString("`")
						enc.writeString(tagBuffer.String(), 0)
					}
					enc.writeString("\n", 0)
				}
			}
		case reflect.Slice, reflect.Array:
			{
				if k != "" { // skip elements with empty names
					enc.writeInnerSliceStruct(k, level)
					vSlice := []interface{}(v.([]interface{}))
					for vSliceEntry := range vSlice {
						enc.browseInterface(map[string]interface{}(vSlice[vSliceEntry].(map[string]interface{})), level+1)
						break
					}
					enc.writeCloseScope(level)
				}
			}
		case reflect.Map:
			{
				if k != "" { // skip elements with empty names
					enc.writeInnerStruct(k, level)
					enc.browseInterface(map[string]interface{}(v.(map[string]interface{})), level+1)
					enc.writeCloseScope(level)
				}
			}
		}
	}
}

// Encode reads the JSON element from data and write the Go type to w.
func (enc *Encoder) Encode(data []byte) error {
	if err := enc.checkUp(); err != nil {
		return err
	}
	var goType interface{}
	if err := json.Unmarshal(data, &goType); err != nil {
		return err
	}
	goTypeKind := reflect.TypeOf(goType).Kind()
	switch goTypeKind {
	case reflect.Map:
		{
			enc.writeStruct(enc.Name, 0)
			enc.browseInterface(map[string]interface{}(goType.(map[string]interface{})), 1)
			enc.writeCloseScope(0)
		}
	case reflect.Slice:
		return fmt.Errorf("expecting JSON Array, got JSON Element")
	default:
		return fmt.Errorf("unknown type")
	}
	return nil
}
