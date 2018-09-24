package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var response events.APIGatewayProxyResponse
	if genre, ok := request.QueryStringParameters["genre"]; ok {
		response = events.APIGatewayProxyResponse{
			Body:       FlibustaRSS(genre),
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
		w.Write([]byte(FlibustaRSS(genre)))
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
