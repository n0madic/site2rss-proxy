package site2rss

import (
	"strings"

	"github.com/n0madic/site2rss-proxy/pkg/flibusta"
	"github.com/n0madic/site2rss-proxy/pkg/gitlab"
	"github.com/n0madic/site2rss-proxy/pkg/golang"
	"github.com/n0madic/site2rss-proxy/pkg/liganet"
	"github.com/n0madic/site2rss-proxy/pkg/pikabu"

	"github.com/gorilla/mux"
)

var Router = NewRouter("")

func NewRouter(prefix string) *mux.Router {
	prefix = strings.TrimSuffix(prefix, "/")
	mux := mux.NewRouter()
	mux.HandleFunc(prefix+"/flibusta/genre/{genre}", flibusta.Handler)
	mux.HandleFunc(prefix+"/gitlab/blog/releases", gitlab.Handler)
	mux.HandleFunc(prefix+"/golang/blog", golang.Handler)
	mux.HandleFunc(prefix+"/liganet/tag/{tag}", liganet.Handler)
	mux.HandleFunc(prefix+"/pikabu/{location}/{name}", pikabu.Handler)

	return mux
}
