
# ckg-pbuf-sr

This repo is compatible with the Confluent SR (some are compatible with other SR's such as landoop). Serializes message as `Protobuf`.

## This demo is now archived
confluent-kafka-go `1.9.1` now supports SchemaRegistry both Avro Protocol Buffers and JsonSchema are available, see the [release note](https://github.com/confluentinc/confluent-kafka-go/releases/tag/v1.9.1).

## Demo & Examples

- To produce and consume `Protobuf` message using `confluent-kafka-go` and `sarama`, please refer to [./demo](./demo). This demo uses the Kafka cluster on [Confluent Cloud](https://confluent.cloud).
- A Protobuf + Confluent Lambda Sink connector example under [./demo_connect_lambda](./demo_connect_lambda).
- A Protobuf + Confluent DynamoDB Sink connector example under [./demo_connect_dynamodb](./demo_connect_dynamodb).
- Some ksqlDB scripts under [./demo_ksql](./demo_ksql) for stream processing to play with data in realtime.
- To modify and re-generate `Protobuf` resources, please refer to [./resources/readme.md](./resources)


## References:
- Refactor based on work done by: [go-kafka-protobuf](https://github.com/xtruder/go-kafka-protobuf) - This project is for experimenting of go - kafka - protobuf - confluent schema
registry integration. It should be later refactored in proper library,
whether merging with upstream project or as separate project.
- Refactor from [markteehan](https://github.com/markteehan) 's private repo.
- [Confluent Kafka Go client](https://docs.confluent.io/clients-confluent-kafka-go/current/overview.html#go-example-code) (Protobuf is not officially supported, for any updates, please refer to the Confluent website.)
- [Sarama](https://github.com/Shopify/sarama)
- [Protocol Buffers](https://developers.google.com/protocol-buffers/docs/gotutorial)