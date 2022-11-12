package main

import (
	"encoding/json"
	"net/http"
)

// Simple func that will create JSON and send it to browser
// Arguments:
// w http.ResponseWriter -> somewhere to write it to (example = w)
// status int -> status of code (example = http.StatusOK)
// data interface {} -> the data to convert to JSON (example = movies)
// wrap string -> will wrap our JSON with some kind of key that describes the content coming out of wrap (example = "movies")
func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, wrap string) error {
	wrapper := make(map[string]interface{})

	// Do wrapper method above to:
	// key = wrap
	// value = data
	wrapper[wrap] = data

	// Convert wrapper result to JSON with Marshal function. sekaligus cek error
	js, err := json.Marshal(wrapper)
	if err != nil {
		return err
	}

	// Set JSON header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

// Simple func that will handle improper/error input
// Arguments:
// w http.ResponseWriter -> somewhere to write it to (example = w)
// err error -> takes the error msg if any
// status ...int ->
func (app *application) errorJSON(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}
	type jsonError struct {
		Message string `json:"message"`
	}

	theError := jsonError{
		Message: err.Error(),
	}

	app.writeJSON(w, statusCode, theError, "error")
}

// For bookHandler
// The ideal way is not to use ioutil.ReadAll, but rather use a decoder on the reader directly.
// Here's a nice function that gets a url and decodes its response onto a target structure.

/*
var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
*/
