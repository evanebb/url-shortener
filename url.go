package main

import (
	"math/rand"
	"net/url"
)

type ShortenedURL struct {
	ID      string
	FullURL *url.URL
}

// GenerateShortenedURL will generate a new ShortenedURL for the given URL, store it and return it.
func GenerateShortenedURL(r ShortenedURLRepository, u *url.URL) (ShortenedURL, error) {
	var id string

	for {
		id = generateID()
		if !r.IsIDAlreadyInUse(id) {
			break
		}
	}

	s := ShortenedURL{ID: id, FullURL: u}
	err := r.StoreShortenedURL(s)
	if err != nil {
		var emptyUrl ShortenedURL
		return emptyUrl, err
	}

	return s, nil
}

// generateID will generate a small, unique base62 ID.
func generateID() string {
	// Base 62, should be enough in terms of uniqueness
	set := []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

	id := make([]rune, 5)
	for i := range id {
		id[i] = set[rand.Intn(len(set))]
	}

	return string(id)
}
