package models

import (
	"encoding/json"
	"log"
	"net/http"
)

// Response - response model
type Response struct {
	Code    int    `json:"-"`
	Message string `json:"-"`
	ReqID   string `json:"-"`

	Payload interface{} `json:"payload" example:"id"`
}

// ErrorResponse - error response model
type ErrorResponse struct {
	Error   bool   `json:"error,omitempty" example:"true"`
	Message string `json:"message,omitempty" example:"error message"`
}

// Send - method of sending a response
func (res *Response) Send(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf8")
	w.WriteHeader(res.Code)

	if res.Message == "" {
		res.Message = http.StatusText(res.Code)
	}

	if res.Payload == nil && res.Code != http.StatusOK {
		res.Payload = ErrorResponse{
			Error:   true,
			Message: res.Message,
		}
	}

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Println(res.ReqID, "Sending response failed:", err)
	}
	log.Println(res.ReqID, "Response Body:", r.URL.Path, r.RemoteAddr, r.Method, res.Payload)
}
