package api

import (
	"net/http"
	"path"

	"github.com/aws/aws-lambda-go/events"
)

func Route(event events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	pth := event.RawPath
	path, id := path.Split(pth)
	method := event.RequestContext.HTTP.Method

	switch {
	case pth == "/events" && method == "POST":
		return putEvents(PutEventsRequest{Body: event.Body}), nil
	case pth == "/events" && method == "GET":
		return getEvents(), nil
	case path == "/events/" && method == "GET":
		return getEvent(id), nil
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: http.StatusNotFound,
	}, nil
}
