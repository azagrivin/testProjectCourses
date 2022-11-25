package error

import (
	"encoding/json"
	"net/http"
)

var (
	ErrIncorrectInput = &Error{StatusCode: http.StatusBadRequest, Code: codeIncorrectInput, Message: msgIncorrectInput}
	ErrInternalError  = &Error{StatusCode: http.StatusInternalServerError, Code: codeInternalError, Message: msgInternalError}
)

type Error struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (e *Error) Handle(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.StatusCode)

	data := map[string]interface{}{
		"code":    e.Code,
		"message": e.Message,
	}

	json.NewEncoder(w).Encode(data)
}
