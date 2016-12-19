package utility

import (
	"strconv"
)

// CheckLimit checks and return if a string is an integer between 0... 50
func CheckLimit(value string) string {
	limit, err := strconv.Atoi(value)
	if err != nil {
		return "50"
	}
	if 0 < limit && limit < 50 {
		return strconv.Itoa(limit)
	}

	return "50"
}
