package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go-kafka-protobuf/protobuf"
	"go-kafka-protobuf/resources/generated"
	"go-kafka-protobuf/srclient"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"math/rand"
	"os"
	"time"
)

func main() {
	topicName := "public.cc.sr.pb.demo"

	// setup credentials from system environment.
	bootstrapServers := os.Getenv("BOOTSTRAP_SERVERS")
	saslUsername := os.Getenv("SASL_USERNAME")
	saslPassword := os.Getenv("SASL_PASSWORD")
	srUrl := os.Getenv("SR_URL")
	srUsername := os.Getenv("SR_USERNAME")
	srPassword := os.Getenv("SR_PASSWORD")

	config := &kafka.ConfigMap{
		"bootstrap.servers":                     bootstrapServers,
		"security.protocol":                     "SASL_SSL",
		"sasl.mechanism":                        "PLAIN",
		"sasl.username":                         saslUsername,
		"sasl.password":                         saslPassword,
		"ssl.endpoint.identification.algorithm": "https",
		"enable.ssl.certificate.verification":   "false",

		// producer configurations.
		"acks": "all",
	}

	p, err := kafka.NewProducer(config)
	if err != nil {
		panic(err)
	}
	defer p.Close()

	// register schemaId by given protobuf structure.
	// will return an existed schemaId if the schema has been registered.
	client := srclient.NewClient(srclient.WithURL(srUrl), srclient.WithCredentials(srUsername, srPassword))
	registrator := protobuf.NewSchemaRegistrator(client)
	schemaID, err := registrator.RegisterValue(context.Background(), topicName, &generated.User{})
	if err != nil {
		panic(fmt.Errorf("error registring schmea: %w", err))
	}
	fmt.Printf("Schema registered with schema ID %d\n", schemaID)

	serde := protobuf.NewProtoSerDe()
	deliveryChan := make(chan kafka.Event)
	for i := 0; i < 10000; i++ {

		// create random demo data.
		user := randomDemoData()

		// serialize message.
		// (magic byte) + (schema id 4 bytes) + (index bytes) + (data bytes)
		messageBytes, err := serde.Serialize(schemaID, user)
		if err != nil {
			panic(fmt.Errorf("error serializing message: %w", err))
		}

		// produce message.
		fmt.Printf("Producing message to %s with key: '%s'\n", topicName, user.UserId)
		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topicName, Partition: kafka.PartitionAny},
			Key:            []byte(user.UserId),
			Value:          messageBytes,
		}, deliveryChan)

		// Sync writes VS Async writes
		// https://docs.confluent.io/clients-confluent-kafka-go/current/overview.html#synchronous-writes
		// https://docs.confluent.io/clients-confluent-kafka-go/current/overview.html#asynchronous-writes
		e := <-deliveryChan
		m := e.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			fmt.Printf("Delivery failed to topic: %v\n", m.TopicPartition.Error)
		} else {
			fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
				*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
		}
		randSleep := rand.Intn(3000)
		time.Sleep(time.Duration(randSleep) * time.Millisecond)
	}
}

// Just a function to create dummy data for the demo.
// Each user instance will nest an Address and a random sized Item array.
func randomDemoData() *generated.User {
	randomItemSize := rand.Intn(5)

	// user information.
	user := &generated.User{
		UserId:    uuid.New().String(),
		FirstName: "Cool",
		LastName:  "Guy",
		Message:   "A normal message.",
		Timestamp: timestamppb.Now(),
		Items:     make([]*generated.Item, randomItemSize),
		Address: &generated.User_Address{
			Street:     "Kafka Street",
			PostalCode: "000-0000",
			City:       "Tokyo",
			Country:    "Japan",
		},
	}

	// add random items.
	for itemIndex := 0; itemIndex < randomItemSize; itemIndex++ {
		user.Items[itemIndex] = &generated.Item{
			Name:  fmt.Sprintf("Item-%d", itemIndex),
			Value: fmt.Sprintf("%d", rand.Intn(10000)),
		}
	}
	return user
}
