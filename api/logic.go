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
	"github.com/te-th/podca-api/utility"
	"golang.org/x/net/context"
	"google.golang.org/appengine/delay"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

// PodcastSearch is the facade that handle the usecase
// of searching podcast for a given term and storing the results into the datastore
type PodcastSearch interface {
	Search(ctx context.Context, term string, limit string) ([]domain.Podcast, error)
}

// PodcastSearcher combines FeedTasks and a SearchEngine
type PodcastSearcher struct {
	FeedTask     feedTask
	SearchEngine PodcastSearchEngine
}

// NewPodcastSearcher creates a new PodcastSearcher instance
func NewPodcastSearcher(feedTask feedTask, searchEngine PodcastSearchEngine) *PodcastSearcher {
	return &PodcastSearcher{
		FeedTask:     feedTask,
		SearchEngine: searchEngine,
	}
}

// Search looks up Podcasts with the given search term
func (podcastSearcher *PodcastSearcher) Search(ctx context.Context,
	term string, limit string) ([]domain.Podcast, error) {
	podcasts, err := podcastSearcher.SearchEngine.Search(ctx, term, limit)
	if err != nil {
		return nil, err
	}

	for _, podcast := range podcasts {
		var delayedTask = delay.Func("feedWorker", func(ctx context.Context, podcast domain.Podcast) {
			podcastSearcher.FeedTask.FetchAndStore(ctx, podcast)
		})

		delayedTask.Call(ctx, podcast)
	}

	return podcasts, nil

}

type feedTask interface {
	FetchAndStore(ctx context.Context, podcast domain.Podcast)
}

// FeedTaskWorker can post-process search results
type FeedTaskWorker struct {
	FeedRepository domain.FeedRepository
}

// NewFeedTaskWorker creates a new FeedTaskWorker
func NewFeedTaskWorker(feedRepo domain.FeedRepository) *FeedTaskWorker {
	return &FeedTaskWorker{
		FeedRepository: feedRepo,
	}
}

// FetchData fetches data from the given url
func (worker *FeedTaskWorker) fetchData(ctx context.Context, url string) ([]byte, error) {

	client := urlfetch.Client(ctx)
	res, err := client.Get(url)
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

// FetchAndStore fetches the Feed assigned to given Podcast and persists it,
func (worker *FeedTaskWorker) FetchAndStore(ctx context.Context, podcast domain.Podcast) {

	data, fetcherr := worker.fetchData(ctx, podcast.FeedURL)
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
	}

	var feed = &rss.Feed

	feed.ID = podcast.CollectionID

	worker.FeedRepository.Save(ctx, feed)
}

// PodcastSearchEngine searches for Podcasts with the given term.
type PodcastSearchEngine interface {
	Search(ctx context.Context, term string, limit string) ([]domain.Podcast, error)
}

// NewSearchEngine creates a new ITunesSearchEngine instance.
func NewSearchEngine() PodcastSearchEngine {
	return &iTunesSearchEngine{}
}

type iTunesSearchEngine struct {
}

// Search finds Podcasts with the given term.
func (searchEngine *iTunesSearchEngine) Search(ctx context.Context,
	term string, limit string) ([]domain.Podcast, error) {
	client := urlfetch.Client(ctx)

	var urlString = "https://itunes.apple.com/search?limit=" +
		url.QueryEscape(utility.CheckLimit(limit)) +
		"&country=DE&entity=podcast&term=" +
		url.QueryEscape(term)

	res, err := client.Get(urlString)
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
