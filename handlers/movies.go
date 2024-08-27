package handlers

import (
	"busha/client"
	u "busha/utils"
	"bytes"
	"encoding/gob"
	"fmt"
	"net/http"
	"sort"
)

type Film struct {
	Count    int         `json:"count"`
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []Results
}

type Results struct {
	Title        string   `json:"title"`
	EpisodeID    int64    `json:"episode_id"`
	OpeningCrawl string   `json:"opening_crawl"`
	Characters   []string `json:"characters"`
	Created      string   `json:"created"`
	ReleaseDate  string   `json:"release_date"`
	Edited       string   `json:"edited"`
	URL          string   `json:"url"`
}

type FilmResponse struct {
	Name         string `json:"name"`
	OpeningCrawl string `json:"opening_crawl"`
	TotalComment int    `json:"total_comment"`
	ReleaseDate  string `json:"release_date"`
}

func (h *Handler) GetMovies(w http.ResponseWriter, r *http.Request) {
	var f Film
	//check if movies exists on cache
	err := h.RedGet("movies", &f)
	if err != nil {
		f, _, err = h.callCacheMovies()
		if err != nil {
			response := u.Message(http.StatusBadRequest, "Something went wrong fetching movies. Try again")
			u.Respond(w, http.StatusBadRequest, response)
			return
		}
	}

	comments, err := h.Repo.GetTotalComments()
	if err != nil {
		response := u.Message(http.StatusBadRequest, "Something went wrong parsing comments")
		u.Respond(w, http.StatusBadRequest, response)
		return
	}

	//Sort movies from earliest to newest
	sortMovies := f.Results
	sort.Slice(sortMovies, func(i, j int) bool {
		return sortMovies[i].ReleaseDate < sortMovies[j].ReleaseDate
	})

	var filmresponse []FilmResponse

	for _, movie := range sortMovies {
		movieID, err2 := u.ResourceId(movie.URL)
		if err2 != nil {
			fmt.Println("Something went wrong parsing ID from URL")
		}
		filmresponse = append(filmresponse, FilmResponse{
			Name:         movie.Title,
			OpeningCrawl: movie.OpeningCrawl,
			TotalComment: comments[movieID],
			ReleaseDate:  movie.ReleaseDate,
		})
	}
	response := u.Message(http.StatusOK, "Movies")
	response["data"] = filmresponse
	u.Respond(w, http.StatusOK, response)
	return
}

func (h *Handler) callCacheMovies() (Film, map[int]int, error) {
	var f Film

	// pull from swapi
	if err := client.Call("/films", &f); err != nil {
		return Film{}, nil, err
	}
	//Store ID's of all movies in am := make([]int, len(f.Results)) map
	m := make(map[int]int)
	for _, movie := range f.Results {
		movieID, err := u.ResourceId(movie.URL)
		if err != nil {
			fmt.Println("Something went wrong parsing ID from URL")
		}
		m[movieID] = movieID
	}

	//store movies on redis
	if err := h.RedSet("movies", f); err != nil {
		fmt.Println("Something went wrong storing cache on redis")
	}

	//store movieID map on redis
	if err := h.RedSet("movieID:list", m); err != nil {
		fmt.Println("Something went wrong storing cache on redis")
	}
	return f, m, nil
}

func (h Handler) RedSet(field string, value interface{}) error {
	var buf bytes.Buffer
	var encoder = gob.NewEncoder(&buf)

	id := "busha"
	err := encoder.Encode(value)
	if err != nil {
		return err
	}
	return h.RedisDB.HSet(id, field, buf.Bytes()).Err()
}

func (h Handler) RedGet(field string, value interface{}) error {
	id := "busha"
	buf, err := h.RedisDB.HGet(id, field).Bytes()
	if err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewReader(buf)).Decode(value)
}
