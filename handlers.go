package main

import (
	"net/http"
	"net/url"
)

func handleGenerateURL(repo ShortenedURLRepository, t Templater, baseUrl *url.URL) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		err = r.ParseForm()
		if err != nil {
			w.WriteHeader(500)
			t.renderError(w)
			return
		}

		parsedUrl, err := url.ParseRequestURI(r.PostFormValue("url"))
		if err != nil {
			w.WriteHeader(500)
			t.renderError(w)
			return
		}

		su, err := GenerateShortenedURL(repo, parsedUrl)
		if err != nil {
			w.WriteHeader(500)
			t.renderError(w)
			return
		}

		var responseURL string
		if baseUrl != nil {
			responseURL = baseUrl.String() + "/u/" + su.ID
		} else {
			scheme := "https://"
			if r.TLS == nil {
				// If the request was made without TLS (so no HTTPS), set the scheme to http://
				scheme = "http://"
			}
			responseURL = scheme + r.Host + "/u/" + su.ID
		}

		w.WriteHeader(200)
		t.render(w, "generated", responseURL)
	})
}

func handleRedirectURL(repo ShortenedURLRepository, t Templater) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		su, err := repo.GetShortenedURLByID(id)
		if err != nil {
			w.WriteHeader(200)
			t.render(w, "unknown", nil)
			return
		}

		http.Redirect(w, r, su.FullURL.String(), http.StatusFound)
	})
}

func handleIndex(t Templater) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(404)
			t.render(w, "errors/404", nil)
			return
		}

		w.WriteHeader(200)
		t.render(w, "index", nil)
	})
}
