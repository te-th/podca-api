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

package feed

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

// NewFeedRepo creates a new FeedRepository
func NewFeedRepo() Repository {
	return &feedRepo{}
}

type feedRepo struct {
}

func (feedRepo *feedRepo) GetAll(ctx context.Context) ([]Feed, error) {
	feeds := []Feed{}
	ks, err := datastore.NewQuery("Feed").Ancestor(feedKey(ctx)).GetAll(ctx, &feeds)
	if err != nil {
		log.Infof(ctx, "query failed with: %v", err)
		return nil, err
	}
	for i := 0; i < len(feeds); i++ {
		feeds[i].ID = ks[i].IntID()
	}
	return feeds, nil
}

func (feedRepo *feedRepo) Get(ctx context.Context, id int64) (*Feed, error) {
	log.Infof(ctx, "FEED: GET")
	feed := new(Feed)
	feed.ID = id
	k := feed.key(ctx)

	if err := datastore.Get(ctx, k, feed); err != nil {
		log.Infof(ctx, "get failed with: %v", err)
		return nil, err
	}

	feed.ID = k.IntID()

	return feed, nil
}

func feedKey(ctx context.Context) *datastore.Key {
	return datastore.NewKey(ctx, "FeedList", "Default", 0, nil)
}

func (feed *Feed) key(ctx context.Context) *datastore.Key {
	if feed.ID == 0 {
		log.Infof(ctx, "New NewIncompleteKey")
		return datastore.NewIncompleteKey(ctx, "Feed", feedKey(ctx))
	}
	return datastore.NewKey(ctx, "Feed", "", feed.ID, feedKey(ctx))
}

func (feedRepo *feedRepo) Save(ctx context.Context, feed *Feed) (*Feed, error) {

	k, err := datastore.Put(ctx, feed.key(ctx), feed)
	if err != nil {
		log.Infof(ctx, "put faild with: %v", err)
		return nil, err
	}
	feed.ID = k.IntID()
	return feed, nil
}
