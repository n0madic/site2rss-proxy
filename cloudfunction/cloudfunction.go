package cloudfunction

import (
	"net/http"

	"github.com/n0madic/site2rss-proxy"
)

func EntryPointHandler(w http.ResponseWriter, r *http.Request) {
	site2rss.Router.ServeHTTP(w, r)
}
