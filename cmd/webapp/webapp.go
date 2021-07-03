package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jamestjw/feedme/instagram"
)

func InstagramUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	data, err := instagram.FetchFeed(&username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if data.Data.User.IsPrivate {
		http.Error(w, "Cannot fetch feed for private account.", http.StatusUnprocessableEntity)
		return
	}

	content, err := instagram.GenerateOutput(data)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	payload, err := xml.MarshalIndent(content, "  ", "    ")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	payload = []byte(xml.Header + string(payload))

	w.Header().Set("Content-Type", "text/xml")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	var port = flag.Int("p", 8080, "Port to launch webapp")
	flag.Parse()

	r := mux.NewRouter()
	// Set up this way so we can support other feeds in the future
	instagramRouter := r.PathPrefix("/instagram").Subrouter()
	instagramRouter.HandleFunc("/user/{username}", InstagramUserHandler).Methods("GET")

	log.Printf("Listening on port %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), logRequest(r)))
}
