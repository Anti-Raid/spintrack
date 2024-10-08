package structstring

import (
	"reflect"
)

type StructFieldsConfig struct {
	// For every field, returns the field name if it should be included, or nil if it should be excluded. A *string can also be returned to override the field itself
	FieldFilter func(field reflect.StructField) (*string, bool)
}

// Given a struct, return a list of all the fields in the struct
//
// Returns an empty slice if `s` is not a struct
func StructFields(s any, cfg StructFieldsConfig) []string {
	refType := reflect.TypeOf(s)
	return findType(refType, cfg)
}

func findType(refType reflect.Type, cfg StructFieldsConfig) []string {
	switch refType.Kind() {
	case reflect.Ptr:
		return findType(refType.Elem(), cfg)
	case reflect.Struct:
		var fields = []string{}

		for i := 0; i < refType.NumField(); i++ {
			field := refType.Field(i)

			var fieldName = field.Name
			if cfg.FieldFilter != nil {
				fieldNameOverride, ok := cfg.FieldFilter(field)

				if !ok {
					continue
				}

				if fieldNameOverride != nil {
					fieldName = *fieldNameOverride
				}
			}

			fields = append(fields, fieldName)
		}

		return fields
	default:
		return []string{}
	}
}
