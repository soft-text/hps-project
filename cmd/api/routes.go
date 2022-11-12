package main

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// wrap is a function to wrap our middleware so middleware can be used more than one bit in our routes
// wrap has one parameter next. next adl http.Handler
// wrap has return httprouter.Handle
func (app *application) wrap(next http.Handler) httprouter.Handle {
	// inside wrap we have return function which takes http.ResponseWriter and it takes a pointer to request r http.Request
	// and also takes ps which is paramaters from httprouter
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", ps)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (app *application) routes() http.Handler {
	router := httprouter.New()
	secure := alice.New(app.checkToken) // you can put as many pieces of middleware into this variable "secure"

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)

	// Setup route to moviesGraphQL handler that will take care of GraphQL request
	router.HandlerFunc(http.MethodPost, "/v1/graphql", app.moviesGraphQL)

	router.HandlerFunc(http.MethodPost, "/v1/signin", app.Signin)

	router.HandlerFunc(http.MethodGet, "/v1/get_book", app.getBook)
	router.HandlerFunc(http.MethodGet, "/v1/movie/:id", app.getOneMovie)
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.getAllMovies)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:genre_id", app.getAllMoviesByGenre)

	router.HandlerFunc(http.MethodGet, "/v1/genres", app.getAllGenres)

	// Below are protected routes using variable secure we just created
	router.POST("/v1/admin/editmovie", app.wrap(secure.ThenFunc(app.editMovie)))
	router.GET("/v1/admin/deletemovie/:id", app.wrap(secure.ThenFunc(app.deleteMovie)))

	return app.enableCORS(router)
}
