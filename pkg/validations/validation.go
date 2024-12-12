package golibvalidations

import (
	"fmt"
	"net/mail"
	"regexp"
	"strconv"
	"strings"

	tld "github.com/jpillora/go-tld"
	"github.com/ttacon/libphonenumber"
	golibconfig "github.com/vivekab/golib/pkg/config"
	goliblogging "github.com/vivekab/golib/pkg/logging"
)

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

func ValidateOFACPhoneNumber(phone string) error {
	// no phone ofac check for empty string
	if phone == "" {
		return nil
	}
	ph, err := libphonenumber.Parse(phone, "")
	if err != nil {
		return err
	}

	countryCode := strconv.Itoa(int(ph.GetCountryCode()))
	for _, isd := range golibconfig.GetStringSlice("ofacBlockedCountry.isd") {
		if isd == countryCode {
			return fmt.Errorf("ofac blocked country isd not supported: +%v", countryCode)
		}
	}

	return nil
}

func ValidateOFACEmail(emailAdd string) error {
	// no email ofac check for empty string
	if emailAdd == "" {
		return nil
	}

	e, err := mail.ParseAddress(emailAdd)
	if err != nil {
		return err
	}

	return ValidateOFACWebsite(e.Address)
}

func ValidateOFACWebsite(websiteURL string) error {
	// no website ofac check for empty string
	if websiteURL == "" {
		return nil
	}

	if !strings.Contains(websiteURL, "://") {
		goliblogging.Info("scheme not supported, parsing as http address")
		websiteURL = fmt.Sprintf("http://%v", websiteURL)
	}

	u, err := tld.Parse(websiteURL)
	if err != nil {
		return err
	}
	tldNames := strings.Split(u.TLD, ".")
	u.TLD = tldNames[len(tldNames)-1]

	for _, blockedCountry := range golibconfig.GetStringSlice("ofacBlockedCountry.domain") {
		if strings.EqualFold(blockedCountry, u.TLD) {
			return fmt.Errorf("ofac blocked country domain not supported: %v", blockedCountry)
		}
	}

	return nil
}

func ValidateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return err
	}
	return nil
}
