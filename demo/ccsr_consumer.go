package main

import (
	"fmt"
	"go-kafka-protobuf/protobuf"
	"go-kafka-protobuf/resources/generated"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"os"
)

func main() {
	topicName := "public.cc.sr.pb.demo"

	// setup credentials from system environment.
	// protobuf messages are self explained, do not need SR configurations.
	bootstrapServers := os.Getenv("BOOTSTRAP_SERVERS")
	saslUsername := os.Getenv("SASL_USERNAME")
	saslPassword := os.Getenv("SASL_PASSWORD")

	config := &kafka.ConfigMap{
		"bootstrap.servers":                     bootstrapServers,
		"security.protocol":                     "SASL_SSL",
		"sasl.mechanism":                        "PLAIN",
		"sasl.username":                         saslUsername,
		"sasl.password":                         saslPassword,
		"ssl.endpoint.identification.algorithm": "https",
		"enable.ssl.certificate.verification":   "false",

		// consumer configurations.
		"group.id":          "something",
		"auto.offset.reset": "earliest",
	}

	c, err := kafka.NewConsumer(config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created new consumer")

	err = c.Subscribe(topicName, nil)

	serde := protobuf.NewProtoSerDe()
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s with key: '%s'\n", msg.TopicPartition, msg.Key)
			user := &generated.User{}

			_, err := serde.Deserialize(msg.Value, user)
			if err != nil {
				fmt.Printf("Error deserializing message: %v\n", err)
			}

			fmt.Printf("Received message: %v\n", user)
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
