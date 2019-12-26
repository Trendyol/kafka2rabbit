package rabbit_test

import (
	"github.com/kafka2rabbit/pkg/rabbit"
	. "github.com/kafka2rabbit/pkg/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strconv"
	"time"
)

var _ = Describe("When RabbitMQ client", func() {
	Context("failed to connect", func() {
		var (
			connectionParameters rabbit.ConnectionParameters
			actualRetryCount     int
			expectedRetryCount   = 5
		)

		BeforeAll(func() {
			connectionParameters = rabbit.ConnectionParameters{
				ConnectionString: "amqp://invalid_connection_String",
				PrefetchCount:    2,
				RetryCount:       expectedRetryCount,
				RetryInterval:    time.Millisecond * 10,
			}

			var client = rabbit.NewClient(connectionParameters)
			client.Connect()

			for actualRetryCount < expectedRetryCount {
				select {
				case isConnected := <-client.NotifyChannel():
					if !isConnected {
						actualRetryCount++
					}
				case <- time.After(time.Second * 50):
					Fail("Maximum wait duration exceeded")
				}
			}
		})

		It("should retry "+strconv.Itoa(expectedRetryCount)+" times to connect again", func() {
			Expect(actualRetryCount).Should(Equal(expectedRetryCount))
		})
	})
})
