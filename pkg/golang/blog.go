package golang

import (
	"net/http"

	"github.com/n0madic/site2rss"
)

func golangBlogRSS() (string, error) {
	return site2rss.NewFeed("https://golang-blog.blogspot.com/", "Golang blog").
		SetParseOptions(&site2rss.FindOnPage{
			Title:       ".post-title",
			Date:        ".date-header > span",
			DateFormat:  "Monday, 2 January 2006 г.",
			Description: ".post-body",
			URL:         ".post-title > a",
		}).
		GetItemsFromQuery(".date-outer", site2rss.ParseQuery).
		FilterItems(site2rss.Filters{
			Selectors: []string{
				".separator",
				".tlglink",
			},
		}).
		GetRSS()
}

// Handler HTTP
func Handler(w http.ResponseWriter, r *http.Request) {
	rss, err := golangBlogRSS()
	if err == nil {
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(rss))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
}
