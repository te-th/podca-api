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

package api

import (
	"net/http"
	"sync"

	"github.com/te-th/podca-api/feed"
	"github.com/te-th/podca-api/middleware"
	"github.com/te-th/podca-api/podcast"
	"github.com/te-th/podca-api/root"
)

var once sync.Once

// Handler creates http.HandlerFunc for endpoints
func Handler() map[string]http.HandlerFunc {
	handlers := make(map[string]http.HandlerFunc)

	// register handler exactly once
	once.Do(func() {
		httpClient := middleware.NewHTTPClient()
		feedRepo := feed.NewFeedRepo()
		feedTask := feed.NewFeedTask(feedRepo, httpClient)
		searchEngine := podcast.NewSearchEngine(httpClient)

		podcastSuggestion := podcast.NewPodcastSuggestion(searchEngine)
		podcastSearch := podcast.NewPodcastSearch(feedTask, searchEngine)

		handlers["/"] = root.Resource()
		handlers["/podcasts/search"] = podcast.SearchHandler(podcastSearch)
		handlers["/podcasts/suggestion"] = podcast.SuggestionHandler(podcastSuggestion)
		handlers["/feeds"] = feed.Handler(feedRepo)

		handlers["/feeds/{feedId}"] = feed.Handler(feedRepo)
		handlers["/feed/{feedId}/episodes/{episodeId}"] = feed.EpisodeHandler()
		handlers["/feed/{feedId}/image/"] = feed.ImageHandler()
	})

	return handlers
}
