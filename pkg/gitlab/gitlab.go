package gitlab

import (
	"net/http"
	"strings"
	"time"

	"github.com/n0madic/site2rss"
)

func gitlabReleaseRSS() (string, error) {
	return site2rss.NewFeed("https://about.gitlab.com/releases/categories/releases/", "GitLab releases").
		GetItemsFromQuery(".blog-card", func(doc *site2rss.Selection, opts *site2rss.FindOnPage) *site2rss.Item {
			url := "https://about.gitlab.com" + doc.Find(".blog-card-title").First().AttrOr("href", "")
			author := doc.Find(".blog-card-author > a").First().Text()
			title := strings.TrimSpace(doc.Find(".blog-card-title > h3").First().Text())
			created, _ := time.Parse("Jan 2, 2006", strings.TrimSpace(doc.Find(".log-card-date").First().Text()))
			desc, _ := doc.Find(".blog-card-excerpt").First().Html()
			desc = strings.TrimSpace(desc)
			return &site2rss.Item{
				Title:       title,
				Author:      &site2rss.Author{Name: author},
				Link:        &site2rss.Link{Href: url},
				Id:          url,
				Description: desc,
				Created:     created,
			}
		}).GetRSS()
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
