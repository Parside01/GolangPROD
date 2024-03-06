package validatate

import (
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
)

func IsValidPaginationParams(l, o string) (int, int, bool) {
	limit, err := strconv.Atoi(l)
	if err != nil {
		return -2, -2, false
	}

	offset, err := strconv.Atoi(o)
	if err != nil {
		return -2, -2, false
	}

	return limit, offset, (validation.Validate(limit, validation.Min(0), validation.Max(50)) == nil) && (validation.Validate(offset, validation.Min(0)) == nil)
}
