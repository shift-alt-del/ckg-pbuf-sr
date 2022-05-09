
## How to run

- Make sure you have Go installed.
- Please config system environment variables as below. You can get all information from https://confluent.cloud.
- You need create the topic `public.cc.sr.pb.demo` before you run the code.
- The following steps will produce to (or consume from) the topic `public.cc.sr.pb.demo` in `Protobuf` format.

```
export BOOTSTRAP_SERVERS=example.confluent.cloud:9092
export SASL_USERNAME=(get from confluent.cloud)
export SASL_PASSWORD=(get from confluent.cloud)
export SR_URL=https://example.confluent.cloud
export SR_USERNAME=(get from confluent.cloud)
export SR_PASSWORD=(get from confluent.cloud)

# for producer
go run ccsr_producer.go
go run ccsr_sarama_producer.go

# for consumer
go run ccsr_consumer.go
go run ccsr_sarama_consumer.go 
```

## Confluent Cloud + Schema-Registry + Protobuf + Kafka Producer

Files:
- `/demo/ccsr_producer.go`
- `/demo/ccsr_sarama_producer.go`

The main steps:
- connects to Confluent Cloud using the (embedded) configuration
- connects to the CCSR using basic auth
- registers a protobuf schema for topic `public.cc.sr.pb.demo`, if one doesnt exist
- produce messages into topic `public.cc.sr.pb.demo`, serialized as pbuf using sync; returning success/failure and the offset

## Confluent Cloud + Protobuf + Kafka Consumer

Files:
- `/demo/ccsr_consumer.go`
- `/demo/ccsr_sarama_consumer.go`

The main steps:
- connects to Confluent Cloud using the (embedded) configuration
- consumes messages from topic `public.cc.sr.pb.demo`, deserialized from pbuf.
