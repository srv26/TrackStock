package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Parse request body and unmarshal the json format to the normal format
func ParseBody(r *http.Request, x interface{}) error {
	if body, err1 := ioutil.ReadAll(r.Body); err1 == nil {
		if err := json.Unmarshal([]byte(body), x); err != nil {
			return err
		}
		return nil
	} else {
		return err1
	}
}

// Sends a json error response
func ResponseWithError(response http.ResponseWriter, code int, payload interface{}) {
	ResponseWithJson(response, code, payload)
}

//Sends a json response with valid data
func ResponseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Contenet-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
