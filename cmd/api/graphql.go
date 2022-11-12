package main

import (
	"backend/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/graphql-go/graphql"
)

// Create package level variable called movies that going to be slice of pointer to models.Movie
var movies []*models.Movie

// GraphQL schema definition. To permit remote users to request whatever data "movie" they want
// For describing the data we're going to send back to frontend after moviesGraphQL func read the requested body
var fields = graphql.Fields{
	"movie": &graphql.Field{
		Type:        movieType,
		Description: "Get movie by id",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			id, ok := p.Args["id"].(int)
			if ok {
				for _, movie := range movies {
					if movie.ID == id {
						return movie, nil
					}
				}
			}
			return nil, nil
		},
	},
	"list": &graphql.Field{
		Type:        graphql.NewList(movieType),
		Description: "Get all movies",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return movies, nil
		},
	},
	"search": &graphql.Field{
		Type:        graphql.NewList(movieType),
		Description: "Search movie by title",
		Args: graphql.FieldConfigArgument{
			"titleContains": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			var theList []*models.Movie
			search, ok := params.Args["titleContains"].(string)
			if ok {
				for _, currentMovie := range movies {
					if strings.Contains(currentMovie.Title, search) {
						log.Println("Found One")
						theList = append(theList, currentMovie)
					}
				}
			}
			return theList, nil
		},
	},
}

// Describes data apa saja yang ada di database di GraphQL service ini
var movieType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Movie",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"year": &graphql.Field{
				Type: graphql.Int,
			},
			"release_date": &graphql.Field{
				Type: graphql.DateTime,
			},
			"runtime": &graphql.Field{
				Type: graphql.Int,
			},
			"rating": &graphql.Field{
				Type: graphql.Int,
			},
			"mpaa_rating": &graphql.Field{
				Type: graphql.String,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updated_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"poster": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

// Create handler and describe the data we're going to handle / sending back to our code
// moviesGraphQL is a function which takes a response writer and a pointer to request
// What moviesGraphQL doing is populating movies variable (data we're going to query)
func (app *application) moviesGraphQL(w http.ResponseWriter, r *http.Request) {
	// Call our backend "All" method to get all movies into slice
	movies, _ = app.models.DB.All()

	// Then read the requested body, which come with JSON format
	q, _ := io.ReadAll(r.Body)
	query := string(q)

	// Then logging it
	log.Println(query)
	// Next step is describing what kind of data we're going to send back to frontend. Go to beginning of code

	// After describing the data. Specify root query
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		app.errorJSON(w, errors.New("failed to create schema"))
		log.Println(err)
		return
	}

	// At this point, need to check what is the program getting in our request
	params := graphql.Params{Schema: schema, RequestString: query}
	resp := graphql.Do(params)
	//log.Println(resp) // resp is the output response of the requested JSON body from frontend
	if len(resp.Errors) > 0 {
		app.errorJSON(w, errors.New(fmt.Sprintf("failed: %+v", resp.Errors)))
	}

	// Convert resp to JSON format
	j, _ := json.MarshalIndent(resp, "", " ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
