package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/n0madic/site2rss-proxy"
)

func main() {
	if _, ok := os.LookupEnv("LAMBDA_TASK_ROOT"); ok {
		muxAdapter := gorillamux.New(site2rss.Router)
		lambda.Start(func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			return muxAdapter.ProxyWithContext(ctx, req)
		})
	} else {
		http.Handle("/", site2rss.Router)
		log.Fatal(http.ListenAndServe(":3000", nil))
	}
}
