package utility_test

import (
	"testing"

	"github.com/te-th/podca-api/utility"
)

var checkLimitTests = map[string]string{
	"0":            "50",
	"1":            "1",
	"50":           "50",
	"":             "50",
	"-50":          "50",
	"notaninteger": "50",
}

func TestCheckLimit(t *testing.T) {
	for given, expect := range checkLimitTests {
		result := utility.CheckLimit(given)
		if expect != utility.CheckLimit(given) {
			t.Errorf("expected: %s but was %s for %s", expect, result, given)
		}

	}
}
