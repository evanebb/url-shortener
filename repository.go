package main

import (
	"errors"
)

var ErrURLAlreadyExists = errors.New("an URL with this ID already exists")
var ErrNoURLFound = errors.New("no URL found with the requested ID")

type ShortenedURLRepository interface {
	// IsIDAlreadyInUse checks whether the given ID is already in use for an existing URL.
	IsIDAlreadyInUse(id string) bool
	// StoreShortenedURL will persist the given ShortenedURL, so it can be retrieved later on.
	//
	// ErrURLAlreadyExists may be returned if the given ID is already in use for a URL, but this can vary per implementation.
	// Use IsIDAlreadyInUse if you need to be sure that the ID is not in use yet before calling this function.
	StoreShortenedURL(ShortenedURL) error
	// GetShortenedURLByID will look up the ShortenedURL for the given ID, and return it.
	//
	// ErrNoURLFound is returned if no URL can be found for the given ID.
	GetShortenedURLByID(id string) (ShortenedURL, error)
}
