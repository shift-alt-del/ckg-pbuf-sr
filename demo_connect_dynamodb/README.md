
## You have to create a DynamoDB table first.
- Make sure you have correct aws cli installed and credentials setup. 

```
./demo_connect_dynamodb/01_init
```


## Deploy DynamoDB Sink to Confluent Cloud
```
confluent connect create --config submit.json
```