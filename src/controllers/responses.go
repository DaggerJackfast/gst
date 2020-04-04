package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Data interface{}

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Data    Data `json:"data,omitempty"`
}

type ResponseError struct {
	Error string `json:"error"`
}

// TODO: add interface and usage as method of Response struct
func JSON(w http.ResponseWriter, statusCode int, data Data) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

func ERROR(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		re := ResponseError{Error: err.Error()}
		JSON(w, statusCode, re)
		return
	}
	JSON(w, http.StatusBadRequest, nil)
}
