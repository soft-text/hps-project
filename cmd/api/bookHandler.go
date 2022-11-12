package main

import (
	"io/ioutil"
	"net/http"
)

func (app *application) getBook(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://openlibrary.org/search.json?title=carl-sagan")
	if err != nil {
		app.logger.Fatalln(err)
	}
	// Read the response body on the line below
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		app.logger.Fatalln(err)
	}
	// Convert the body to type string
	book := string(body)
	app.logger.Printf(book)
	// Convert to JSON using function writeJSON
	err = app.writeJSON(w, http.StatusOK, body, "book")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

}
