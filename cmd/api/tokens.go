package main

import (
	"backend/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"
)

// validUser is variable for all valid user in sign in page
var validUser = models.User{
	ID:    10,
	Email: "admin@here.com",
	// Password need to be a hash. Generated with lib: golang.org/x/crypto/bcrypt
	Password: "$2a$12$4sCixyf7V8BQEsMNTUmp5OWDLom7sVHQ1x5yGlrySvc1jtci5P3Ee", // Generated from https://go.dev/play/p/uKMMCzJWGsW with kurosak1 as password
}

// Credentials is an object for JSON communications
type Credentials struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

// Below is a Signin function that has a purpose of handler which we will call in routes.go
// And has receiver of app and pointer of application
func (app *application) Signin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	// Get JSON r.Body that posted to backend and decode into &creds
	err := json.NewDecoder(r.Body).Decode(&creds)
	// If error happens with the process above. Means that data from frontend not valid
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"))
		return
	}

	// Get actual Username and Password and query to database to check if it is valid
	hashedPassword := validUser.Password

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password))
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"))
		return
	}
	// At this point we have valid user and password

	// Now we create JWT and send it back to frontend
	var claims jwt.Claims
	claims.Subject = fmt.Sprint(validUser.ID)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(1 * time.Hour)) // Expires time = 1 Hour
	claims.Issuer = "mydomain.com"
	claims.Audiences = []string{"mydomain.com"}

	// Generate/create JWT token
	// jwtBytes, because we're getting slice of bytes, is assign value of claims.HMACSign with algorithm of HS256 and
	// we're going to pass it to a slice of byte with "secret" inside app.config.jwt
	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secret))
	if err != nil {
		app.errorJSON(w, errors.New("error signing"))
		return
	}

	// Write JSON file to user
	app.writeJSON(w, http.StatusOK, string(jwtBytes), "response")
}
