
# ckg-pbuf-sr
This repo is
* compatible with the Confluent SR (some are compatible with other SR's such as landoop)
* using confluent-kafka-go client
* serializes as protobuf


## Project structure

**Protobuf resources**

`/resources`
- protobuf files
- run `protoc -I ./resources --go_out=. item.proto user.proto` to generate code, source code will be generated to `/resorces/generated` folder.

**Producer & Consumer demo**

`/demo/ccsr_producer.go`
- connects to Confluent Cloud using the (embedded) configuration
- connects to the CCSR using basic auth
- registers a protobuf schema for topic `public.cc.sr.pb.demo`, if one doesnt exist
- produce messages into topic `public.cc.sr.pb.demo`, serialized as pbuf using sync; returning success/failure and the offset

`/demo/ccsr_consumer.go`
- connects to Confluent Cloud using the (embedded) configuration
- consumes messages from topic `public.cc.sr.pb.demo`, deserialized from pbuf.

**AWS Lambda sink connector demo**

`demo_connect_lambda/pbuf-sr`
- a demo Lambda function using python chalice package, it prints out both parameters and shows batch.size.

`demo_connect_lambda/submit.json`
- a json file for Lambda sink connector deploy.
- run `confluent connect create --config submit.json` to deploy connector.

**AWS Dynamodb sink connector demo** 

`demo_connect_dynamodb/01_init`
- creates dynamodb table for dynamodb sink demo.

`demo_connect_dynamodb/submit.json`
- a json file for Dynamodb sink connector deploy.
- run `confluent connect create --config submit.json` to deploy connector.

**ksqlDB demo**

`demo_ksql/01_init`
- a simple script to help create cluster and ksqlDB cluster on CC.

`demo_ksql/statements.sql`
- some simple KSQL scripts to process data in realtime.

## Verify result:
```
TOPIC=public.cc.sr.pb.demo
KAFKA=pkc...
SR_URL=https://psrc...
SR_AUTH=xxx:xxx
CLIENT_CONFIG=/var/tmp/pbuf_go.config.2

kafka-protobuf-console-consumer --bootstrap-server $KAFKA                                  \ 
                                --topic ${TOPIC}                                           \ 
                                --property schema.registry.url=${SR_URL}                   \
                                --property basic.auth.credentials.source=USER_INFO         \
                                --property schema.registry.basic.auth.user.info="$SR_AUTH" \
                                --from-beginning                                           \
                                --consumer.config ${CLIENT_CONFIG}
```

## References:
- Refactor based on work done by: https://github.com/xtruder/go-kafka-protobuf - This project is for experimenting of go - kafka - protobuf - confluent schema
registry integration. It should be later refactored in proper library,
whether merging with upstream project or as separate project.
- Refactor from [markteehan](https://github.com/markteehan) 's private repo.
- [Confluent Kafka Go client](https://docs.confluent.io/clients-confluent-kafka-go/current/overview.html#go-example-code) (Protobuf is not officially supported, for any updates, please refer to the Confluent website.)
- [Protocol Buffers](https://developers.google.com/protocol-buffers/docs/gotutorial)