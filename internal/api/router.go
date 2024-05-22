package api

import (
	"net/http"
	"path"

	"github.com/aws/aws-lambda-go/events"
)

func Route(event events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	pth := event.RequestContext.HTTP.Path
	method := event.RequestContext.HTTP.Method

	switch {
	case pth == "/events" && method == "POST":
		return putEvents(PutEventsRequest{Body: event.Body}), nil
	case pth == "/events" && method == "GET":
		return getEvents(), nil
	default:
		path, id := path.Split(pth)
		if path == "/events/" && method == "GET" {
			return getEvent(id), nil
		}
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: http.StatusNotFound,
	}, nil
}
