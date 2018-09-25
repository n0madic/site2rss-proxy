package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gorilla/feeds"
	"github.com/n0madic/site2rss"
)

const flibustaURL = "http://flibusta.is"

func flibustaRSS(genre string) string {
	rss, err := site2rss.NewFeed(fmt.Sprintf("%s/g/%s/Time", flibustaURL, genre), fmt.Sprintf("Flibusta %s feed", genre)).
		GetLinks("#main > form > ol > a").
		GetFeedItems(func(book *goquery.Document) *feeds.Item {
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
			return &feeds.Item{
				Title:       fmt.Sprintf("%s :: %s", title, author),
				Link:        &feeds.Link{Href: book.Url.String()},
				Id:          book.Url.String(),
				Author:      &feeds.Author{Name: author},
				Description: desc,
				Created:     created,
			}
		}).GetRSS()
	if err != nil {
		return ""
	}
	return rss
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var response events.APIGatewayProxyResponse
	if genre, ok := request.QueryStringParameters["genre"]; ok {
		response = events.APIGatewayProxyResponse{
			Body:       flibustaRSS(genre),
			StatusCode: 200,
			Headers:    map[string]string{"content-type": "application/xml"},
		}
	} else {
		response = events.APIGatewayProxyResponse{
			Body:       "ERROR: genre required!",
			StatusCode: 400,
			Headers:    map[string]string{"content-type": "text/plain"},
		}
	}
	return response, nil
}

func rssRequest(w http.ResponseWriter, r *http.Request) {
	genre := r.URL.Query().Get("genre")
	if genre != "" {
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(flibustaRSS(genre)))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ERROR: genre required!"))
	}
}

func main() {
	if _, ok := os.LookupEnv("LAMBDA_TASK_ROOT"); ok {
		lambda.Start(handleRequest)
	} else {
		http.HandleFunc("/", rssRequest)
		log.Fatal(http.ListenAndServe(":3000", nil))
	}
}
