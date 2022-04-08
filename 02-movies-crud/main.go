package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/hidapple/isbn-gen/isbn"
)

type Movie struct {
	ID       string    `json:"id"`
	ISBN     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

var movies []Movie

func generateMovieID() string {
	return uuid.New().String()
}

func generateISBN() string {
	isbn, err := isbn.NewISBN("Spain", "")
	if err != nil {
		return ""
	}
	return isbn.Number()
}

func loadMovies() {
	movies = append(movies, Movie{ID: generateMovieID(), ISBN: generateISBN(), Title: "The Go Programming Language", Director: &Director{FirstName: "Robert", LastName: "Griesemer"}})
	movies = append(movies, Movie{ID: generateMovieID(), ISBN: generateISBN(), Title: "The Go Programming Language 2", Director: &Director{FirstName: "John", LastName: "Doe"}})
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(movies[index])
			return
		}
	}
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.Error(w, "Movie not found", http.StatusNotFound)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = generateMovieID()
	movies = append(movies, movie)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/movies/%s", movie.ID))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
	http.Error(w, "Movie not found", http.StatusNotFound)
}

func main() {
	router := mux.NewRouter()
	loadMovies()
	router.HandleFunc("/movies", getMovies).Methods("GET")
	router.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	router.HandleFunc("/movies", createMovie).Methods("POST")
	router.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	router.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")
	fmt.Println("Server is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
