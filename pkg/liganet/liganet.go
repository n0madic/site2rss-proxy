package liganet

import (
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/n0madic/site2rss"
)

func ligaNewsRSS(tag string) (string, error) {
	return site2rss.NewFeed("https://www.liga.net/tag/"+tag, "Liga.net News").
		GetLinks("a.news-card__title").
		GetItemsFromLinks(func(doc *site2rss.Document, opts *site2rss.FindOnPage) *site2rss.Item {
			link := doc.Url.String()
			title := strings.TrimSpace(doc.Find(".article-header__title").First().Text())
			author := strings.TrimSpace(doc.Find(".article-header__user .user__name").First().Text())
			created, _ := time.Parse(time.RFC3339, doc.Find("time.article-header__date").First().AttrOr("datetime", ""))
			desc, _ := doc.Find(".article-body").First().Html()
			return &site2rss.Item{
				Title:       title,
				Link:        &site2rss.Link{Href: link},
				Id:          link,
				Author:      &site2rss.Author{Name: author},
				Description: desc,
				Created:     created,
			}
		}).
		FilterItems(site2rss.Filters{
			Text: []string{
				"Читайте нас в Telegram:",
				"Читайте нас у Telegram:",
				"Читайте также:",
				"Читайте також:",
				"Ctrl+Enter",
			},
		}).
		GetRSS()
}

// Handler HTTP
func Handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["tag"] != "" {
		rss, err := ligaNewsRSS(vars["tag"])
		if err == nil {
			w.Header().Set("Content-Type", "application/xml")
			w.Write([]byte(rss))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ERROR: tag required!"))
	}
}
