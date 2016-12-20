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

package podcast

import (
	"net/http"

	"github.com/te-th/podca-api/feed"
	"github.com/te-th/podca-api/middleware"
	"golang.org/x/net/context"
)

// SearchHandler creates a http.HandlerFunc searching for Podcasts on the given PodcastSearch.
func SearchHandler(podcastSearcher *Search) http.HandlerFunc {
	return middleware.ServeHTTP(func(ctx context.Context, responseWriter http.ResponseWriter, request *http.Request) {
		s := &search{podcastSearcher: podcastSearcher, ctx: ctx}
		s.search(responseWriter, request)
	})
}

// SuggestionHandler is a http.HandlerFunc serving search results
func SuggestionHandler(podcastSuggestion SuggestionFacade) http.HandlerFunc {
	return middleware.ServeHTTP(func(ctx context.Context, responseWriter http.ResponseWriter, request *http.Request) {
		s := &suggestion{facade: podcastSuggestion, ctx: ctx}
		s.suggest(responseWriter, request)
	})
}

// Podcast struct is strong coupled to the Apple iTunes format.
type Podcast struct {
	ID             int64    `json:"id"`
	ArtistName     string   `json:"artistName"`
	CollectionName string   `json:"collectionName"`
	FeedURL        string   `json:"feedUrl"`
	CollectionID   int64    `json:"collectionId"`
	TrackID        int64    `json:"trackId"`
	Genres         []string `json:"genres"`
}

// SearchEngine searches for Podcasts with the given term.
type SearchEngine interface {
	Search(ctx context.Context, term string, limit string) ([]Podcast, error)
}

// SearchFacade is the facade that handle the usecase
// of searching podcast for a given term and storing the results into the datastore
type SearchFacade interface {
	Search(ctx context.Context, term string, limit string) ([]Podcast, error)
}

// Search encapsulates Task and SearchEngine capabilities
type Search struct {
	FeedTask     feed.TaskFacade
	SearchEngine SearchEngine
}

// SuggestionFacade exposes a search result
type SuggestionFacade interface {
	Suggestion(ctx context.Context, term string, limit string) ([]Podcast, error)
}
