package validation

import (
	"testing"
)

func TestValidationRejectsUnknownStatus(t *testing.T) {
	if Status("burning") == nil {
		t.Fatal("unknown status accepted")
	}
}
