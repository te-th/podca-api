package podca_api

import(
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func rootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("API Podca "))
	}

}


func podcastSearchHandler(podcastSearcher PodcastSearch) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		ctx := appengine.NewContext(r)
		var term = r.FormValue("term")
		podcasts, err := podcastSearcher.Search(ctx, term) ; if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			panic(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if jsonerr := json.NewEncoder(w).Encode(podcasts); jsonerr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			panic(jsonerr)
		}
	}

}


func feedHandler(feedRepo FeedRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "GET":
			ctx := appengine.NewContext(r)
			vars := mux.Vars(r)

			feedId, _ := strconv.Atoi(vars["feedId"])
			if feedId > 0 {
				log.Infof(ctx, "FEEDID> %d", int64(feedId))
				result, _ := feedRepo.get(ctx, int64(feedId))
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				if jsonerr := json.NewEncoder(w).Encode(result); jsonerr != nil {
					panic(jsonerr)
				}
			} else {
				result, _ :=feedRepo.getAll(ctx)
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				if jsonerr := json.NewEncoder(w).Encode(result); jsonerr != nil {
					panic(jsonerr)
				}
			}

		case "PUT":
			// update
			w.Write([]byte("NOT YET IMPLEMENTED "))
			w.WriteHeader(http.StatusNoContent)
		case "POST":
			// create
			w.Write([]byte("NOT YET IMPLEMENTED "))
			w.WriteHeader(http.StatusNoContent)
		case "DELETE":
			// remove
			w.Write([]byte("NOT YET IMPLEMENTED "))
			w.WriteHeader(http.StatusNoContent)

		}

	}
}

func feedImageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API Podca "))
	}
}

func feedEpisodeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API Podca "))
	}
}

