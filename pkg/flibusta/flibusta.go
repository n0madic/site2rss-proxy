package flibusta

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/n0madic/site2rss"
)

func flibustaRSS(genre string) string {
	rss, err := site2rss.NewFeed(fmt.Sprintf("%s/g/%s/Time", "http://flibusta.is", genre), fmt.Sprintf("Flibusta %s feed", genre)).
		GetLinks("#main > form > ol > a").
		GetItemsFromLinks(func(book *site2rss.Document, opts *site2rss.FindOnPage) *site2rss.Item {
			reAdded := regexp.MustCompile(`Добавлена: (\S+)`)
			author := book.Find("#main > a:nth-child(5)").First().Text()
			title := strings.TrimSuffix(book.Find("#main > h1").First().Text(), " (fb2)")
			created, _ := time.Parse("02.01.2006", reAdded.FindStringSubmatch(book.Text())[1])
			annotation, _ := book.Find("#main > p").First().Html()
			rating := book.Find("#newann > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(1) > p:nth-child(1)").First().Text()
			desc := fmt.Sprintf("%s<br><br><h5>%s</h5>", annotation, rating)
			if cover, ok := book.Find("#main > img").First().Attr("src"); ok {
				desc = fmt.Sprintf(`<img src="%s" width="400" align="left" hspace="10"> %s`, cover, desc)
			}
			return &site2rss.Item{
				Title:       fmt.Sprintf("%s :: %s", title, author),
				Link:        &site2rss.Link{Href: book.Url.String()},
				Id:          book.Url.String(),
				Author:      &site2rss.Author{Name: author},
				Description: desc,
				Created:     created,
			}
		}).GetRSS()
	if err != nil {
		return ""
	}
	return rss
}

func Handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["genre"] != "" {
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(flibustaRSS(vars["genre"])))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ERROR: genre required!"))
	}
}
