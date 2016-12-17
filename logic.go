package podca_api

import (
	"net/url"
	"io/ioutil"
	"encoding/json"
	"golang.org/x/net/context"
	"encoding/xml"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/delay"
)

type FeedCore interface {

}



// PodcastSearch is the facade that handle the usecase
// of searching podcast for a given term and storing the results into the datastore
type PodcastSearch interface {
	Search(ctx context.Context, term string) ([]Podcast, error)
}

type PodcastSearcher struct {
	FeedTask FeedTask
	SearchEngine PodcastSearchEngine
}

func NewPodcastSearcher(feedTask FeedTask, searchEngine PodcastSearchEngine) *PodcastSearcher {
	return &PodcastSearcher{
		FeedTask: feedTask,
		SearchEngine: searchEngine,
	}
}

func (podcastSearcher *PodcastSearcher) Search(ctx context.Context, term string) ([]Podcast, error) {
	podcasts, err := podcastSearcher.SearchEngine.Search(ctx, term) ; if err != nil {
		return nil, err
	}

	for _, podcast := range podcasts {
		var delayedTask = delay.Func("feedWorker", func(ctx context.Context, podcast Podcast){
			podcastSearcher.FeedTask.FetchAndStore(ctx, podcast)
		})

		delayedTask.Call(ctx, podcast)
	}

	return podcasts, nil

}

type FeedTask interface {
	FetchAndStore(ctx context.Context,  podcast Podcast)

}

type FeedTaskWorker struct {
	FeedRepository FeedRepository
}

func NewFeedTaskWorker(feedRepo FeedRepository) *FeedTaskWorker {
	return &FeedTaskWorker{
		FeedRepository: feedRepo,
	}
}

func (worker *FeedTaskWorker) FetchData(ctx context.Context, url string) ([]byte, error) {

	client := urlfetch.Client(ctx)
	res, err := client.Get(url); if err != nil {
		return  nil, err
	}

	log.Infof(ctx, "FETCHED URL: %s  WITH HTTP STATUSCODE %d ", url, res.StatusCode)

	xmlResponse, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return  nil, err
	}
	return xmlResponse, nil
}

func (worker *FeedTaskWorker) FetchAndStore(ctx context.Context, podcast Podcast) {

	data, fetcherr := worker.FetchData(ctx, podcast.FeedUrl); if fetcherr != nil {
		panic(fetcherr)
	}

	type Rss struct {
		Feed	Feed `xml:"channel"`
	}

	var rss Rss

	xmlerr := xml.Unmarshal(data, &rss); if xmlerr != nil {
		panic(xmlerr)
		return
	}

	var feed = &rss.Feed

	// Use CollectionId as Feed.id
	feed.Id = podcast.CollectionId

	worker.FeedRepository.save(ctx,feed)
}

type PodcastSearchEngine interface {
	Search(ctx context.Context, term string) ([]Podcast, error)
}

func NewSearchEngine() *ITunesSearchEngine {
	return &ITunesSearchEngine{}
}

type ITunesSearchEngine struct {

}

func (searchEngine *ITunesSearchEngine) Search(ctx context.Context, term string) ([]Podcast, error) {
	client := urlfetch.Client(ctx)

	var urlString = "https://itunes.apple.com/search?term="+ url.QueryEscape(term) + "&country=DE&entity=podcast"

	res, err := client.Get(urlString); if err != nil {
		return  nil, err
	}

	jsonResponse, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return  nil, err
	}

	type SearchResponse struct {
		ResultCount int
		Results []Podcast
	}

	var result SearchResponse
	if err := json.Unmarshal(jsonResponse, &result); err != nil {
		return  nil, err
	}

	var podcasts []Podcast = result.Results

	return podcasts, nil
}