package podca_api



import(
	"net/http"
	"github.com/golang/appengine"
	"encoding/json"
	"github.com/golang/appengine/log"
	"github.com/gorilla/mux"
	"strconv"
)

func rootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("API Podca "))
	}

}


func podcastSearchHandler(searchEngine *ITunesSearchEngine, worker *FeedWorker) http.HandlerFunc {

	//var laterFunc = delay.Func("key", worker.Retrieve())

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		var term = r.FormValue("term")
		podcasts, err := searchEngine.Search(ctx, term) ; if err != nil {
			panic(err)
		}

		for _, podcast := range podcasts {
			worker.Retrieve(ctx, podcast)

		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if jsonerr := json.NewEncoder(w).Encode(podcasts); jsonerr != nil {
			panic(jsonerr)
		}
	}

}


/*func feedParseHandler(feedRepo *FeedRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := appengine.NewContext(r)
		var feedUrl = "http://www1.wdr.de/radio/podcasts/wdr5/zeitzeichen244.podcast"

		data, fetcherr := fetchData(ctx, feedUrl); if fetcherr != nil {
			panic(fetcherr)
		}

		type Rss struct {
			Feed	Feed `xml:"channel"`
		}

		var rss Rss

		xmlerr := xml.Unmarshal(data, &rss); if xmlerr != nil {
			panic(xmlerr)
			return
		}
		feedRepo.save(ctx, &rss.Feed)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if jsonerr := json.NewEncoder(w).Encode(rss.Feed); jsonerr != nil {
			panic(jsonerr)
		}
	}

}*/



func feedHandler(feedRepo *FeedRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "GET":

			// Get all and return
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
			w.WriteHeader(http.StatusNoContent)
		case "POST":
			// create
			w.WriteHeader(http.StatusNoContent)
		case "DELETE":
			// remove
			w.WriteHeader(http.StatusNoContent)

		}

		//w.Write([]byte("API Podca "))
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

