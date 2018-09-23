package main

import (
	"fmt"
	"log"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
)

const flibustaURL = "http://flibusta.is"
const maxFeedItems = 10

func absoluteURL(rpath string) string {
	u, err := url.Parse(flibustaURL)
	if err != nil {
		panic("invalid url")
	}
	u.Path = path.Join(u.Path, rpath)
	return u.String()
}

// FlibustaRSS create RSS feed by genre
func FlibustaRSS(genre string) string {
	doc, err := goquery.NewDocument(fmt.Sprintf("%s/g/%s/Time", flibustaURL, genre))
	if err != nil {
		log.Fatal(err)
	}

	genreTitle := doc.Find("h1.title").First().Text()

	feed := &feeds.Feed{
		Title:       fmt.Sprintf("Flibusta %s feed", genre),
		Link:        &feeds.Link{Href: flibustaURL},
		Description: genreTitle,
	}

	links := doc.Find("#main > form > ol > a").Map(func(i int, s *goquery.Selection) string {
		bookLink, _ := s.Attr("href")
		return bookLink
	})[:maxFeedItems]

	reAdded, _ := regexp.Compile(`Добавлена: (\S+)`)

	for _, link := range links {
		url := absoluteURL(link)
		book, err := goquery.NewDocument(url)
		if err != nil {
			log.Fatal(err)
		}
		author := book.Find("#main > a:nth-child(5)").First().Text()
		title := strings.TrimSuffix(book.Find("#main > h1").First().Text(), " (fb2)")
		added := reAdded.FindStringSubmatch(book.Text())[1]
		created, _ := time.Parse("02.01.2006", added)
		annotation, _ := book.Find("#main > p").First().Html()
		rating := book.Find("#newann > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(1) > p:nth-child(1)").First().Text()
		desc := fmt.Sprintf("%s<br><br><h5>%s</h5>", annotation, rating)
		if cover, ok := book.Find("#main > img").First().Attr("src"); ok {
			desc = fmt.Sprintf(`<img src="%s" width="400" align="left" hspace="10"> %s`, absoluteURL(cover), desc)
		}
		feed.Add(&feeds.Item{
			Title:       fmt.Sprintf("%s :: %s", title, author),
			Link:        &feeds.Link{Href: url},
			Id:          url,
			Author:      &feeds.Author{Name: author},
			Description: desc,
			Created:     created,
		})
	}

	rss, err := feed.ToRss()
	if err != nil {
		log.Fatal(err)
	}
	return rss
}
