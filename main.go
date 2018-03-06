package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	SLACK_URL = "slack-url"
	CHANNEL   = "my-channel"
)

var (
	ErrNameNotProvided = errors.New("no name was provided in the HTTP body")
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	if len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{}, ErrNameNotProvided
	}

	var iJsonStruct interface{}

	if err := json.Unmarshal([]byte(request.Body), &iJsonStruct); err != nil {
		panic(err)
	}

	parsed, _ := iJsonStruct.(map[string]interface{})
	parsed["channel"] = CHANNEL

	send(parsed)

	return events.APIGatewayProxyResponse{
		Body:       "ok",
		StatusCode: 200,
	}, nil

}

func send(data map[string]interface{}) error {
	rawJson, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(SLACK_URL, "application/json", bytes.NewBuffer(rawJson))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Something bad happened: %s", body)
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
