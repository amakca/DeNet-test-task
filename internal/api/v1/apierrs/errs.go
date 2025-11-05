package apierrs

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	ErrInvalidAuthHeader = fmt.Errorf("invalid auth header")
	ErrCannotParseToken  = fmt.Errorf("cannot parse token")
)

func NewErrorResponseHTTP(w http.ResponseWriter, errStatus int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errStatus)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": message})
}
