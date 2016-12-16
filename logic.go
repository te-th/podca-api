package podca_api

import (
	"net/url"
	"io/ioutil"
	"encoding/json"
	"golang.org/x/net/context"
	"github.com/golang/appengine/urlfetch"
	"encoding/xml"
	"github.com/golang/appengine/log"
)

//type FeedWorker interface {
//	Retrieve(ctx context.Context, feedUrl string)
//
//}
type FeedWorker struct {
	FeedRepo *FeedRepo
}

func NewFeedWorker(feedRepo *FeedRepo) *FeedWorker {
	return &FeedWorker{
		FeedRepo: feedRepo,
	}
}

func (worker *FeedWorker) FetchData(ctx context.Context, url string) ([]byte, error) {

	client := urlfetch.Client(ctx)
	res, err := client.Get(url)

	if err != nil {
		return  nil, err
	}

	log.Infof(ctx, "HTTP STATUS> %d ", res.StatusCode)
	log.Infof(ctx, "HTTP STATUS> %s ", res.Header)

	xmlResponse, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return  nil, err
	}
	return xmlResponse, nil
}

func (worker *FeedWorker)  Retrieve(ctx context.Context, podcast Podcast) {

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

	worker.FeedRepo.save(ctx,feed)
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

	res, err := client.Get(urlString)

	if err != nil {
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