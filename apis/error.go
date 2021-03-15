package apis

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rancher/wrangler/pkg/schemas/validation"
)

type ResponseAPIError struct {
	validation.ErrorCode
	Message string
}

type APIErrorCode interface {
	ErrorCode() *validation.ErrorCode
}

type ResponseError struct {
	RespCode         int
	RespBody         []byte
	ResponseAPIError ResponseAPIError
}

func (e *ResponseError) ErrorCode() *validation.ErrorCode {
	return &validation.ErrorCode{
		Code:   e.ResponseAPIError.Code,
		Status: e.ResponseAPIError.Status,
	}
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("respCode： %d, respBody: %s", e.RespCode, string(e.RespBody))
}

func NewResponseError(respCode int, respBody []byte) error {
	var responseAPIError ResponseAPIError
	if err := json.Unmarshal(respBody, &responseAPIError); err != nil {
		return err
	}
	return &ResponseError{
		RespCode:         respCode,
		RespBody:         respBody,
		ResponseAPIError: responseAPIError,
	}
}

func CodeForError(err error) *validation.ErrorCode {
	if errorCode := APIErrorCode(nil); errors.As(err, &errorCode) {
		return errorCode.ErrorCode()
	}
	return nil
}

func IsNotFound(err error) bool {
	return CodeForError(err).Code == validation.NotFound.Code
}
