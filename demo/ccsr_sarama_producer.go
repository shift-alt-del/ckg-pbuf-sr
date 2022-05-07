// Refactor from:
// https://github.com/Shopify/sarama/blob/main/examples/sasl_scram_client/main.go

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"go-kafka-protobuf/protobuf"
	"go-kafka-protobuf/resources/generated"
	"go-kafka-protobuf/srclient"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/Shopify/sarama"
)

func init() {
	sarama.Logger = log.New(os.Stdout, "[Sarama] ", log.LstdFlags)
}

func main() {
	brokers := flag.String("brokers", os.Getenv("BOOTSTRAP_SERVERS"), "The Kafka brokers to connect to, as a comma separated list")
	userName := flag.String("username", os.Getenv("SASL_USERNAME"), "The SASL username")
	passwd := flag.String("passwd", os.Getenv("SASL_PASSWORD"), "The SASL password")
	topic := flag.String("topic", "public.cc.sr.pb.demo", "The Kafka topic to use")
	logger := log.New(os.Stdout, "[Producer] ", log.LstdFlags)

	srUrl := os.Getenv("SR_URL")
	srUsername := os.Getenv("SR_USERNAME")
	srPassword := os.Getenv("SR_PASSWORD")

	// The broker list must be an array type.
	splitBrokers := strings.Split(*brokers, ",")

	// Create sarama configuration
	conf := sarama.NewConfig()
	conf.Metadata.Full = true
	conf.Version = sarama.V0_10_0_0
	conf.ClientID = "ccsr_sarama"
	conf.Metadata.Full = true

	conf.Net.SASL.Mechanism = "" // SASL_PLAIN for the default SASL mechanism
	conf.Net.SASL.Enable = true
	conf.Net.SASL.User = *userName
	conf.Net.SASL.Password = *passwd
	conf.Net.TLS.Enable = true

	// Producer configurations
	conf.Producer.Retry.Max = 1
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Return.Successes = true

	// create producer instance.
	syncProducer, err := sarama.NewSyncProducer(splitBrokers, conf)
	if err != nil {
		logger.Fatalln("failed to create producer: ", err)
	}

	// register schemaId by given protobuf structure.
	// will return an existed schemaId if the schema has been registered.
	client := srclient.NewClient(srclient.WithURL(srUrl), srclient.WithCredentials(srUsername, srPassword))
	registrator := protobuf.NewSchemaRegistrator(client)
	schemaID, err := registrator.RegisterValue(context.Background(), *topic, &generated.User{})
	if err != nil {
		panic(fmt.Errorf("error registring schmea: %w", err))
	}
	fmt.Printf("Schema registered with schema ID %d\n", schemaID)

	serde := protobuf.NewProtoSerDe()

	// create random demo data.
	user := randomDemoDataSarama()

	// serialize message.
	// (magic byte) + (schema id 4 bytes) + (index bytes) + (data bytes)
	messageBytes, err := serde.Serialize(schemaID, user)

	// produce protobuf message.
	// instead of using StringEncoder, use ByteEncoder to send raw bytes which has been serialized as protobuf.
	partition, offset, err := syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: *topic,
		Value: sarama.ByteEncoder(messageBytes),
	})
	if err != nil {
		logger.Fatalln("failed to send message to ", *topic, err)
	}
	logger.Printf("wrote message at partition: %d, offset: %d", partition, offset)
	_ = syncProducer.Close()

	logger.Println("Bye now !")
}

// Just a function to create dummy data for the demo.
// Each user instance will nest an Address and a random sized Item array.
func randomDemoDataSarama() *generated.User {
	randomItemSize := rand.Intn(5)

	// user information.
	user := &generated.User{
		UserId:    uuid.New().String(),
		FirstName: "Sa",
		LastName:  "Rama",
		Message:   "A normal message from Sarama client.",
		Timestamp: timestamppb.Now(),
		Items:     make([]*generated.Item, randomItemSize),
		Address: &generated.User_Address{
			Street:     "Kafka Street - Sarama",
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
