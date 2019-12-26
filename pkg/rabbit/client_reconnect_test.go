package rabbit_test

import (
	"github.com/kafka2rabbit/pkg/rabbit"
	. "github.com/kafka2rabbit/pkg/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("When RabbitMQ client connection is dropped", func() {
	Context("and tries to reconnect again", func() {
		var (
			connectionParameters   rabbit.ConnectionParameters
			actualRetryCount       int
			isClientConnectedAgain bool
			maxRetryCount          = 30
		)

		BeforeAll(func() {
			connectionParameters = rabbit.ConnectionParameters{
				ConnectionString: rabbitMQContainer.Address(),
				PrefetchCount:    2,
				RetryCount:       maxRetryCount,
				RetryInterval:    time.Second,
			}

			var client = rabbit.NewClient(connectionParameters)
			client.Connect()
			<-client.NotifyChannel()

			rabbitMQContainer.Stop()

			rabbitMQContainer.Run()

			for isClientConnectedAgain == false {
				select {
				case status := <-client.NotifyChannel():
					actualRetryCount++
					if actualRetryCount > maxRetryCount {
						Fail("Maximum retry count exceeded")
					}
					if status {
						isClientConnectedAgain = true
					}
				case <-time.After(time.Second * 50):
					Fail("Maximum wait duration exceeded")
				}
			}
		})

		It("should retry to connect again", func() {
			Expect(actualRetryCount > 0).Should(BeTrue())
		})

		It("should reconnect successfully", func() {
			Expect(isClientConnectedAgain).Should(BeTrue())
		})
	})
})
