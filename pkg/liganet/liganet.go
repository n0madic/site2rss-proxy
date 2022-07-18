package liganet

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/n0madic/site2rss"
)

func ligaNewsRSS(tag string) (string, error) {
	return site2rss.NewFeed("https://www.liga.net/tag/"+tag, "Liga.net News").
		GetLinks("div.title > a:nth-child(1)").
		SetParseOptions(&site2rss.FindOnPage{
			Title:       ".article-content > h1",
			Author:      ".author-redactor",
			Date:        ".article-time",
			DateFormat:  "02.01.2006, 15:04",
			Description: "#news-text",
		}).
		GetItemsFromLinks(site2rss.ParseItem).
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
