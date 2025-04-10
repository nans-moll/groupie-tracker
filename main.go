package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Movie struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Year  string `json:"year"`
}

type TMDbMovie struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Year  string `json:"release_date"`
}

type TMDbResponse struct {
	Results []TMDbMovie `json:"results"`
}

func fetchMoviesFromTMDb(apiKey string) ([]Movie, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/popular?api_key=%s&language=en-US&page=1", apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch movies: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	var tmdbResponse TMDbResponse
	if err := json.NewDecoder(resp.Body).Decode(&tmdbResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	var movies []Movie
	for _, tmdbMovie := range tmdbResponse.Results {
		year := ""
		if len(tmdbMovie.Year) >= 4 {
			year = tmdbMovie.Year[:4]
		}
		movies = append(movies, Movie{
			ID:    tmdbMovie.ID,
			Title: tmdbMovie.Title,
			Year:  year,
		})
	}

	return movies, nil
}

func moviesHandler(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		http.Error(w, "TMDB API key is missing", http.StatusInternalServerError)
		return
	}

	movies, err := fetchMoviesFromTMDb(apiKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("error fetching movies: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	http.HandleFunc("/api/movies", moviesHandler)

	log.Println("Server started on :", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
