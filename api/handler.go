/*
Copyright [2016] TE,TH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/te-th/podca-api/domain"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// RootHandler is a http.HandlerFunc serving request against "/".
func RootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API Podca"))
	}
}

// PodcastSuggestionHandler is a http.HandlerFunc serving search results
func PodcastSuggestionHandler(podcastSuggestion PodcastSuggestionFacade) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		var term = r.FormValue("term")
		var limit = r.FormValue("limit")
		podcasts, err := podcastSuggestion.Suggestion(ctx, term, limit)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			panic(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if jsonerr := json.NewEncoder(w).Encode(podcasts); jsonerr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			panic(jsonerr)
		}
	}
}

// PodcastSearchHandler is a http.HandlerFunc looking up Podcast through the given searcher and exposes the result.
func PodcastSearchHandler(podcastSearcher *PodcastSearch) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		var term = r.FormValue("term")

		if term == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("parameter term is missing e.g. term=WDR "))
		} else {

			var limit = r.FormValue("limit")
			podcasts, err := podcastSearcher.SearchEngine.Search(ctx, term, limit)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				panic(err)
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if jsonerr := json.NewEncoder(w).Encode(podcasts); jsonerr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				panic(jsonerr)
			}
		}
	}
}

// FeedHandler is a http.HandlerFunc serving a Feed with a given ID.
func FeedHandler(feedRepo domain.FeedRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "GET":
			ctx := appengine.NewContext(r)
			vars := mux.Vars(r)

			feedID, _ := strconv.Atoi(vars["feedId"])
			if feedID > 0 {
				log.Infof(ctx, "FEEDID> %d", int64(feedID))
				result, _ := feedRepo.Get(ctx, int64(feedID))
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				if jsonerr := json.NewEncoder(w).Encode(result); jsonerr != nil {
					panic(jsonerr)
				}
			} else {
				result, _ := feedRepo.GetAll(ctx)
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				if jsonerr := json.NewEncoder(w).Encode(result); jsonerr != nil {
					panic(jsonerr)
				}
			}

		case "PUT":
			// update
			w.Write([]byte("NOT YET IMPLEMENTED "))
			w.WriteHeader(http.StatusNoContent)
		case "POST":
			// create
			w.Write([]byte("NOT YET IMPLEMENTED "))
			w.WriteHeader(http.StatusNoContent)
		case "DELETE":
			// remove
			w.Write([]byte("NOT YET IMPLEMENTED "))
			w.WriteHeader(http.StatusNoContent)
		}
	}
}

// FeedImageHandler is a http.HandlerFunc serving FeedImages.
func FeedImageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API Podca "))
	}
}

// FeedEpisodeHandler is a http.HandlerFunc serving FeedEpisodes.
func FeedEpisodeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API Podca "))
	}
}
