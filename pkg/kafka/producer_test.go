package kafka_test

import (
	"github.com/Shopify/sarama"
	"github.com/kafka2rabbit/pkg/kafka"
	. "github.com/kafka2rabbit/pkg/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("When publishing a message", func() {

	Context("and the broker is reachable", func() {
		var (
			message = "test"
			conf    = kafka.ConnectionParameters{
				Version:         "2.2.0",
				ConsumerGroupID: "some-id",
				ClientID:        "oms-event-generator",
				Topic:           "createClaim",
			}
			receivedPayload string
			err             error
			producer        sarama.SyncProducer
		)

		BeforeAll(func() {
			conf.Brokers = KafkaContainer.Address()
			producer, err = kafka.NewProducer(conf)
			Expect(err).NotTo(HaveOccurred())
			_, _, err = producer.SendMessage(&sarama.ProducerMessage{
				Value: sarama.StringEncoder(message),
				Topic: conf.Topic,
			})
			time.Sleep(5*time.Second)
			_, receivedPayload, _ = Consume(conf)
		})

		It("should not produce an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("should send the message with expected payload", func() {
			Expect(receivedPayload).Should(Equal(message))
		})
	})

	Context("and the broker is unreachable", func() {
		var (
			wrongConf = kafka.ConnectionParameters{
				Version: "2.2.0",
				Brokers: "localhost:9093",
				Topic:   "createClaim",
			}
			expectedError error
		)

		BeforeAll(func() {
			_, expectedError = kafka.NewProducer(wrongConf)
		})

		It("should produce an error", func() {
			Expect(expectedError).To(HaveOccurred())
		})

		It("should produce the expected error", func() {
			Expect(expectedError.Error()).Should(Equal("kafka: client has run out of available brokers to talk to (Is your cluster reachable?)"))
		})
	})

})
