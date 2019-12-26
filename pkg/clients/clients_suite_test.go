package slack_api_client_test

import (
	"github.com/durmaze/gobank"
	"github.com/kafka2rabbit/internal/integrationcontainers/mountebank"
	slack_api_client "github.com/kafka2rabbit/pkg/clients"
	"github.com/kafka2rabbit/pkg/rest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestClients(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Clients Suite")
}

var (
	userName            = "userName"
	slackChannel        = "slackChannel"
	mountebankContainer *mountebank.Container
	slackClient         slack_api_client.Client
	Mountebank          *gobank.Client
)

var _ = BeforeSuite(func() {
	mountebankContainer = mountebank.NewContainer("andyrbell/mountebank:2.1.0")
	mountebankContainer.Run()

	Mountebank = gobank.NewClient(mountebankContainer.AdminUri())
	slackClient = slack_api_client.NewClient(mountebankContainer.Uri(), userName, slackChannel, rest.NewClient())
})

var _ = AfterSuite(func() {
	mountebankContainer.ForceRemoveAndPrune()
})
