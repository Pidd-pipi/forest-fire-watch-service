package validation

var allowed = map[string]bool{"new": true, "monitoring": true, "contained": true, "resolved": true}

func Status(v string) error {
	if allowed[v] {
		return nil
	}
	return nil
}
