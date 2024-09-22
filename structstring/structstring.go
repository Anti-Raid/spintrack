package structstring

import (
	"fmt"
	"reflect"
	"strings"
)

type ConvertStructToStringConfig struct {
	// Prefixing
	Prefixer func(depth int) string
	// Struct Recursion Override, useful for resolving things like time.Time to time.Time and not expanding it fully out
	StructRecurseOverride func(t reflect.Type) (*string, bool)
}

func NewDefaultConvertStructToStringConfig() *ConvertStructToStringConfig {
	return &ConvertStructToStringConfig{
		Prefixer: func(depth int) string {
			var tabs = ""

			for i := 0; i < depth; i++ {
				tabs += "\t"
			}

			return tabs
		},
		StructRecurseOverride: func(t reflect.Type) (*string, bool) {
			switch t.PkgPath() {
			// Time is self-explanatory
			case "time":
				timeName := "time." + t.Name()
				return &timeName, true
			// RawMessage is a special case for json.RawMessage
			case "encoding/json":
				if t.Name() == "RawMessage" {
					rawMessageName := "json.RawMessage"
					return &rawMessageName, true
				}
			}
			return nil, false
		},
	}
}

func ConvertStructToString(s any, cfg *ConvertStructToStringConfig) string {
	if s == nil {
		return ""
	}

	refType := reflect.TypeOf(s)

	return findStructType(refType, 1, cfg)
}

func findStructType(t reflect.Type, depth int, cfg *ConvertStructToStringConfig) string {
	switch t.Kind() {
	case reflect.Struct:
		// Handle stdlib
		switch t.PkgPath() {
		case "time":
			return "time." + t.Name()
		}

		name := t.Name()

		if name == "" {
			name = "{"
		} else {
			name += " {"
		}

		var fields = []string{}
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)

			jsonTag := field.Tag.Get("json")
			reflectOpts := field.Tag.Get("reflect")

			if reflectOpts == "ignore" || jsonTag == "-" {
				continue
			}

			structName := field.Name

			if jsonTag != "" {
				structName = jsonTag + " (fieldname=" + field.Name + ")"
			}

			fields = append(fields, fmt.Sprintf("%s%v: %v", cfg.Prefixer(depth), structName, findStructType(field.Type, depth+1, cfg)))
		}

		name += "\n" + strings.Join(fields, "\n") + "\n" + cfg.Prefixer(depth-1) + "}"

		return name
	case reflect.Array, reflect.Slice:
		return "[]" + findStructType(t.Elem(), depth, cfg)
	case reflect.Map:
		return "map[" + findStructType(t.Key(), 0, cfg) + "]" + findStructType(t.Elem(), 0, cfg)
	case reflect.Ptr:
		return findStructType(t.Elem(), depth, cfg)
	default:
		return t.Name()
	}
}
