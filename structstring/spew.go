package structstring

import "github.com/spewerspew/spew"

// Simple opinionated spew-er
func SpewStruct(s any) string {
	return spew.Sdump(s)
}
