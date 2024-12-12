package golibtwilio

import (
	"fmt"
	"os"
	"strconv"

	golibconstants "github.com/vivekab/golib/pkg/constants"

	"github.com/ttacon/libphonenumber"
)

// IsDigitPhoneNumber checks only number
func IsDigitOnlyPhoneNumber(phone string) bool {

	for i := 0; i < len(phone); i++ {
		if i == 0 && phone[i] == '+' {
			continue
		}
		if phone[i] >= '0' && phone[i] <= '9' {
			continue
		}
		return false
	}

	return true
}

// ValidatePhone validates phone number
func ValidatePhone(phone string) error {
	if len(phone) < 7 || len(phone) > 16 {
		return fmt.Errorf("invalid phone")
	}

	// in sandbox environments, exclude +15555555555 from validation, as this number is used for testing purposes
	if os.Getenv("API_ENV") != golibconstants.EnvProd && phone == "+15555555555" {
		return nil
	}

	ph, err := libphonenumber.Parse(phone, "")
	if err != nil {
		return err
	}

	phStr := strconv.Itoa(int(ph.GetNationalNumber()))
	countryCode := strconv.Itoa(int(ph.GetCountryCode()))

	fullPhone := countryCode + phStr
	if len(fullPhone) < 6 || len(fullPhone) > 15 {
		return fmt.Errorf("length not satisfied")
	}
	if !IsDigitOnlyPhoneNumber(phone) {
		return fmt.Errorf("not a valid number")
	}
	return nil
}
