package podca_api

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func NewFeedRepo() *FeedRepo {
	return &FeedRepo{}
}

type FeedRepo struct {

}

func (feedRepo *FeedRepo) getAll(ctx context.Context) ([]Feed, error){
	feeds := []Feed{}
	ks, err := datastore.NewQuery("Feed").Ancestor(feedKey(ctx)).GetAll(ctx, &feeds)
	if err != nil {
		log.Infof(ctx, "query failed with: %v",err)
		return nil, err
	}
	for i := 0; i < len(feeds); i++ {
		feeds[i].Id = ks[i].IntID()
	}
	return feeds, nil
}

func (feedRepo *FeedRepo) get(ctx context.Context, id int64) (*Feed, error){
	log.Infof(ctx, "FEED: GET")
	feed := new(Feed)
	feed.Id = id
	k := feed.key(ctx)

	if err := datastore.Get(ctx,k, feed); err != nil {
		log.Infof(ctx, "get failed with: %v",err)
		return nil, err
	}

	feed.Id = k.IntID()

	return feed, nil
}

func feedKey(ctx context.Context) *datastore.Key  {
	return datastore.NewKey(ctx, "FeedList", "Default", 0, nil)
}

func (feed *Feed) key(ctx context.Context) *datastore.Key {
	if feed.Id == 0 {
		log.Infof(ctx, "New NewIncompleteKey")
		return datastore.NewIncompleteKey(ctx, "Feed", feedKey(ctx))
	}
	return datastore.NewKey(ctx, "Feed", "", feed.Id, feedKey(ctx))
}


func (feedRepo *FeedRepo) save(ctx context.Context, feed *Feed) (*Feed, error) {

	k, err := datastore.Put(ctx, feed.key(ctx), feed)
	if err != nil {
		log.Infof(ctx, "put faild with: %v",err)
		return nil, err
	}
	feed.Id = k.IntID()
	return feed, nil
}

