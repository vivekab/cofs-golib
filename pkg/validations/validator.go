package golibvalidations

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	goliblocale "github.com/vivekab/golib/pkg/locale"

	tld "github.com/jpillora/go-tld"
	golibtwilio "github.com/vivekab/golib/pkg/twilio"
	golibtypes "github.com/vivekab/golib/pkg/types"
	"github.com/vivekab/golib/protobuf/protoroot"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	idPattern               = "^[a-z]{3,5}_[0-9a-fA-F]{32}$"
	requiredMessage         = "%s is required"
	FORMAT_VALIDATION_ERROR = "Field <fieldName> is invalid, <fieldName> must be in %s format"
)

func ValidateProtoMessage(ctx context.Context, msg proto.Message, isUpdate bool) *protoroot.GrpcError {

	msgDescriptor := msg.ProtoReflect().Descriptor()
	golibValidationOption := protoroot.E_Validate.TypeDescriptor().Number()

	valueMap := getValueMap(msg)

	// Iterate over all fields in the message
	for i := 0; i < msgDescriptor.Fields().Len(); i++ {
		fieldDesc := msgDescriptor.Fields().Get(i)

		var golibValidations *protoroot.GolibValidation
		ok := false

		proto.RangeExtensions(fieldDesc.Options(), func(typ protoreflect.ExtensionType, i interface{}) bool {
			if golibValidationOption != typ.TypeDescriptor().Number() {
				return true
			}
			golibValidations, ok = i.(*protoroot.GolibValidation)

			return !ok
		})

		if golibValidations != nil {

			value := msg.ProtoReflect().Get(fieldDesc)

			if rslt := ValidateGolibValidations(FieldValidationParams{
				fieldDescriptor: fieldDesc,
				validations:     golibValidations,
				value:           value,
			}, valueMap, isUpdate); rslt != nil {
				return rslt
			}

		}
	}

	solidMsgLevelValidations := proto.GetExtension(msg.ProtoReflect().Descriptor().Options(), protoroot.E_MessageValidate).([]*protoroot.MessageValidation)

	for _, msgValidation := range solidMsgLevelValidations {
		isPassed := EvaluvateGolibExp(msgValidation.Exp, valueMap)
		if !isPassed {
			return getGrpcErrorForGolibMessageValidation(msgValidation)
		}

	}

	return nil
}

func getGrpcErrorForGolibMessageValidation(msgValidation *protoroot.MessageValidation) *protoroot.GrpcError {
	errCode := protoroot.ErrorCode(protoroot.ErrorCode_value[msgValidation.ErrorCode])

	return &protoroot.GrpcError{
		Code:     errCode,
		Message:  msgValidation.Message,
		HttpCode: 400,
	}

}

type FieldValidationParams struct {
	fieldDescriptor protoreflect.FieldDescriptor
	validations     *protoroot.GolibValidation
	value           protoreflect.Value
}

func ValidateGolibValidations(fieldValidation FieldValidationParams, msgValue map[string]string, isUpdate bool) *protoroot.GrpcError {

	fieldDescriptor := fieldValidation.fieldDescriptor
	validations := fieldValidation.validations
	value := fieldValidation.value

	if isUpdate {
		validations.Required = validations.RequiredOnUpdate
	}

	if validations.RequiredIf != "" {
		validations.Required = EvaluvateGolibExp(validations.RequiredIf, msgValue)
	}

	if fieldDescriptor.Cardinality() == protoreflect.Repeated {
		if validations.Required && value.List().Len() == 0 {
			return &protoroot.GrpcError{
				Code:      protoroot.ErrorCode_ERROR_CODE_MIN_ALLOWED_LENGTH_NOT_MET,
				Message:   "invalid length",
				FieldName: string(fieldDescriptor.Name()),
				HttpCode:  400,
			}
		}
		for i := 0; i < value.List().Len(); i++ {
			fieldName := string(fieldDescriptor.Name())
			if err := validateFieldOfKind(FieldValidationParams{
				fieldDescriptor: fieldDescriptor,
				validations:     validations,
				value:           value.List().Get(i),
			}, isUpdate); err != nil {
				// Replace the field name with the index of the repeated field
				err.FieldName = strings.Replace(err.FieldName, fieldName, fmt.Sprintf("%s[%d]", fieldName, i), 1)
				return err
			}
		}
		return nil
	}

	if err := validateFieldOfKind(fieldValidation, isUpdate); err != nil {
		return err
	}

	return nil
}

func validateFieldOfKind(fieldValidation FieldValidationParams, isUpdate bool) *protoroot.GrpcError {
	fieldDescriptor := fieldValidation.fieldDescriptor
	validations := fieldValidation.validations
	value := fieldValidation.value

	switch fieldDescriptor.Kind() {
	case protoreflect.Int64Kind:
		return validateAmount(validations, value.Int(), fieldDescriptor)
	case protoreflect.StringKind:
		return validateString(validations, value.String(), fieldDescriptor)

	case protoreflect.MessageKind:
		return validateMessage(validations, value.Message(), fieldDescriptor, isUpdate)

	case protoreflect.EnumKind:
		return validateEnum(validations, int32(value.Enum()), fieldDescriptor)
	case protoreflect.BytesKind:
		return validateBytes(validations, value.Bytes(), fieldDescriptor)

	}
	return nil
}

func validateAmount(validations *protoroot.GolibValidation, value int64, fieldDescriptor protoreflect.FieldDescriptor) *protoroot.GrpcError {
	if validations.Amount && validations.Required {
		if value <= 0 {
			return &protoroot.GrpcError{
				Code:      protoroot.ErrorCode_ERROR_CODE_INVALID_FIELD,
				Message:   "amount should be greater than zero",
				FieldName: string(fieldDescriptor.Name()),
				HttpCode:  400,
			}
		}
	}
	return nil
}

func validateBytes(validations *protoroot.GolibValidation, value []byte, fieldDescriptor protoreflect.FieldDescriptor) *protoroot.GrpcError {
	if validations.Required && (len(value) == 0) {
		return &protoroot.GrpcError{
			Code:      protoroot.ErrorCode_ERROR_CODE_REQUIRED_FIELD_MISSING,
			Message:   fmt.Sprintf(requiredMessage, fieldDescriptor.Name()),
			FieldName: string(fieldDescriptor.Name()),
			HttpCode:  400,
		}
	}

	return nil
}

func validateString(validations *protoroot.GolibValidation, value string, fieldDescriptor protoreflect.FieldDescriptor) *protoroot.GrpcError {

	if validations.Required && (len(value) == 0) {
		return &protoroot.GrpcError{
			Code:      protoroot.ErrorCode_ERROR_CODE_REQUIRED_FIELD_MISSING,
			Message:   fmt.Sprintf(requiredMessage, fieldDescriptor.Name()),
			FieldName: string(fieldDescriptor.Name()),
			HttpCode:  400,
		}
	}

	if validations.Id && (len(value) > 0 || validations.Required) {
		pattern := idPattern
		if validations.Pattern != "" {
			pattern = validations.Pattern
		}
		match, _ := regexp.MatchString(pattern, value)

		if !match {
			return &protoroot.GrpcError{
				Code:      protoroot.ErrorCode_ERROR_CODE_INVALID_FIELD,
				Message:   fmt.Sprintf("%s is invalid", fieldDescriptor.Name()),
				FieldName: string(fieldDescriptor.Name()),
				HttpCode:  400,
			}
		}
	}

	if validations.Phone && (len(value) > 0 || validations.Required) {
		err := golibtwilio.ValidatePhone(value)
		if err != nil {
			return &protoroot.GrpcError{
				Code:      protoroot.ErrorCode_ERROR_CODE_INVALID_FIELD,
				Message:   "phone is invalid. phone must be in E.164 format.",
				FieldName: string(fieldDescriptor.Name()),
				HttpCode:  400,
			}
		}
	}

	if validations.Email && (len(value) > 0 || validations.Required) {
		isValid := golibtypes.Email(value).IsValidEmail()
		if !isValid {
			return &protoroot.GrpcError{
				Code:      protoroot.ErrorCode_ERROR_CODE_INVALID_FIELD,
				Message:   "email is invalid",
				FieldName: string(fieldDescriptor.Name()),
				HttpCode:  400,
			}
		}
	}

	if validations.Url && (len(value) > 0 || validations.Required) {
		url, err := tld.Parse(value)
		if err != nil || url.Hostname() == "" {
			return &protoroot.GrpcError{
				Code:      protoroot.ErrorCode_ERROR_CODE_INVALID_FIELD,
				Message:   "url is invalid",
				FieldName: string(fieldDescriptor.Name()),
				HttpCode:  400,
			}
		}
	}
	if validations.Website && (len(value) > 0 || validations.Required) {
		url, err := tld.Parse(value)
		if err != nil || url.Hostname() == "" {
			return &protoroot.GrpcError{
				Code:      protoroot.ErrorCode_ERROR_CODE_INVALID_FIELD,
				Message:   "website is invalid",
				FieldName: string(fieldDescriptor.Name()),
				HttpCode:  400,
			}
		}
	}

	if validations.NumberString && (len(value) > 0 || validations.Required) {
		match, _ := regexp.MatchString(`^[0-9]+$`, value)
		if !match {
			return &protoroot.GrpcError{
				Code:      protoroot.ErrorCode_ERROR_CODE_INVALID_FIELD,
				Message:   "invalid number",
				FieldName: string(fieldDescriptor.Name()),
				HttpCode:  400,
			}
		}
	}

	if validations.CountryCode && (len(value) > 0 || validations.Required) {
		_, ok := goliblocale.CountryAlpha2ToAlpha3[goliblocale.Country(value)]
		if !ok {
			return &protoroot.GrpcError{
				Code:      protoroot.ErrorCode_ERROR_CODE_INVALID_FIELD,
				Message:   "invalid country code",
				FieldName: string(fieldDescriptor.Name()),
				HttpCode:  400,
			}
		}
	}

	if len(validations.Pattern) > 0 && (len(value) > 0 || validations.Required) {
		match, _ := regexp.MatchString(validations.Pattern, value)
		if !match {
			return &protoroot.GrpcError{
				Code:      protoroot.ErrorCode_ERROR_CODE_INVALID_FIELD,
				Message:   fmt.Sprintf("value should match pattern %s", validations.Pattern),
				FieldName: string(fieldDescriptor.Name()),
				HttpCode:  400,
			}
		}
	}

	if validations.MinLen > 0 && validations.MinLen == validations.MaxLen && (len(value) > 0 || validations.Required) {
		valid := len(value) == int(validations.MinLen)

		if !valid {
			return &protoroot.GrpcError{
				Code:      protoroot.ErrorCode_ERROR_CODE_INVALID_FIELD,
				Message:   fmt.Sprintf("value must be exactly %d characters", validations.MinLen),
				FieldName: string(fieldDescriptor.Name()),
				HttpCode:  400,
			}
		}
	}

	if validations.MinLen > 0 && (len(value) > 0 || validations.Required) {

		valid := len(value) >= int(validations.MinLen)

		if !valid {
			return &protoroot.GrpcError{
				Code:      protoroot.ErrorCode_ERROR_CODE_INVALID_FIELD,
				Message:   fmt.Sprintf("value must be at least %d characters", validations.MinLen),
				FieldName: string(fieldDescriptor.Name()),
				HttpCode:  400,
			}
		}
	}

	if validations.MaxLen > 0 && (len(value) > 0 || validations.Required) {
		valid := len(value) <= int(validations.MaxLen)

		if !valid {
			return &protoroot.GrpcError{
				Code:      protoroot.ErrorCode_ERROR_CODE_INVALID_FIELD,
				Message:   fmt.Sprintf("value must be at max %d characters", validations.MaxLen),
				FieldName: string(fieldDescriptor.Name()),
				HttpCode:  400,
			}
		}
	}

	return nil
}

func validateEnum(validations *protoroot.GolibValidation, value int32, fieldDescriptor protoreflect.FieldDescriptor) *protoroot.GrpcError {

	if validations.Required && value == 0 {
		return &protoroot.GrpcError{
			Code:      protoroot.ErrorCode_ERROR_CODE_REQUIRED_FIELD_MISSING,
			Message:   fmt.Sprintf("%s is a required field.", fieldDescriptor.Name()),
			FieldName: string(fieldDescriptor.Name()),
			HttpCode:  400,
		}
	}
	return nil
}

// func validateAddress(validations *protoroot.GolibValidation, value *protoroot.Address, fieldDescriptor protoreflect.FieldDescriptor) *protoroot.GrpcError {

// 	if value == nil && validations.Required {
// 		return &protoroot.GrpcError{
// 			Code:      protoroot.ErrorCode_ERROR_CODE_INVALID_FIELD,
// 			Message:   "address is invalid",
// 			FieldName: string(fieldDescriptor.Name()),
// 			HttpCode:  400,
// 		}
// 	} else if value == nil {
// 		return nil
// 	}

// 	err := golibutils.ValidateAddressRequest(context.TODO(), value)

// 	grpcErr := goliberror.GetGrpcError(err)

// 	if err != nil {
// 		return mergeSubObjectError(&protoroot.GrpcError{
// 			Code:      protoroot.ErrorCode_ERROR_CODE_INVALID_FIELD,
// 			Message:   "invalid address",
// 			FieldName: string(fieldDescriptor.Name()),
// 			HttpCode:  400,
// 		}, grpcErr)
// 	}

// 	return nil

// }

func validateMessage(validations *protoroot.GolibValidation, value protoreflect.Message, fieldDescriptor protoreflect.FieldDescriptor, isUpdate bool) (grpcError *protoroot.GrpcError) {

	if validations.Required && !value.IsValid() {
		return &protoroot.GrpcError{
			Code:      protoroot.ErrorCode_ERROR_CODE_REQUIRED_FIELD_MISSING,
			Message:   fmt.Sprintf(requiredMessage, string(fieldDescriptor.Name())),
			FieldName: string(fieldDescriptor.Name()),
			HttpCode:  400,
		}

	}

	if !value.IsValid() {
		return nil
	}

	if validations.DateOfBirth && (value.Interface().(*timestamppb.Timestamp) != nil || validations.Required) {
		if value.Interface().(*timestamppb.Timestamp).AsTime().After(timestamppb.Now().AsTime()) {
			return &protoroot.GrpcError{
				Code:      protoroot.ErrorCode_ERROR_CODE_INVALID_FIELD,
				Message:   "date_of_birth is invalid. date_of_birth must be in YYYY-MM-DD format and must not be a future date.",
				FieldName: string(fieldDescriptor.Name()),
				HttpCode:  400,
			}
		}
	}

	if validations.FormationDate && (value.Interface().(*timestamppb.Timestamp) != nil || validations.Required) {
		if value.Interface().(*timestamppb.Timestamp).AsTime().After(timestamppb.Now().AsTime()) {
			return &protoroot.GrpcError{
				Code:      protoroot.ErrorCode_ERROR_CODE_INVALID_FIELD,
				Message:   "formation_date is invalid. formation_date must be in YYYY-MM-DD format and must not be a future date.",
				FieldName: string(fieldDescriptor.Name()),
				HttpCode:  400,
			}
		}
	}

	messageName := fieldDescriptor.Message().FullName()

	// well known types

	switch messageName {
	case "google.protobuf.StringValue":
		grpcError = validateString(validations, value.Interface().(*wrapperspb.StringValue).Value, fieldDescriptor)
	case "google.protobuf.ListValue":
		listValue := value.Interface().(*structpb.ListValue) // Type assert to ListValue
		for i, v := range listValue.Values {
			switch v.Kind.(type) {
			case *structpb.Value_StringValue:
				// Validate number value
				if grpcError = validateString(validations, v.GetStringValue(), fieldDescriptor); grpcError != nil {
					grpcError.FieldName = fmt.Sprintf("%s[%d]", grpcError.FieldName, i)
					break
				}
			// Add more cases for other value types (e.g., BoolValue, StructValue, etc.)
			default:
				// Handle unsupported types or return an error if necessary
				return &protoroot.GrpcError{
					Code:    protoroot.ErrorCode_ERROR_CODE_INTERNAL_ERROR,
					Message: "Unsupported type in ",
				}
			}
		}
	}

	if grpcError != nil {
		return grpcError
	}

	// Add more message level validations here, keep sub_object validations at the end

	if validations.SubObjectValidate || validations.Array {
		valid := ValidateProtoMessage(context.TODO(), proto.Message(value.Interface()), isUpdate)

		if valid != nil {
			return mergeSubObjectError(&protoroot.GrpcError{
				Code:      protoroot.ErrorCode_ERROR_CODE_INVALID_FIELD,
				Message:   fmt.Sprintf("invalid %s", string(fieldDescriptor.Name())),
				FieldName: string(fieldDescriptor.Name()),
				HttpCode:  400,
			}, valid)
		}

	}

	return nil
}

func mergeSubObjectError(parentErr, subObjectGrpcError *protoroot.GrpcError) *protoroot.GrpcError {

	if subObjectGrpcError == nil {
		return parentErr
	}

	mergedErr := parentErr

	if subObjectGrpcError.Code != 0 {
		mergedErr.Code = subObjectGrpcError.Code
	}

	if subObjectGrpcError.Message != "" {
		mergedErr.Message = subObjectGrpcError.Message
	}

	if subObjectGrpcError.FieldName != "" {
		mergedErr.FieldName = fmt.Sprintf("%s.%s", parentErr.FieldName, subObjectGrpcError.FieldName)
	}

	return mergedErr
}

func getValueMap(msg proto.Message) map[string]string {
	msgDescriptor := msg.ProtoReflect().Descriptor()
	valueMap := make(map[string]string)

	// Iterate over all fields in the message
	for i := 0; i < msgDescriptor.Fields().Len(); i++ {

		fieldDesc := msgDescriptor.Fields().Get(i)
		fieldValue := msg.ProtoReflect().Get(fieldDesc)

		// Single Value
		if fieldDesc.Cardinality() != protoreflect.Repeated {
			valueMap[string(fieldDesc.Name())] = getValueForField(fieldDesc, fieldValue)
			continue
		}

		// Repeated Value
		stringValues := []string{}
		for i := 0; i < fieldValue.List().Len(); i++ {
			value := fieldValue.List().Get(i)
			stringValues = append(stringValues, getValueForField(fieldDesc, value))
		}
		valueMap[string(fieldDesc.Name())] = strings.Join(stringValues, ",")

	}
	return valueMap
}

func getValueForField(fieldDesc protoreflect.FieldDescriptor, value protoreflect.Value) string {
	switch fieldDesc.Kind() {
	case protoreflect.StringKind:
		return value.String()
	case protoreflect.Int64Kind, protoreflect.Int32Kind, protoreflect.Uint32Kind, protoreflect.Uint64Kind:
		return fmt.Sprintf("%d", value.Int())
	case protoreflect.EnumKind:
		return fmt.Sprintf("%d", value.Enum())
	case protoreflect.MessageKind:
		return handleMessageKind(value.Message())
	default:
		return "unknown"
	}
}

func handleMessageKind(msg protoreflect.Message) string {
	if !msg.IsValid() {
		return "null"
	}

	// Check if the message is a well-known type
	switch msg.Descriptor().FullName() {
	case "google.protobuf.ListValue":
		listValue := msg.Interface().(*structpb.ListValue) // Type assert to ListValue
		values := []string{}
		for _, v := range listValue.Values {
			values = append(values, getStructValue(v)) // Recursively handle the value
		}
		return strings.Join(values, ",")
	case "google.protobuf.Struct":
		structMsg := msg.Interface().(*structpb.Struct)
		if structMsg == nil {
			return "null"
		}
		return "valid"
	case "google.protobuf.StringValue":
		return msg.Interface().(*wrapperspb.StringValue).Value
	case "google.protobuf.Int32Value":
		return fmt.Sprintf("%d", msg.Interface().(*wrapperspb.Int32Value).Value)
	case "google.protobuf.Int64Value":
		return fmt.Sprintf("%d", msg.Interface().(*wrapperspb.Int64Value).Value)
	case "google.protobuf.FloatValue":
		return fmt.Sprintf("%f", msg.Interface().(*wrapperspb.FloatValue).Value)
	case "google.protobuf.DoubleValue":
		return fmt.Sprintf("%f", msg.Interface().(*wrapperspb.DoubleValue).Value)
	case "google.protobuf.BoolValue":
		return fmt.Sprintf("%t", msg.Interface().(*wrapperspb.BoolValue).Value)
	default:
		return "valid"
	}
}

// Helper function to handle values inside ListValue (which are of type google.protobuf.Value)
func getStructValue(v *structpb.Value) string {
	switch kind := v.Kind.(type) {
	case *structpb.Value_StringValue:
		return kind.StringValue
	case *structpb.Value_NumberValue:
		return fmt.Sprintf("%f", kind.NumberValue)
	case *structpb.Value_BoolValue:
		return fmt.Sprintf("%t", kind.BoolValue)
	case *structpb.Value_ListValue:
		return "valid"
	case *structpb.Value_NullValue:
		return "null"
	case *structpb.Value_StructValue:
		return "valid"
	default:
		return "valid"
	}
}
