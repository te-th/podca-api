package utility

import (
	"testing"
)

var checkLimitTests = map[string]string {
	"0" : "50",
	"1":"1",
	"50":"50",
	"":"50",
	"-50":"50",
	"notaninteger":"50",
}

func TestCheckLimit(t *testing.T) {
	for given, expect := range checkLimitTests {
		result := CheckLimit(given)
		if expect !=CheckLimit(given) {
			t.Errorf("expected: %s but was %s for %s", expect, result, given)
		}

	}
}