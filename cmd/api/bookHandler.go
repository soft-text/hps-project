package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// A Response struct to map the entire response
type Response struct {
	NumFound int    `json:"numFound"`
	Start    int    `json:"start"`
	Docs     []Docs `json:"docs"`
}

// A Docs struct to map every docs to
type Docs struct {
	Key        string      `json:"key"`
	Title      string      `json:"title"`
	AuthorName interface{} `json:"author_name"` // Using interface{} isn't the best practice, because it's actually a slice data type.
}

var olibUrl = "https://openlibrary.org"

// var param = "author"
// var value = "carl sagan"
// var page = 2

// Agar GET request-nya ditentukan secara dinamis (nnti oleh frontend)
func dynamicUrl(param string, value string, page int) string {
	switch param {
	// Search by Title
	case "title":
		if page > 0 {
			value = strings.ReplaceAll(value, " ", "+")
			page := strconv.Itoa(page)
			result := olibUrl + "/search.json?title=" + value + "&page=" + page
			return result
		} else {
			value = strings.ReplaceAll(value, " ", "+")
			result := olibUrl + "/search.json?title=" + value
			return result
		}
	// Search by Author
	case "author":
		if page > 0 {
			value = strings.ReplaceAll(value, " ", "+")
			page := strconv.Itoa(page)
			result := olibUrl + "/search.json?author=" + value + "&page=" + page
			return result
		} else {
			value = strings.ReplaceAll(value, " ", "+")
			result := olibUrl + "/search.json?author=" + value
			return result
		}
	// Search by Author and Title
	default:
		if page > 0 {
			value = strings.ReplaceAll(value, " ", "+")
			page := strconv.Itoa(page)
			result := olibUrl + "/search.json?q=" + value + "&page=" + page
			return result
		} else {
			value = strings.ReplaceAll(value, " ", "+")
			result := olibUrl + "/search.json?q=" + value
			return result
		}
	}
}

// Function below are a fetching feature from HTTP REQ endpoint to HTTP WRITTER endpoint
func (app *application) fetchBook(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(dynamicUrl("title", "harry potter", 3))
	// If any error will exit status 1
	if err != nil {
		app.logger.Fatalln(err)
	}

	// debugger
	app.logger.Println(dynamicUrl("title", "harry potter", 3))

	// The client must close the response body when finished with it
	defer resp.Body.Close()

	// Read the response body on the line below
	// body result is in []byte or slice format
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		app.logger.Fatalln(err)
	}

	// from https://stackoverflow.com/questions/17156371/how-to-get-json-response-from-http-get
	// Unmarshaling the json body
	// book result is in []byte or slice format
	var book Response
	err = json.Unmarshal(body, &book)
	if err != nil {
		app.logger.Fatalln(err)
	}

	// debugger
	//app.logger.Println(body)

	// === PAGINATION ===
	// Karena JSON yg diterima adalah bulk (lgsg semua data)
	// kalo mau bikin pagination harus dimasukin ke database dulu
	// else, kalo mau milih top 10 berarti harus ngebuang JSON selain top 10 tersebut
	// dan tiap ke page-N harus selalu request GET yg mana bikin lambat
	// ====
	// Solusi lain: return JSON sebanyak bulk tp disimpen dicache front end dgn per 10 buku

	// Convert to JSON using function writeJSON
	err = app.writeJSON(w, http.StatusOK, book, "book")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

}
