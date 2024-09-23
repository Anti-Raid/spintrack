package structstring

import "reflect"

type StructFieldsConfig struct {
	FieldFilter func(field reflect.StructField) bool
}

// Given a struct, return a list of all the fields in the struct
//
// Returns an empty slice if `s` is not a struct
func StructFields(s any, cfg StructFieldsConfig) []string {
	if s == nil {
		return []string{}
	}

	refType := reflect.TypeOf(s)

	if refType.Kind() != reflect.Struct {
		return []string{}
	}

	var fields = []string{}

	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)

		if cfg.FieldFilter != nil && !cfg.FieldFilter(field) {
			continue
		}

		fields = append(fields, field.Name)
	}

	return fields
}
