package podca_api

import (
	"github.com/gorilla/mux"
	"net/http"
)

func init() {
	// Dependencies
	feedRepo := NewFeedRepo()
	searchEngine := NewSearchEngine()
	feedWorker := NewFeedWorker(feedRepo)


	router := mux.NewRouter()

	router.HandleFunc("/", rootHandler())
	router.HandleFunc("/podcasts/search", podcastSearchHandler(searchEngine, feedWorker))
	router.HandleFunc("/feeds", feedHandler(feedRepo))
	router.HandleFunc("/feeds/{feedId}", feedHandler(feedRepo))
	router.HandleFunc("/feed/{feedId}/episodes/{episodeId}", feedEpisodeHandler())
	router.HandleFunc("/feed/{feedId}/image/", feedImageHandler())

	http.Handle("/", router)


}

