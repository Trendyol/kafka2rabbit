package event_executor_test

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/durmaze/gobank"
	"github.com/durmaze/gobank/predicates"
	"github.com/durmaze/gobank/responses"
	slack_api_client "github.com/kafka2rabbit/pkg/clients"
	"github.com/kafka2rabbit/pkg/rest"
	. "github.com/kafka2rabbit/pkg/util"
	"github.com/kafka2rabbit/services/event_executor"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("ErrorExecuteBehavioral", func() {
	Context("", func() {
		var (
			topic        = "error"
			userName     = "userName"
			slackChannel = "channel"
			request      = slack_api_client.SlackAlert{
				Channel:   slackChannel,
				Username:  userName,
				IconEmoji: "warning",
				Text:      "message",
			}
			hitCount = 0
		)

		BeforeAll(func() {
			marshalledRequest, _ := json.Marshal(request)
			arrangeMountebank(string(marshalledRequest))

			restClient := rest.NewClient()
			claimApiClient := slack_api_client.NewClient(mountebankContainer.Uri(), userName, slackChannel, restClient)

			executor := event_executor.NewEventExecutor()
			executor.SetStrategy(event_executor.NewErrorExecutor(claimApiClient))

			msg := &sarama.ConsumerMessage{
				Value: []byte("message"),
				Topic: topic,
				Key:   []byte("some-aggregateId"),
			}
			err := executor.Execute(msg)
			Expect(err).NotTo(HaveOccurred())

			hitCount, err = Mountebank.NumberOfRequests(mountebankContainer.Port())

		})
		It("it should one time hit to slack api", func() {
			Expect(hitCount).Should(Equal(1))
		})
	})

})

func arrangeMountebank(body string) {
	stub := gobank.Stub().
		Predicates(predicates.Matches().
			Method("POST").
			Path("/").
			Body(body).
			Build()).
		Responses(responses.
			Is().
			StatusCode(http.StatusOK).
			Build()).
		Build()

	HttpStub(stub)
}

func HttpStub(stubs ...gobank.StubElement) {
	imposter := gobank.NewImposterBuilder().
		Protocol("http").
		Port(mountebankContainer.Port()).
		Stubs(stubs...).
		Build()
	Mountebank.DeleteAllImposters()
	if _, err := Mountebank.CreateImposter(imposter); err != nil {
		fmt.Println("Imposter error", err.Error())
	}
}
