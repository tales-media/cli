/*
Copyright 2025 shio solutions GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package formatter

import (
	"errors"
	"fmt"
	"io"
	"reflect"

	yaml "sigs.k8s.io/yaml/goyaml.v3"

	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	HumanStructTagKey            = "human"
	HumanStructTagWideOnlyOption = "wideonly"
)

var TableStyle = table.Style{
	Name: "talesctl",
	Box: table.BoxStyle{
		BottomLeft:       "+",
		BottomRight:      "+",
		BottomSeparator:  "+",
		EmptySeparator:   " ",
		Left:             "|",
		LeftSeparator:    "+",
		MiddleHorizontal: "-",
		MiddleSeparator:  "+",
		MiddleVertical:   "|",
		PaddingLeft:      "",
		PaddingRight:     "   ",
		PageSeparator:    "\n",
		Right:            "|",
		RightSeparator:   "+",
		TopLeft:          "+",
		TopRight:         "+",
		TopSeparator:     "+",
		UnfinishedRow:    " ~",
	},
	Color:   table.ColorOptions{},
	Format:  table.FormatOptionsDefault,
	HTML:    table.DefaultHTMLOptions,
	Options: table.OptionsNoBordersAndSeparators,
	Size:    table.SizeOptionsDefault,
	Title:   table.TitleOptionsDefault,
}

type Human struct {
	Wide bool
}

var _ Formatter = &Human{}

func (f *Human) List(w io.Writer, list any) error {
	listVal := reflect.ValueOf(list)
	listType := listVal.Type()
	if listType.Kind() != reflect.Slice {
		return errors.New("input is not a slice")
	}

	itemType := listType.Elem()
	if itemType.Kind() != reflect.Struct {
		return errors.New("input is not a slice of structs")
	}

	fieldNames := make([]string, 0, itemType.NumField())
	header := make(table.Row, 0, itemType.NumField())
	for i := 0; i < itemType.NumField(); i++ {
		fieldStructType := itemType.Field(i)

		if !fieldStructType.IsExported() {
			continue
		}

		fieldName := fieldStructType.Name
		headerName, fieldOptions := ParseStructTag(fieldStructType, HumanStructTagKey)

		if !f.Wide && fieldOptions.Contains(HumanStructTagWideOnlyOption) {
			continue
		}

		if f.shouldIncludeStructFieldType(fieldStructType.Type) {
			fieldNames = append(fieldNames, fieldName)
			header = append(header, headerName)
		}
	}

	rows := make([]table.Row, 0, listVal.Len())
	for i := 0; i < listVal.Len(); i++ {
		itemVal := listVal.Index(i)
		row := make(table.Row, 0, itemType.NumField())
		for _, name := range fieldNames {
			fieldVal := itemVal.FieldByName(name)
			row = append(row, f.getStringValue(fieldVal))
		}
		rows = append(rows, row)
	}

	tbl := table.NewWriter()
	tbl.SetOutputMirror(w)
	tbl.SetStyle(TableStyle)
	tbl.AppendHeader(header)
	tbl.AppendRows(rows)
	tbl.Render()
	return nil
}

func (f *Human) getStringValue(val reflect.Value) string {
	valType := val.Type()
	if _, ok := valType.MethodByName("String"); ok {
		return fmt.Sprintf("%v", val.Interface())
	}

	switch valType.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.String:
		return fmt.Sprintf("%v", val.Interface())
	case reflect.Pointer:
		if !val.IsNil() {
			return f.getStringValue(val.Elem())
		}
	}
	return ""
}

func (f *Human) shouldIncludeStructFieldType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.String:
		return true
	default:
		if _, ok := t.MethodByName("String"); ok {
			return true
		}
		if t.Kind() == reflect.Pointer {
			return f.shouldIncludeStructFieldType(t.Elem())
		}
	}
	return false
}

func (f *Human) Object(w io.Writer, obj any) error {
	enc := yaml.NewEncoder(w)
	return enc.Encode(obj)

	// TODO: fix human (object) formatter
	// return f.writeValue(w, reflect.ValueOf(obj), "")
}

func (f *Human) Error(w io.Writer, err error) error {
	if _, err := fmt.Fprintf(w, "Error: %s\n", err.Error()); err != nil {
		return err
	}
	return nil
}

func (f *Human) writeValue(w io.Writer, val reflect.Value, prefix string) error {
	valType := val.Type()

	// use String() method if available
	if _, ok := valType.MethodByName("String"); ok {
		_, err := fmt.Fprintf(w, "%v", val.Interface())
		return err
	}

	// custom print logic
	switch valType.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.String:
		_, err := fmt.Fprintf(w, "%v", val.Interface())
		return err

	case reflect.Struct:
		// TODO: improve this implementation
		p2 := prefix + "  "
		if prefix != "" {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}

		for i := 0; i < valType.NumField(); i++ {
			fieldStructType := valType.Field(i)
			if !fieldStructType.IsExported() {
				continue
			}

			fieldName, _ := ParseStructTag(fieldStructType, HumanStructTagKey)
			fieldVal := val.Field(i)

			if _, err := fmt.Fprintf(w, "%s%s: ", prefix, fieldName); err != nil {
				return err
			}
			if err := f.writeValue(w, fieldVal, p2); err != nil {
				return err
			}
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}
		return nil

	case reflect.Slice, reflect.Array:
		// TODO: improve this implementation
		p2 := prefix + "  "
		if prefix != "" {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}

		for i := 0; i < val.Len(); i++ {
			if _, err := fmt.Fprintf(w, "%s- ", prefix); err != nil {
				return err
			}
			if err := f.writeValue(w, val.Index(i), p2); err != nil {
				return err
			}
			if i+1 < val.Len() {
				if _, err := fmt.Fprintln(w); err != nil {
					return err
				}
			}
		}
		return nil

	case reflect.Map:
		// TODO: improve this implementation
		p2 := prefix + "  "
		if prefix != "" {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}

		for _, mapKey := range val.MapKeys() {
			mapVal := val.MapIndex(mapKey)
			if _, err := fmt.Fprintf(w, "%s%s: ", prefix, mapKey); err != nil {
				return err
			}
			if err := f.writeValue(w, mapVal, p2); err != nil {
				return err
			}
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}
		return nil

	case reflect.Interface, reflect.Pointer:
		// TODO: guard against loops
		return f.writeValue(w, val.Elem(), prefix)

	default:
		panic(fmt.Sprintf("BUG: %s not allowed", valType.Kind()))
	}
}
