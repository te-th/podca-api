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
	"golang.org/x/net/context"
	"github.com/te-th/podca-api/domain"
	"google.golang.org/appengine/delay"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

type FeedCore interface {
}

// PodcastSearch is the facade that handle the usecase
// of searching podcast for a given term and storing the results into the datastore
type PodcastSearch interface {
	Search(ctx context.Context, term string) ([]domain.Podcast, error)
}

type PodcastSearcher struct {
	FeedTask     FeedTask
	SearchEngine PodcastSearchEngine
}

func NewPodcastSearcher(feedTask FeedTask, searchEngine PodcastSearchEngine) *PodcastSearcher {
	return &PodcastSearcher{
		FeedTask:     feedTask,
		SearchEngine: searchEngine,
	}
}

func (podcastSearcher *PodcastSearcher) Search(ctx context.Context, term string) ([]domain.Podcast, error) {
	podcasts, err := podcastSearcher.SearchEngine.Search(ctx, term)
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

type FeedTask interface {
	FetchAndStore(ctx context.Context, podcast domain.Podcast)
}

type FeedTaskWorker struct {
	FeedRepository domain.FeedRepository
}

func NewFeedTaskWorker(feedRepo domain.FeedRepository) *FeedTaskWorker {
	return &FeedTaskWorker{
		FeedRepository: feedRepo,
	}
}

func (worker *FeedTaskWorker) FetchData(ctx context.Context, url string) ([]byte, error) {

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

func (worker *FeedTaskWorker) FetchAndStore(ctx context.Context, podcast domain.Podcast) {

	data, fetcherr := worker.FetchData(ctx, podcast.FeedURL)
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

	// Use CollectionId as Feed.id
	feed.ID = podcast.CollectionID

	worker.FeedRepository.Save(ctx, feed)
}

type PodcastSearchEngine interface {
	Search(ctx context.Context, term string) ([]domain.Podcast, error)
}

func NewSearchEngine() *ITunesSearchEngine {
	return &ITunesSearchEngine{}
}

type ITunesSearchEngine struct {
}

func (searchEngine *ITunesSearchEngine) Search(ctx context.Context, term string) ([]domain.Podcast, error) {
	client := urlfetch.Client(ctx)

	var urlString = "https://itunes.apple.com/search?term=" + url.QueryEscape(term) + "&country=DE&entity=podcast"

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
