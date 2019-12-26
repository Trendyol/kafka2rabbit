package slack_api_client_test

import (
	"encoding/json"
	"github.com/durmaze/gobank"
	"github.com/durmaze/gobank/predicates"
	"github.com/durmaze/gobank/responses"
	slack_api_client "github.com/kafka2rabbit/pkg/clients"
	. "github.com/kafka2rabbit/pkg/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("Client", func() {
	Context("when sent slack alert request", func() {
		var (
			request = slack_api_client.SlackAlert{
				Channel:   slackChannel,
				Username:  userName,
				IconEmoji: "warning",
				Text:      "message",
			}
			err          error
			requestCount = 0
		)
		BeforeAll(func() {
			body, err := json.Marshal(request)
			Expect(err).NotTo(HaveOccurred())
			claimImposter := gobank.NewImposterBuilder().
				Protocol("http").
				Port(9999).
				Stubs(gobank.Stub().
					Predicates(predicates.Matches().
						Method("Post").
						Path("/").
						Body(string(body)).
						Build()).
					Responses(responses.
						Is().
						StatusCode(http.StatusOK).
						Header("Content-Type", "application/json").
						Build()).
					Build(),
					NotFoundStub()).
				Build()
			Mountebank.DeleteAllImposters()
			_, err = Mountebank.CreateImposter(claimImposter)
			Expect(err).NotTo(HaveOccurred())
			err = slackClient.PushNotification("message")
			requestCount, err = Mountebank.NumberOfRequests(mountebankContainer.Port())
		})
		It("should not have any error ", func() {
			Expect(err).NotTo(HaveOccurred())
		})
		It("should one time hit", func() {
			Expect(requestCount).Should(Equal(1))
		})
	})
})