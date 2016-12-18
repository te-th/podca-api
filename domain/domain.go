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

package domain

import "golang.org/x/net/context"

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

type Episode struct {
	Title       string `json:"title" xml:"title"`
	Description string `json:"description" xml:"description" datastore:",noindex"`
	Author      string `json:"author" xml:"author"`
	GUID        string `json:"guid" xml:"guid"`
	PubDate     string `json:"pubDate" xml:"pubDate"`
}

type Image struct {
	//Id 	int64	`datastore:"-"`
	FeeURL string `json:"url" xml:"url"`
	Title  string `json:"title" xml:"title"`
	Link   string `json:"link" xml:"link"`
}

// Podcast struct is strong coupled to the Apple iTunes format
type Podcast struct {
	ID             int64    `json:"id"`
	ArtistName     string   `json:"artistName"`
	CollectionName string   `json:"collectionName"`
	FeedURL        string   `json:"feedUrl"`
	CollectionID   int64    `json:"collectionId"`
	TrackID        int64    `json:"trackId"`
	Genres         []string `json:"genres"`
}

type FeedRepository interface {
	Save(ctx context.Context, feed *Feed) (*Feed, error)
	Get(ctx context.Context, id int64) (*Feed, error)
	GetAll(ctx context.Context) ([]Feed, error)
}
