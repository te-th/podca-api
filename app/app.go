package app

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/te-th/podca-api/api"

	"github.com/te-th/podca-api/domain"
)

func init() {
	// Dependencies PodcastSearcher
	feedRepo := domain.NewFeedRepo()
	feedWorker := api.NewFeedTaskWorker(feedRepo)
	searchEngine := api.NewSearchEngine()

	podcastSearcher := api.NewPodcastSearcher(feedWorker, searchEngine)

	router := mux.NewRouter()

	router.HandleFunc("/", api.RootHandler())
	router.HandleFunc("/podcasts/search", api.PodcastSearchHandler(podcastSearcher))
	router.HandleFunc("/feeds", api.FeedHandler(feedRepo))
	router.HandleFunc("/feeds/{feedId}", api.FeedHandler(feedRepo))
	router.HandleFunc("/feed/{feedId}/episodes/{episodeId}", api.FeedEpisodeHandler())
	router.HandleFunc("/feed/{feedId}/image/", api.FeedImageHandler())

	http.Handle("/", router)
}
