package site2rss

import (
	"strings"

	"github.com/n0madic/site2rss-proxy/pkg/flibusta"
	"github.com/n0madic/site2rss-proxy/pkg/pikabu"

	"github.com/gorilla/mux"
)

var Router = NewRouter("")

func NewRouter(prefix string) *mux.Router {
	prefix = strings.TrimSuffix(prefix, "/")
	mux := mux.NewRouter()
	mux.HandleFunc(prefix+"/flibusta/genre/{genre}", flibusta.Handler)
	mux.HandleFunc(prefix+"/pikabu/{location}/{name}", pikabu.Handler)

	return mux
}
