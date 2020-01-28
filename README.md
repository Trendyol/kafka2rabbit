# Description
**It transfers messages from kafka topics to rabbit exchanges with given configuration . Rabbit producer sends messages with publish confirmation so that it provides eventual consistency . If an error occurs while sending messages to rabbit , message goes to kafka retry topic then failed messages are retried for a given number of times via a retry topic. If it still fails to send the message with retires, it moves to message to an error topic and sends a slack notification about the failure with given configuration.**

# Required environment variables
	APP_PORT : exposed port , for health check
	BROKERS :  kafka cluster address with port (ex: 1.0.0.1:9092,1.0.0.2:9092)
	USERNAME : kafka username
	PASSWORD : kafka password
	KAFKA_VERSION : kafka version (ex: 2.0.0)
	RABBIT_ADDRESS : rabbit address (ex: amqp://usr:passw@addr:5672/virtual-host)
	SLACK_URL: which address to be sent
	SLACK_USERNAME : which user to be sent
	SLACK_CHANNEL : which channel to be sent

# Required application configs
    Configs : [ {
    Topic : “kafka topic”,
    Exchange : “rabbit exchange”
    ExchangeKind: “topic-direct-fanout”
    RoutingKey: “exchange routing key if exist”
    }]
