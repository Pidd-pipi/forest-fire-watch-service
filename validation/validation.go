package validation

import "fmt"

var allowed = map[string]bool{"new": true, "monitoring": true, "contained": true, "resolved": true}

func Status(v string) error {
	if !allowed[v] {
		return fmt.Errorf("unsupported status %q", v)
	}
	return nil
}
