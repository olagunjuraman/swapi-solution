package handlers

import (
	"busha/client"
	u "busha/utils"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
)

var SORT = map[string]string{
	"asc":  "ASC",
	"desc": "DESC",
}

type ByName []CharacterData
type ByGender []CharacterData
type ByHeight []CharacterData

func (a ByName) Len() int           { return len(a) }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (a ByHeight) Len() int           { return len(a) }
func (a ByHeight) Less(i, j int) bool { return a[i].Height < a[j].Height }
func (a ByHeight) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (a ByGender) Len() int           { return len(a) }
func (a ByGender) Less(i, j int) bool { return a[i].Gender < a[j].Gender }
func (a ByGender) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type Character struct {
	Name      string   `json:"name"`
	Height    string   `json:"height"`
	Mass      string   `json:"mass"`
	HairColor string   `json:"hair_color"`
	SkinColor string   `json:"skin_color"`
	EyeColor  string   `json:"eye_color"`
	BirthYear string   `json:"birth_year"`
	Gender    string   `json:"gender"`
	Homeworld string   `json:"homeworld"`
	Films     []string `json:"films"`
	URL       string   `json:"url"`
}

type CharacterData struct {
	Gender string `json:"gender"`
	Name   string `json:"name"`
	Height string `json:"height"`
}

type Metadata struct {
	Count             int    `json:"count"`
	TotalHeightMeters string `json:"total_height_meters"`
	TotalHeightFeet   string `json:"total_height_feet"`
}

type CharacterResponse struct {
	Metadata   `json:"metadata"`
	Characters []CharacterData `json:"characters"`
}

func (h *Handler) GetMovieCharacters(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movieID := params["id"]
	sortby := r.URL.Query().Get("sort_by")
	order_by := r.URL.Query().Get("order_by")
	filter_by := r.URL.Query().Get("filter_by")

	var f Film
	mId, _ := strconv.Atoi(movieID)
	var c []Character

	//check if characters for this movie exists on redis
	characterCacheKey := "characters:" + strconv.Itoa(mId)
	err := h.RedGet(characterCacheKey, &c)
	if err != nil {

		//check for movies on redis
		err = h.RedGet("movies", &f)
		if err != nil {
			f, _, err = h.callCacheMovies()
			if err != nil {
				response := u.Message(http.StatusBadRequest, "Something went wrong. Try again")
				u.Respond(w, http.StatusBadRequest, response)
				return
			}
		}
		for _, movie := range f.Results {
			id, err2 := u.ResourceId(movie.URL)
			if err2 != nil {
				log.Println("Something went wrong parsing id from URL ")
			}
			//If movieID is equals to ID passed
			if id == mId {
				c, err = h.getCharacters(movie.Characters, characterCacheKey)
				if err != nil {
					response := u.Message(http.StatusBadRequest, "Something went wrong fetching characters. Try again")
					u.Respond(w, http.StatusBadRequest, response)
					return
				}
				break
			}
		}
	}
	var characterData []CharacterData
	var totalheight = 0
	for _, char := range c {
		if filter_by != "" {
			if char.Gender == filter_by {
				characterData = append(characterData, CharacterData{
					Gender: char.Gender,
					Name:   char.Name,
					Height: char.Height,
				})
				totalheight += convertHeighttoInt(char.Height)
			}
		} else {
			characterData = append(characterData, CharacterData{
				Gender: char.Gender,
				Name:   char.Name,
				Height: char.Height,
			})
			totalheight += convertHeighttoInt(char.Height)
		}
	}
	count := len(characterData)
	response := u.Message(http.StatusOK, "Characters in this movie")
	metadata := Metadata{
		Count:             count,
		TotalHeightMeters: strconv.Itoa(totalheight),
		TotalHeightFeet:   convertHeighttoFeet(totalheight),
	}
	if sortby != "" {
		if order_by == "" {
			order_by = "ASC"
		}
		if SORT[order_by] == "" {
			order_by = "ASC"
		} else {
			order_by = SORT[order_by]
		}
		characterData = sortCharacter(characterData, sortby, order_by)
	}
	response["data"] = CharacterResponse{
		Metadata:   metadata,
		Characters: characterData,
	}
	u.Respond(w, http.StatusOK, response)
	return

}

func (h *Handler) getCharacters(urls []string, cacheKey string) (people []Character, err error) {
	for _, url := range urls {
		var c Character
		if err = client.Call(url, &c); err != nil {
			return nil, err
		}
		people = append(people, c)

		//store characters on redis
		if err = h.RedSet(cacheKey, people); err != nil {
			log.Println("Something went wrong storing characters cache on redis")
		}
	}
	return people, nil
}

func convertHeighttoInt(height string) int {
	h, err := strconv.Atoi(height)
	if err != nil {
		log.Println("Height could not be converted to int")
		return 0
	}
	return h
}

func convertHeighttoFeet(height int) string {
	if height != 0 {
		inches := float64(height) / 2.54
		feet := int(inches / 12)
		inches = inches - float64(12*feet)
		resp := fmt.Sprintf(" %vft and %v inches", feet, math.Round(inches*100)/100)
		return resp
	}
	return ""
}

func sortCharacter(array []CharacterData, key string, method string) []CharacterData {
	if method == "ASC" {
		if key == "name" {
			sort.Sort(ByName(array))
		}
		if key == "gender" {
			sort.Sort(ByGender(array))
		}
		if key == "height" {
			sort.Sort(ByHeight(array))
		}
		return array
	}
	if method == "DESC" {
		if key == "name" {
			sort.Sort(sort.Reverse(ByName(array)))
		}
		if key == "gender" {
			sort.Sort(sort.Reverse(ByGender(array)))
		}
		if key == "height" {
			sort.Sort(sort.Reverse(ByHeight(array)))
		}
		return array
	}
	return nil
}
