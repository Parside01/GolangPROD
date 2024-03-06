package validatate

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
)

var (
	regions = []interface{}{"Europe", "Africa", "Americas", "Oceania", "Asia"}
)

func IsValidAlpha2(alpha2 string) bool {
	return validation.Validate(alpha2, validation.Match(regexp.MustCompile("[a-zA-Z]{2}"))) == nil || alpha2 == ""
}

func IsValidRegion(region string) bool {
	return validation.Validate(region, validation.In(regions...)) == nil || region == ""
}
