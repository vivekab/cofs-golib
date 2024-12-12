package golibtypes

import (
	"fmt"
	"net/mail"

	"github.com/vivekab/golib/protobuf/protoroot"
)

// var regEx = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[^-][A-Za-z0-9-]+(\\.[A-Za-z0-9-]+)*(\\.[A-Za-z]{2,})$"

type Email string

func (e Email) String() string {
	return string(e)
}

func (e Email) IsValidEmail() bool {
	_, err := mail.ParseAddress(e.String())
	return err == nil
}

func (e Email) IsValid(args map[string]string) *protoroot.GrpcError {

	if _, err := mail.ParseAddress(e.String()); err != nil {
		return &protoroot.GrpcError{
			Code:    protoroot.ErrorCode_ERROR_CODE_INVALID_FIELD,
			Message: fmt.Sprintf(FORMAT_VALIDATION_ERROR, "email"),
		}
	}

	return nil
}
