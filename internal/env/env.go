package env

import (
	"os"
	"path/filepath"
)

// Var returns the value of the environment variable `key`.
// If no env var is found, it returns:
//
// - If no extra args are given, ""
// - If 1 extra arg is given, the value of the arg (for defaulting!)
// - If >1 extra args are given, filepath.Join(s)
func Var(key string, s ...string) string {
	env := os.Getenv(key)
	if env != "" {
		return env
	}

	if s == nil {
		return ""
	}
	if len(s) == 0 {
		return ""
	}

	if len(s) == 1 {
		return s[0]
	}

	return filepath.Join(s...)
}
