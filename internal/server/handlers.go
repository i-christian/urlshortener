package server

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"

	"github.com/redis/go-redis/v9"
)

// GetLongUrl handler method accepts a url with an identifier key.
// Retrieves long url from the database using the key pathvalue
// Then initiates a url redirect to the long url.
func (s *State) GetLongUrl(w http.ResponseWriter, r *http.Request) {
	urlKey := r.PathValue("identifier")

	resp, err := s.redisClient.HGet(r.Context(), fmt.Sprintf("urls:%s", urlKey), "long_url").Result()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	http.Redirect(w, r, resp, http.StatusFound)
}

// ShortenUrl handler method accepts a form with a long url.
// Creates a short url, and saves it into redis
func (s *State) ShortenUrl(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	longUrl := r.FormValue("long_url")

	encodedUrl, err := s.createShortCode(r, longUrl)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	hashFields := []string{
		"key", encodedUrl,
		"short_url", fmt.Sprintf("http://localhost/%s", encodedUrl),
		"long_url", longUrl,
	}

	_, err = s.redisClient.HSet(r.Context(), fmt.Sprintf("urls:%s", encodedUrl), hashFields).Result()
	if err != nil {
		slog.Error("internal server error", "error", err.Error())
		return
	} else {
		s.newKey = encodedUrl
	}

	w.WriteHeader(http.StatusCreated)
}

// Latest Method displays the latest created url result json
func (s *State) Latest(w http.ResponseWriter, r *http.Request) {
	resp, err := s.redisClient.HGetAll(r.Context(), fmt.Sprintf("urls:%s", s.newKey)).Result()
	if err != nil {
		return
	}

	respBody, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")

	w.Write(respBody)
}

func (s *State) createShortCode(r *http.Request, longUrl string) (string, error) {
	urlHash := md5.Sum([]byte(longUrl))
	var parsedBase62 string

	// handle hash collisions
	for idx := 5; idx < 8; idx++ {
		slicedHash := urlHash[:idx]

		// encode md5 hash to a base62 string
		var i big.Int
		i.SetBytes(slicedHash)
		parsedBase62 = i.Text(62)

		res, err := s.redisClient.HGet(r.Context(), fmt.Sprintf("urls:%s", parsedBase62), "long_url").Result()
		if err == redis.Nil {
			return parsedBase62, nil
		} else if res == longUrl {
			return "", errors.New("long url already exists")
		} else {
			return "", errors.New("unexpected error occured")
		}

	}

	return "", errors.New("failed to create a shortCode")
}
