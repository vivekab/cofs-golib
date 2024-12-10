package golibtypes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/vivekab/golib/protobuf/protoroot"
)

const (
	ENUM_VALIDATION_ERROR   = "Field <fieldName> is invalid, allowed values are %s. Refer to https://docs.solidfi.com/v2/api-reference for valid <fieldName> values."
	FORMAT_VALIDATION_ERROR = "Field <fieldName> is invalid, <fieldName> must be in %s format"
)

func validateAllowed[T ~string](args map[string]string, value T) *protoroot.GrpcError {

	if value == "" {
		return nil
	}

	if _, ok := args["allowed"]; !ok {
		return nil
	}

	allowed := strings.Split(args["allowed"], "|")

	for _, a := range allowed {
		if a == string(value) {
			return nil
		}
	}

	return &protoroot.GrpcError{
		Code:     protoroot.ErrorCode_ERROR_CODE_INVALID_FIELD,
		Message:  fmt.Sprintf(ENUM_VALIDATION_ERROR, getAllowedValues(args, "")),
		HttpCode: http.StatusBadRequest,
	}
}

func getAllowedValues(args map[string]string, defaultAllowed string) string {
	allowed := defaultAllowed

	if _, ok := args["allowed"]; ok {
		allowed = strings.TrimSpace(args["allowed"])
	}

	allowed = strings.ReplaceAll(allowed, "|", ", ")

	return allowed
}

// Used to construct a GrpcError when an invalid enum is passed
// args: map[string]string - the map of arguments passed to the IsValid function
// allallowedValues: string - all the allowed values for the enum, separated by '|'
func getInvalidEnumError(args map[string]string, allallowedValues string) *protoroot.GrpcError {
	allowedValues := getAllowedValues(args, allallowedValues)

	message := fmt.Sprintf(ENUM_VALIDATION_ERROR, allowedValues)
	if allowedValues == "" {
		message = "<fieldName> is invalid. Refer to https://docs.solidfi.com/v2/api-reference for valid <fieldName> values."
	}

	return &protoroot.GrpcError{
		Code:    protoroot.ErrorCode_ERROR_CODE_INVALID_FIELD,
		Message: message,
	}
}

func ValidateEnum[T ~string, Enum any](args map[string]string, toProto map[T]Enum, value T) *protoroot.GrpcError {
	if allowedErr := validateAllowed(args, value); allowedErr != nil {
		return allowedErr
	}

	args = map[string]string{}

	keys := make([]string, 0, len(toProto))
	for k := range toProto {
		if string(k) != "" {
			keys = append(keys, string(k))
		}
	}

	if _, ok := toProto[value]; !ok {
		return getInvalidEnumError(args, strings.Join(keys, "|"))
	}

	return nil
}
