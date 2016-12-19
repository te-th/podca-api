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
	"encoding/xml"
	"io/ioutil"
	"net/url"

	"github.com/te-th/podca-api/domain"
	"github.com/te-th/podca-api/networking"
	"github.com/te-th/podca-api/utility"
	"golang.org/x/net/context"
	"google.golang.org/appengine/delay"
	"google.golang.org/appengine/log"
)

// PodcastSuggestionFacade exposes a search result
type PodcastSuggestionFacade interface {
	Suggestion(ctx context.Context, term string, limit string) ([]domain.Podcast, error)
}

type podcastSuggestion struct {
	SearchEngine PodcastSearchEngine
}

// NewPodcastSuggestion creates a new PodcastSuggestionFacade instance
func NewPodcastSuggestion(searchEngine PodcastSearchEngine) PodcastSuggestionFacade {
	return &podcastSuggestion{
		SearchEngine: searchEngine,
	}
}

func (suggestion *podcastSuggestion) Suggestion(ctx context.Context,
	term string, limit string) ([]domain.Podcast, error) {

	podcasts, err := suggestion.SearchEngine.Search(ctx, term, limit)
	if err != nil {
		return nil, err
	}
	return podcasts, nil
}

/*
type CacheFacade interface {
	GetOrSet(ctx context.Context, key string, podcasts []domain.Podcast) ([]domain.Podcast, error)
}

type cache struct {
}

func (cache *cache) GetOrSet(ctx context.Context, key string, podcasts []domain.Podcast) ([]domain.Podcast, error) {

	return nil, nil
}
*/

// PodcastSearchFacade is the facade that handle the usecase
// of searching podcast for a given term and storing the results into the datastore
type PodcastSearchFacade interface {
	Search(ctx context.Context, term string, limit string) ([]domain.Podcast, error)
}

// PodcastSearch encapsulates FeedTask and SearchEngine capabilities
type PodcastSearch struct {
	FeedTask     FeedTaskFacade
	SearchEngine PodcastSearchEngine
}

// NewPodcastSearch creates a new PodcastSearchFacade instance
func NewPodcastSearch(feedTask FeedTaskFacade, searchEngine PodcastSearchEngine) *PodcastSearch {
	return &PodcastSearch{
		FeedTask:     feedTask,
		SearchEngine: searchEngine,
	}
}

// Search searches for and persists Podcasts
func (podcastSearch *PodcastSearch) Search(ctx context.Context, term string, limit string) ([]domain.Podcast, error) {
	podcasts, err := podcastSearch.SearchEngine.Search(ctx, term, limit)
	if err != nil {
		return nil, err
	}

	for _, podcast := range podcasts {
		var delayedTask = delay.Func("feedTask", func(ctx context.Context, podcast domain.Podcast) {
			podcastSearch.FeedTask.FetchAndStore(ctx, podcast)
		})

		delayedTask.Call(ctx, podcast)
	}

	return podcasts, nil

}

// FeedTaskFacade combines fetching and storing of Podcasts
type FeedTaskFacade interface {
	FetchAndStore(ctx context.Context, podcast domain.Podcast)
}

// FeedTask fetches and stores Podcast Feeds
type FeedTask struct {
	FeedRepository domain.FeedRepository
	HTTPClient     networking.HTTPClientFacade
}

// NewFeedTask creates a new FeedTask
func NewFeedTask(feedRepo domain.FeedRepository, httpClient networking.HTTPClientFacade) *FeedTask {
	return &FeedTask{
		FeedRepository: feedRepo,
		HTTPClient:     httpClient,
	}
}

// FetchData fetches data from the given url.
func (task *FeedTask) FetchData(ctx context.Context, url string) ([]byte, error) {

	res, err := task.HTTPClient.Get(ctx, url)
	if err != nil {
		return nil, err
	}

	log.Infof(ctx, "FETCHED URL: %s  WITH HTTP STATUSCODE %d ", url, res.StatusCode)

	xmlResponse, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	return xmlResponse, nil
}

// FetchAndStore fetches the Feed assigned to given Podcast and persists it.
func (task *FeedTask) FetchAndStore(ctx context.Context, podcast domain.Podcast) {

	data, fetcherr := task.FetchData(ctx, podcast.FeedURL)
	if fetcherr != nil {
		panic(fetcherr)
	}

	type Rss struct {
		Feed domain.Feed `xml:"channel"`
	}

	var rss Rss

	xmlerr := xml.Unmarshal(data, &rss)
	if xmlerr != nil {
		panic(xmlerr)
		return
	}

	var feed = &rss.Feed

	feed.ID = podcast.CollectionID

	task.FeedRepository.Save(ctx, feed)
}

// PodcastSearchEngine searches for Podcasts with the given term.
type PodcastSearchEngine interface {
	Search(ctx context.Context, term string, limit string) ([]domain.Podcast, error)
}

// NewSearchEngine creates a new ITunesSearchEngine instance.
func NewSearchEngine(httpClient networking.HTTPClientFacade) PodcastSearchEngine {
	return &iTunesSearchEngine{
		HTTPClient: httpClient,
	}
}

type iTunesSearchEngine struct {
	HTTPClient networking.HTTPClientFacade
}

// Search finds Podcasts with the given term.
func (searchEngine *iTunesSearchEngine) Search(ctx context.Context,
	term string, limit string) ([]domain.Podcast, error) {

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
		Results     []domain.Podcast
	}

	var result SearchResponse
	if err := json.Unmarshal(jsonResponse, &result); err != nil {
		return nil, err
	}

	return result.Results, nil
}
