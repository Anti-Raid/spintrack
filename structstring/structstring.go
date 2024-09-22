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
	// The tags to lookup/show on each struct
	Tags []string
	// Debug mode
	Debug bool
}

func BaseStructRecurseOverride(t reflect.Type) (*string, bool) {
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
}

func BasePrefixer(depth int) string {
	var tabs = ""

	for i := 0; i < depth; i++ {
		tabs += "\t"
	}

	return tabs
}

func NewDefaultConvertStructToStringConfig() *ConvertStructToStringConfig {
	return &ConvertStructToStringConfig{
		Prefixer:              BasePrefixer,
		StructRecurseOverride: BaseStructRecurseOverride,
		Tags:                  []string{"json", "validate", "description"},
	}
}

func ConvertStructToString(s any, cfg *ConvertStructToStringConfig) string {
	if s == nil {
		return ""
	}

	refType := reflect.TypeOf(s)

	return findStructType(refType, 1, make(map[reflect.Type]string), cfg)
}

func findStructType(t reflect.Type, depth int, visited map[reflect.Type]string, cfg *ConvertStructToStringConfig) string {
	if cfg.Debug {
		fmt.Println("findStructType", t, depth)
	}

	switch t.Kind() {
	case reflect.Struct:
		name := t.Name()

		// Handle override and recursion
		override, overrideOk := cfg.StructRecurseOverride(t)

		if override != nil {
			name = *override
		}

		// Check if in visited to avoid infinite recursion
		if cachedValue, haveVisited := visited[t]; haveVisited {
			return cachedValue + " [self-reference]"
		}

		// Mark as visited
		visited[t] = ""

		if overrideOk {
			return name
		}

		// Add braces
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

			fieldVal := fmt.Sprintf("%s%v: %v", cfg.Prefixer(depth), structName, findStructType(field.Type, depth+1, visited, cfg))

			if len(cfg.Tags) > 0 {
				var tagData = []string{}

				for _, tag := range cfg.Tags {
					tagVal := field.Tag.Get(tag)

					if tagVal != "" {
						tagData = append(tagData, tag+"="+tagVal)
					}
				}

				if len(tagData) > 0 {
					fieldVal += " [" + strings.Join(tagData, ", ") + "]"
				}
			}

			fields = append(fields, fieldVal)
		}

		name += "\n" + strings.Join(fields, "\n") + "\n" + cfg.Prefixer(depth-1) + "}"

		visited[t] = name

		return name
	case reflect.Array, reflect.Slice:
		return "[]" + findStructType(t.Elem(), depth, visited, cfg)
	case reflect.Map:
		return "map[" + findStructType(t.Key(), 0, visited, cfg) + "]" + findStructType(t.Elem(), 0, visited, cfg)
	case reflect.Ptr:
		return "*" + findStructType(t.Elem(), depth, visited, cfg)
	default:
		name := t.Name()

		// e.g. nil interfaces
		if name == "" {
			name = "any"
		}

		return name
	}
}
