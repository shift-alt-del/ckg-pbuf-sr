// Refactor from:
// https://github.com/Shopify/sarama/blob/main/examples/sasl_scram_client/main.go

package main

import (
	"flag"
	"fmt"
	"go-kafka-protobuf/protobuf"
	"go-kafka-protobuf/resources/generated"
	"log"
	"os"
	"os/signal"
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

	// Consumer configurations
	// OffsetNewest int64 = -1
	// OffsetOldest int64 = -2
	conf.Consumer.Offsets.Initial = -2

	// create consumer instance.
	consumer, err := sarama.NewConsumer(splitBrokers, conf)
	if err != nil {
		panic(err)
	}
	log.Println("consumer created")
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	log.Println("commence consuming")
	partitionConsumer, err := consumer.ConsumePartition(*topic, 2, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumed := 0

	// protobuf serde
	serde := protobuf.NewProtoSerDe()

ConsumerLoop:
	for {
		log.Println("in the for")
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Consumed message offset %d\n", msg.Offset)

			// deserialize protobuf message
			user := &generated.User{}
			_, err := serde.Deserialize(msg.Value, user)
			if err != nil {
				fmt.Printf("Error deserializing message: %v\n", err)
			}
			log.Printf("KEY: %s VALUE: %s", msg.Key, user)

			consumed++
		case <-signals:
			break ConsumerLoop
		}
	}

	log.Printf("Consumed: %d\n", consumed)

	logger.Println("Bye now !")
}
