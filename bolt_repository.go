package main

import (
	bolt "go.etcd.io/bbolt"
	"net/url"
)

type BoltShortenedURLRepository struct {
	db *bolt.DB
}

func newBoltShortenedURLRepository(db *bolt.DB) (*BoltShortenedURLRepository, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		// We do not need the bucket here, we just want to create it
		_, err := tx.CreateBucketIfNotExists([]byte("urls"))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &BoltShortenedURLRepository{db: db}, nil
}

func (r *BoltShortenedURLRepository) IsIDAlreadyInUse(id string) bool {
	exists := false

	// We ignore the error here, since we do not return any error from the function passed to View()
	_ = r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("urls"))
		v := b.Get([]byte(id))
		if v != nil {
			exists = true
		}
		return nil
	})

	return exists
}

func (r *BoltShortenedURLRepository) StoreShortenedURL(u ShortenedURL) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("urls"))
		return b.Put([]byte(u.ID), []byte(u.FullURL.String()))
	})
}

func (r *BoltShortenedURLRepository) GetShortenedURLByID(id string) (ShortenedURL, error) {
	var rawUrl string
	var err error

	err = r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("urls"))
		v := b.Get([]byte(id))
		if v == nil {
			return ErrNoURLFound
		}
		rawUrl = string(v)
		return nil
	})

	var u ShortenedURL
	if err != nil {
		return u, err
	}

	parsedUrl, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		return u, err
	}

	return ShortenedURL{ID: id, FullURL: parsedUrl}, nil
}
