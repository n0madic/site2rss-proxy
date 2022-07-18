package pikabu

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/n0madic/site2rss"
)

func pikabuRSS(location string) string {
	rss, err := site2rss.NewFeed("https://pikabu.ru/"+location, fmt.Sprintf("Pikabu %s feed", location)).
		GetItemsFromQuery(".story__main", func(doc *site2rss.Selection, opts *site2rss.FindOnPage) *site2rss.Item {
			sponsor := false
			doc.Find(".story__sponsor").Each(func(i int, c *site2rss.Selection) {
				sponsor = true
			})
			if !sponsor {
				url := doc.Find(".story__title > a").First().AttrOr("href", "")
				author := doc.Find(".user__nick").First().Text()
				title := doc.Find(".story__title-link").First().Text()
				title = site2rss.ConvertToUTF8(title, "windows-1251")
				created, _ := time.Parse(time.RFC3339, doc.Find(".story__datetime").First().AttrOr("datetime", ""))
				desc, _ := doc.Find(".story__content-inner").Each(func(i int, sel *site2rss.Selection) {
					sel.Find(".story-image__image").Each(func(i int, sel *site2rss.Selection) {
						sel.SetAttr("src", sel.AttrOr("data-src", ""))
					})
				}).Html()
				desc = site2rss.ConvertToUTF8(strings.TrimSpace(desc), "windows-1251")
				return &site2rss.Item{
					Title:       title,
					Author:      &site2rss.Author{Name: author},
					Link:        &site2rss.Link{Href: url},
					Id:          url,
					Description: desc,
					Created:     created,
				}
			} else {
				return nil
			}
		}).GetRSS()
	if err != nil {
		return ""
	}
	return rss
}

func Handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["location"] != "" && vars["name"] != "" {
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(pikabuRSS(vars["location"] + "/" + vars["name"])))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ERROR: location required!"))
	}
}
