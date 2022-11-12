package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pascaldekloe/jwt"
)

// A function so 2 origin urls able to share resource
// Has 1 argument/parameter and has to return to http.Handler
func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		next.ServeHTTP(w, r)
	})
}

func (app *application) checkToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			//cound set anonymous user

		}

		// headerParts has purpose to check if authHeader is correct
		// will error if non authorized user try to edit
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			app.errorJSON(w, errors.New("invalid auth header"))
			return
		}

		if headerParts[0] != "Bearer" {
			app.errorJSON(w, errors.New("unauthorized - no bearer"))
			return
		}

		token := headerParts[1]

		claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.jwt.secret))
		if err != nil {
			app.errorJSON(w, errors.New("unauthorized - failed HMACCheck"), http.StatusForbidden)
			return
		}

		if !claims.Valid(time.Now()) {
			app.errorJSON(w, errors.New("unauthorized - token expired"), http.StatusForbidden)
			return
		}

		if !claims.AcceptAudience("mydomain.com") {
			app.errorJSON(w, errors.New("unauthorized - invalid audience"), http.StatusForbidden)
			return
		}

		if claims.Issuer != ("mydomain.com") {
			app.errorJSON(w, errors.New("unauthorized - invalid issuer"), http.StatusForbidden)
			return
		}

		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			app.errorJSON(w, errors.New("unauthorized"), http.StatusForbidden)
			return
		}

		log.Println("Valid user:", userID)

		next.ServeHTTP(w, r)
	})
}
