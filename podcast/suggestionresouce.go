package podcast

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
)

type suggestion struct {
	facade SuggestionFacade
	ctx    context.Context
}

func (suggestion *suggestion) suggest(responseWriter http.ResponseWriter, request *http.Request) {
	var term = request.FormValue("term")
	var limit = request.FormValue("limit")
	podcasts, err := suggestion.facade.Suggestion(suggestion.ctx, term, limit)

	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if jsonerr := json.NewEncoder(responseWriter).Encode(podcasts); jsonerr != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		panic(jsonerr)
	}
}

type podcastSuggestion struct {
	SearchEngine SearchEngine
}

// NewPodcastSuggestion creates a new PodcastSuggestionFacade instance
func NewPodcastSuggestion(searchEngine SearchEngine) SuggestionFacade {
	return &podcastSuggestion{
		SearchEngine: searchEngine,
	}
}

func (suggestion *podcastSuggestion) Suggestion(ctx context.Context,
	term string, limit string) ([]Podcast, error) {

	podcasts, err := suggestion.SearchEngine.Search(ctx, term, limit)
	if err != nil {
		return nil, err
	}
	return podcasts, nil
}
