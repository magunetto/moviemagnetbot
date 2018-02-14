package http

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/magunetto/moviemagnetbot/pkg/user"
)

// RunServer register handlers and start HTTP server
func RunServer() {

	// http handlers
	r := mux.NewRouter()
	r.HandleFunc("/tasks/{user}.xml", feedHandler)

	// http loop
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}

func feedHandler(w http.ResponseWriter, r *http.Request) {

	// query user by feed id
	u := &user.User{FeedID: mux.Vars(r)["user"]}
	u, err := u.GetByFeedID()
	if err != nil {
		log.Printf("error while getting user: %s", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	rss, err := u.GenerateFeed()
	if err != nil {
		log.Printf("error while generating feed: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// send the feed
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(rss))
	if err != nil {
		log.Printf("error while sending feed: %s", err)
	}

	err = u.RenewFeedChecked()
	if err != nil {
		log.Printf("error while updating user: %s", err)
	}
}
