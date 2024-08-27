package handlers

import (
	u "busha/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type CommentRequest struct {
	FilmID  int    `json:"film_id"`
	Comment string `json:"comment"`
}

func (h *Handler) AddComment(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var c CommentRequest
	err := decoder.Decode(&c)
	if err != nil {
		response := u.Message(http.StatusBadRequest, "Something went wrong with input passed")
		u.Respond(w, http.StatusBadRequest, response)
		return
	}
	//500chars limit
	if len(c.Comment) > 500 {
		response := u.Message(http.StatusBadRequest, "Comment character(s) should not exceed 500 chars.")
		u.Respond(w, http.StatusBadRequest, response)
		return
	}

	//Check that filmID exists on our record
	var n map[int]int
	err = h.RedGet("movieID:list", &n)
	if err != nil {
		_, n, err = h.callCacheMovies()
		if err != nil {
			response := u.Message(http.StatusBadRequest, "Something went wrong with fetching movies to map comment. Try again")
			u.Respond(w, http.StatusBadRequest, response)
			return
		}
	}

	_, exists := n[c.FilmID]
	if !exists {
		response := u.Message(http.StatusBadRequest, "This Film ID is Invalid. Kindly check ID and try again")
		u.Respond(w, http.StatusBadRequest, response)
		return
	}

	_, err = h.Repo.CreateComment(c.Comment, c.FilmID, u.GetIp(r))
	if err != nil {
		response := u.Message(http.StatusBadRequest, "Something went wrong with creating comment")
		u.Respond(w, http.StatusBadRequest, response)
		return
	}
	response := u.Message(http.StatusOK, "Comment was successfully created")
	u.Respond(w, http.StatusOK, response)
	return
}

func (h *Handler) GetMovieComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movieID := params["id"]
	comments, err := h.Repo.GetComment(movieID)
	if err != nil {
		response := u.Message(http.StatusBadRequest, "Something went wrong fetching comments for this movie. Try again.")
		u.Respond(w, http.StatusBadRequest, response)
		return
	}
	response := u.Message(http.StatusOK, "Comments")
	response["data"] = comments
	u.Respond(w, http.StatusOK, response)
	return
}
