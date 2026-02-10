package errz

import (
	"errors"
	"log"
	"net/http"

	"github.com/Mahaveer86619/Hearth/pkg/views"
	"github.com/gin-gonic/gin"
)

type ErrzType string

const (
	BadRequest          ErrzType = "bad_request"
	NotFound            ErrzType = "not_found"
	Conflict            ErrzType = "conflict"
	Unauthorized        ErrzType = "unauthorized"
	Forbidden           ErrzType = "forbidden"
	InternalServerError ErrzType = "internal_server_error"
)

var statusMap = map[ErrzType]int{
	BadRequest:          http.StatusBadRequest,
	NotFound:            http.StatusNotFound,
	Conflict:            http.StatusConflict,
	Unauthorized:        http.StatusUnauthorized,
	Forbidden:           http.StatusForbidden,
	InternalServerError: http.StatusInternalServerError,
}

type HearthError struct {
	Type    ErrzType
	Message string
	Err     error
}

func (e *HearthError) Error() string {
	return e.Message
}

func New(errType ErrzType, msg string, err error) error {
	return &HearthError{
		Type:    errType,
		Message: msg,
		Err:     err,
	}
}

// HandleErrors now takes *gin.Context to write the response
func HandleErrors(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var bkErr *HearthError
	if errors.As(err, &bkErr) {
		statusCode := statusMap[bkErr.Type]
		if statusCode == 0 {
			statusCode = http.StatusInternalServerError
		}

		if statusCode == http.StatusInternalServerError {
			log.Printf("Internal Error: %v | Source: %v", bkErr.Message, bkErr.Err)
		}

		resp := &views.Failure{
			StatusCode: statusCode,
			Message:    bkErr.Message,
		}

		_ = resp.Send(c)
		return
	}

	log.Printf("Unknown Error: %v", err)
	resp := &views.Failure{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal Server Error",
	}
	_ = resp.Send(c)
}
