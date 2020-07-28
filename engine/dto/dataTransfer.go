package dto

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Msg  string      `json:"message"`
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

//SendResponse function for sending response
func SendResponse(w http.ResponseWriter, r *http.Request, code int, message string, data interface{}) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	res := response{
		Msg:  message,
		Code: code,
		Data: data,
	}

	json.NewEncoder(w).Encode(res)

}
