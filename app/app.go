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
