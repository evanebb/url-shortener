package main

type InMemoryShortenedURLRepository struct {
	urls []ShortenedURL
}

func newInMemoryShortenedURLRepository() *InMemoryShortenedURLRepository {
	return &InMemoryShortenedURLRepository{urls: make([]ShortenedURL, 0)}
}

func (r *InMemoryShortenedURLRepository) IsIDAlreadyInUse(id string) bool {
	for _, u := range r.urls {
		if u.ID == id {
			return true
		}
	}
	return false
}

func (r *InMemoryShortenedURLRepository) StoreShortenedURL(u ShortenedURL) error {
	if r.IsIDAlreadyInUse(u.ID) {
		return ErrURLAlreadyExists
	}

	r.urls = append(r.urls, u)
	return nil
}

func (r *InMemoryShortenedURLRepository) GetShortenedURLByID(id string) (ShortenedURL, error) {
	for _, u := range r.urls {
		if u.ID == id {
			return u, nil
		}
	}

	var emptyUrl ShortenedURL
	return emptyUrl, ErrNoURLFound
}
