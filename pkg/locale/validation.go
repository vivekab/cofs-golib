package goliblocale

import "regexp"

const (
	regExPostalCodeUS = `^\d{4,5}(?:[-\s]\d{4})?$`
	regExPostalCode   = `^[a-zA-Z0-9- ]{0,12}$` //a to z, A to Z, 0 to 9, - and space allowed
)

func ValidPostalCodeUS(postalCode string) bool {

	pattern := regexp.MustCompile(regExPostalCodeUS)
	return pattern.MatchString(postalCode)

}

func ValidPostalCode(postalCode string) bool {

	pattern := regexp.MustCompile(regExPostalCode)
	return pattern.MatchString(postalCode)

}
