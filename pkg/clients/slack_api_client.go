package slack_api_client

import (
	"encoding/json"
	"fmt"
	"github.com/kafka2rabbit/pkg/rest"
	"net/http"
)

const alertEmoji = "warning"

type slackApiClient struct {
	address      string
	userName     string
	slackChannel string
	restClient   *rest.RestClient
}

type Client interface {
	PushNotification(text string) error
}

func NewClient(address, userName, slackChannel string, restClient *rest.RestClient) Client {
	return &slackApiClient{
		address:      address,
		userName:     userName,
		slackChannel: slackChannel,
		restClient:   restClient,
	}
}

func (s *slackApiClient) PushNotification(text string) error {
	var (
		resp *http.Response
		err  error
	)
	request := SlackAlert{
		Channel:   s.slackChannel,
		Username:  s.userName,
		IconEmoji: alertEmoji,
		Text:      text,
	}
	marshalledRequest, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("request marshalling error , err:%v", err)
	}
	resp, err = s.restClient.Post(s.address, marshalledRequest).End()
	if err != nil || resp == nil {
		return fmt.Errorf("error from server")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status not ok")
	}

	return nil
}
