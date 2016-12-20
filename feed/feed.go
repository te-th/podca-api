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
	"net/http"

	"github.com/te-th/podca-api/middleware"
	"golang.org/x/net/context"
)

// Handler is a http.HandlerFunc serving a Feed with a given ID.
func Handler(feedRepository Repository) http.HandlerFunc {
	return middleware.ServeHTTP(func(ctx context.Context, responseWriter http.ResponseWriter, request *http.Request) {
		feed := &feedresource{ctx: ctx, repository: feedRepository}
		feed.resource(responseWriter, request)
	})
}

// ImageHandler is a http.HandlerFunc serving FeedImages.
func ImageHandler() http.HandlerFunc {
	return middleware.ServeHTTP(func(ctx context.Context, responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Write([]byte("API Podca "))
	})
}

// EpisodeHandler is a http.HandlerFunc serving FeedEpisodes.
func EpisodeHandler() http.HandlerFunc {
	return middleware.ServeHTTP(func(ctx context.Context, responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Write([]byte("API Podca "))
	})
}

// Feed represents a Podcast Feed.
type Feed struct {
	ID          int64     `json:"id" datastore:"-"`
	Title       string    `json:"title" xml:"title"`
	Link        string    `json:"link" xml:"link"`
	Description string    `json:"description" xml:"description" datastore:",noindex"`
	Language    string    `json:"language" xml:"language"`
	Copyright   string    `json:"copyright" xml:"copyright"`
	PubDate     string    `json:"pubDate" xml:"pubDate"`
	Image       Image     `json:"image" xml:"image"`
	Episodes    []Episode `json:"episodes" xml:"item"`
}

// Episode represents a Podcast Episode.
type Episode struct {
	Title       string `json:"title" xml:"title"`
	Description string `json:"description" xml:"description" datastore:",noindex"`
	Author      string `json:"author" xml:"author"`
	GUID        string `json:"guid" xml:"guid"`
	PubDate     string `json:"pubDate" xml:"pubDate"`
}

// Image is attached to Feed.
type Image struct {
	//Id 	int64	`datastore:"-"`
	FeeURL string `json:"url" xml:"url"`
	Title  string `json:"title" xml:"title"`
	Link   string `json:"link" xml:"link"`
}

// Repository is responsible for CRUD operations on Feeds.
type Repository interface {
	Save(ctx context.Context, feed *Feed) (*Feed, error)
	Get(ctx context.Context, id int64) (*Feed, error)
	GetAll(ctx context.Context) ([]Feed, error)
}
