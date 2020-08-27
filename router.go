package site2rss

import (
	"github.com/n0madic/site2rss-proxy/pkg/flibusta"
	"github.com/n0madic/site2rss-proxy/pkg/pikabu"

	"github.com/gorilla/mux"
)

var Router = newMux()

func newMux() *mux.Router {
	mux := mux.NewRouter()
	mux.HandleFunc("/site2rss/flibusta/genre/{genre}", flibusta.Handler)
	mux.HandleFunc("/site2rss/pikabu/{location}/{name}", pikabu.Handler)

	return mux
}
