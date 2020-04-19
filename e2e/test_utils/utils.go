package test_utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func GetResponseBodyJson(response *http.Response) map[string]string {
	bodyBytes := GetResponseBodyBytes(response)
	var resp map[string]string
	err := json.Unmarshal(bodyBytes, &resp)
	if err != nil {
		log.Fatal(err)
	}
	err = response.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func GetResponseBodyString(response *http.Response) string {
	bodyBytes := GetResponseBodyBytes(response)
	err := response.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	resp := string(bodyBytes)
	return resp
}

func GetResponseBodyBytes(response *http.Response) []byte {
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	return bodyBytes
}
