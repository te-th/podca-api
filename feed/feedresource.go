/*
Copyright [2016] TE,TH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://wwresponseWriter.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package feed

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strconv"

	"google.golang.org/appengine/log"

	"github.com/gorilla/mux"
	"github.com/te-th/podca-api/middleware"
	"golang.org/x/net/context"
)

type feedresource struct {
	ctx        context.Context
	repository Repository
}

func (f *feedresource) resource(responseWriter http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		vars := mux.Vars(request)
		feedID, _ := strconv.Atoi(vars["feedId"])
		if feedID > 0 {
			log.Infof(f.ctx, "FEEDID> %d", int64(feedID))
			result, _ := f.repository.Get(f.ctx, int64(feedID))
			responseWriter.WriteHeader(http.StatusOK)
			responseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if jsonerr := json.NewEncoder(responseWriter).Encode(result); jsonerr != nil {
				panic(jsonerr)
			}
		} else {
			result, _ := f.repository.GetAll(f.ctx)
			responseWriter.WriteHeader(http.StatusOK)
			responseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if jsonerr := json.NewEncoder(responseWriter).Encode(result); jsonerr != nil {
				panic(jsonerr)
			}
		}

	case "PUT":
		// update
		responseWriter.Write([]byte("NOT YET IMPLEMENTED "))
		responseWriter.WriteHeader(http.StatusNoContent)
	case "POST":
		// create
		responseWriter.Write([]byte("NOT YET IMPLEMENTED "))
		responseWriter.WriteHeader(http.StatusNoContent)
	case "DELETE":
		// remove
		responseWriter.Write([]byte("NOT YET IMPLEMENTED "))
		responseWriter.WriteHeader(http.StatusNoContent)
	}
}

// TaskFacade combines fetching and storing of Podcasts
type TaskFacade interface {
	FetchAndStore(ctx context.Context, feedURL string, collectionID int64)
}

// Task fetches and stores Podcast Feeds
type Task struct {
	FeedRepository Repository
	HTTPClient     middleware.HTTPClientFacade
}

// NewFeedTask creates a new FeedTask
func NewFeedTask(feedRepo Repository, httpClient middleware.HTTPClientFacade) *Task {
	return &Task{
		FeedRepository: feedRepo,
		HTTPClient:     httpClient,
	}
}

// FetchData fetches data from the given url.
func (task *Task) FetchData(ctx context.Context, url string) ([]byte, error) {

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
func (task *Task) FetchAndStore(ctx context.Context, feedURL string, collectionID int64) {

	data, fetcherr := task.FetchData(ctx, feedURL)
	if fetcherr != nil {
		panic(fetcherr)
	}

	type Rss struct {
		Feed `xml:"channel"`
	}

	var rss Rss

	xmlerr := xml.Unmarshal(data, &rss)
	if xmlerr != nil {
		panic(xmlerr)
	}

	var feed = &rss.Feed

	feed.ID = collectionID

	task.FeedRepository.Save(ctx, feed)
}
