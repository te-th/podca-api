package podca_api

import (
	"github.com/gorilla/mux"
	"net/http"
)

func init() {


	// Dependencies PodcastSearcher
	feedRepo := NewFeedRepo()
	feedWorker := NewFeedTaskWorker(feedRepo)
	searchEngine := NewSearchEngine()

	podcastSearcher := NewPodcastSearcher(feedWorker, searchEngine)

	router := mux.NewRouter()

	router.HandleFunc("/", rootHandler())
	router.HandleFunc("/podcasts/search", podcastSearchHandler(podcastSearcher))
	router.HandleFunc("/feeds", feedHandler(feedRepo))
	router.HandleFunc("/feeds/{feedId}", feedHandler(feedRepo))
	router.HandleFunc("/feed/{feedId}/episodes/{episodeId}", feedEpisodeHandler())
	router.HandleFunc("/feed/{feedId}/image/", feedImageHandler())

	http.Handle("/", router)


}

