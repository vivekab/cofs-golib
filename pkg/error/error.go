package goliberror

import (
	"context"
	"net/http"
	"time"

	logging "github.com/vivekab/golib/pkg/logging"

	golibconstants "github.com/vivekab/golib/pkg/constants"
	golibcontext "github.com/vivekab/golib/pkg/context"
	goliblocalise "github.com/vivekab/golib/pkg/localise"
	"github.com/vivekab/golib/protobuf/protoroot"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Error struct {
	Origin  string `json:"origin"`
	Message string `json:"message"`
	Type    int    `json:"type"`
	Code    int    `json:"code"`
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(ctx context.Context, grpcErr *protoroot.GrpcError) error {
	locAttributes := goliblocalise.GetErrorAtrributes(ctx, grpcErr.Code)
	if grpcErr.Message == "" {
		grpcErr.Message = locAttributes.Value
	}
	grpcErr.HttpCode = locAttributes.HttpCode
	status, _ := grpcstatus.New(grpccodes.Internal, "").WithDetails(
		grpcErr,
		&protoroot.RequestId{RequestId: golibcontext.GetFromContext(ctx, golibconstants.HeaderRequestID)},
	)
	logging.ErrorD(ctx, "error", status.Err(), logging.Fields{
		"error_code":  grpcErr.Code.String(),
		"message":     grpcErr.Message,
		"sys_message": grpcErr.SysMessage,
		"field_name":  grpcErr.FieldName,
		"http_code":   grpcErr.HttpCode,
	})
	return status.Err()
}

func NewErrorWithCode(ctx context.Context, errCode protoroot.ErrorCode, sysMessage string) error {
	return NewError(ctx, &protoroot.GrpcError{
		Code:       errCode,
		SysMessage: sysMessage,
	})
}

func GetHttpErrCode(err error) grpccodes.Code {
	status := grpcstatus.Convert(err)
	errDetails := status.Details()
	var httpCode int32
	if len(errDetails) > 1 {
		msg, ok := errDetails[1].(*protoroot.GrpcError)
		if ok {
			httpCode = msg.HttpCode
		}
	}
	return grpccodes.Code(httpCode)
}

func GetErrCode(err error) protoroot.ErrorCode {
	status := grpcstatus.Convert(err)
	if status != nil {
		details := status.Details()
		if len(details) > 0 {
			if rootError, ok := details[0].(*protoroot.GrpcError); ok {
				return rootError.Code
			}
		}
	}
	return protoroot.ErrorCode_ERROR_CODE_UNSPECIFIED
}

func GetGrpcError(err error) *protoroot.GrpcError {
	grpcErr := &protoroot.GrpcError{}
	if status, ok := grpcstatus.FromError(err); ok {
		details := status.Details()
		if len(details) > 0 {
			rootError, ok := details[0].(*protoroot.GrpcError)
			if ok {
				return rootError
			}
		}
	} else {
		grpcErr.Message = err.Error()
		grpcErr.Code = protoroot.ErrorCode_ERROR_CODE_INTERNAL_ERROR
	}
	return grpcErr
}

func UpdateGrpcErrorMessage(ctx context.Context, err error, errMsg string) error {
	grpcErr := GetGrpcError(err)
	if grpcErr != nil {
		grpcErr.Message = errMsg
		return NewError(ctx, grpcErr)
	}
	return err
}

// NewHTTPError is a http standard error for api responses
func NewHTTPError(ctx context.Context, httpMethod string, grpcErr *protoroot.GrpcError) *protoroot.HTTPError {
	locAttributes := goliblocalise.GetErrorAtrributes(ctx, grpcErr.Code)
	if grpcErr.Message == "" {
		grpcErr.Message = locAttributes.Value
		grpcErr.HttpCode = locAttributes.HttpCode
	}
	grpcErr.HttpCode = locAttributes.HttpCode

	return &protoroot.HTTPError{
		RequestId: golibcontext.GetFromContext(ctx, golibconstants.HeaderRequestID),
		Method:    httpMethod,
		Status:    int32(grpcErr.HttpCode),
		Error: &protoroot.Error{
			Message:   grpcErr.Message,
			Code:      grpcErr.Code.String(),
			FieldName: grpcErr.FieldName,
		},
		CreatedAt: &timestamppb.Timestamp{
			Seconds: time.Now().UTC().Unix(),
		},
	}
}

// NewHTTPErrorFromError creates a http error from a generic error - optimized for grpc status errors
func NewHTTPErrorFromGrpcError(ctx context.Context, httpMethod string, err error) *protoroot.HTTPError {
	errResponse := &protoroot.HTTPError{
		Status:    http.StatusInternalServerError,
		Method:    httpMethod,
		RequestId: golibcontext.GetFromContext(ctx, golibconstants.HeaderRequestID),
		CreatedAt: &timestamppb.Timestamp{
			Seconds: time.Now().UTC().Unix(),
		},
	}

	if status, ok := grpcstatus.FromError(err); ok {
		details := status.Details()
		if len(details) > 0 {
			rootError, ok := details[0].(*protoroot.GrpcError)
			if ok {
				errResponse.Status = rootError.HttpCode
				errResponse.Error = &protoroot.Error{
					Code:      rootError.Code.String(),
					Message:   rootError.Message,
					FieldName: rootError.FieldName,
				}
			}
		}
	} else {
		errResponse.Error = &protoroot.Error{
			Message: err.Error(),
			Code:    protoroot.ErrorCode_ERROR_CODE_INTERNAL_ERROR.String(),
		}
	}
	return errResponse
}

/*
GetSysMessageFromError return sysMessage embedded in the error returned by gRPC services.
Errors returned by go-lib.NewError which is used by most of gRPC services
returns a gRPC compatible error which can be converted to grpc Status type
*/
func GetSysMessageFromError(err error) string {
	status := grpcstatus.Convert(err)
	errDetails := status.Details()
	sysMessage := ""
	if len(errDetails) > 1 {
		msg, ok := errDetails[0].(*protoroot.GrpcError)
		if ok {
			sysMessage = msg.Message
		}
	}
	return sysMessage
}
