package gitlab

import (
	"net/http"

	"github.com/n0madic/site2rss"
)

func gitlabReleaseRSS() (string, error) {
	return site2rss.NewFeed("https://about.gitlab.com/releases/categories/releases/", "GitLab releases").
		GetLinks(".blog-hero-title, .blog-card-title").
		SetParseOptions(&site2rss.FindOnPage{
			Title:       ".blog.article > div.wrapper.body > h1",
			Author:      "span.author > a:nth-child(1)",
			Date:        "span.date",
			DateFormat:  "Jan 2, 2006",
			Description: ".content",
		}).
		GetItemsFromLinks(site2rss.ParseItem).
		GetRSS()
}

// Handler HTTP
func Handler(w http.ResponseWriter, r *http.Request) {
	rss, err := gitlabReleaseRSS()
	if err == nil {
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(rss))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
}
