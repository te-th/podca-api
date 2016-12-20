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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"google.golang.org/appengine/delay"

	"github.com/te-th/podca-api/feed"
	"github.com/te-th/podca-api/middleware"
	"github.com/te-th/podca-api/utility"
	"golang.org/x/net/context"
)

type search struct {
	podcastSearcher *Search
	ctx             context.Context
}

func (s *search) search(responseWriter http.ResponseWriter, request *http.Request) {
	var term = request.FormValue("term")
	if term == "" {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte("parameter term is missing e.g. term=WDR "))
	} else {

		var limit = request.FormValue("limit")
		podcasts, err := s.podcastSearcher.SearchEngine.Search(s.ctx, term, limit)
		if err != nil {
			responseWriter.WriteHeader(http.StatusInternalServerError)
			panic(err)
		}

		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if jsonerr := json.NewEncoder(responseWriter).Encode(podcasts); jsonerr != nil {
			responseWriter.WriteHeader(http.StatusInternalServerError)
			panic(jsonerr)
		}
	}
}

// NewSearchEngine creates a new ITunesSearchEngine instance.
func NewSearchEngine(httpClient middleware.HTTPClientFacade) SearchEngine {
	return &iTunesSearchEngine{
		HTTPClient: httpClient,
	}
}

type iTunesSearchEngine struct {
	HTTPClient middleware.HTTPClientFacade
}

// Search finds Podcasts with the given term.
func (searchEngine *iTunesSearchEngine) Search(ctx context.Context,
	term string, limit string) ([]Podcast, error) {

	var urlString = "https://itunes.apple.com/search?limit=" +
		url.QueryEscape(utility.CheckLimit(limit)) +
		"&country=DE&entity=podcast&term=" +
		url.QueryEscape(term)

	res, err := searchEngine.HTTPClient.Get(ctx, urlString)
	if err != nil {
		return nil, err
	}

	jsonResponse, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	type SearchResponse struct {
		ResultCount int
		Results     []Podcast
	}

	var result SearchResponse
	if err := json.Unmarshal(jsonResponse, &result); err != nil {
		return nil, err
	}

	return result.Results, nil
}

// NewPodcastSearch creates a new PodcastSearchFacade instance
func NewPodcastSearch(feedTask feed.TaskFacade, searchEngine SearchEngine) *Search {
	return &Search{
		FeedTask:     feedTask,
		SearchEngine: searchEngine,
	}
}

// Search searches for and persists Podcasts
func (podcastSearch *Search) Search(ctx context.Context, term string, limit string) ([]Podcast, error) {
	podcasts, err := podcastSearch.SearchEngine.Search(ctx, term, limit)
	if err != nil {
		return nil, err
	}

	for _, podcast := range podcasts {
		var delayedTask = delay.Func("feedTask", func(ctx context.Context, podcast Podcast) {
			podcastSearch.FeedTask.FetchAndStore(ctx, podcast.FeedURL, podcast.CollectionID)
		})

		delayedTask.Call(ctx, podcast)
	}

	return podcasts, nil
}
