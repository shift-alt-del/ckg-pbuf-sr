

## Write a lambda function in Python (If you don't have one)

Documents:
- https://aws.github.io/chalice/quickstart.html
- https://aws.github.io/chalice/topics/purelambda.html

```
pip install -r requirements.txt

$ mkdir ~/.aws
$ cat >> ~/.aws/config
[default]
aws_access_key_id=YOUR_ACCESS_KEY_HERE
aws_secret_access_key=YOUR_SECRET_ACCESS_KEY
region=YOUR_REGION (such as us-west-2, us-west-1, etc)
```

## Deploy lambda function

Run `chalice deploy` to deploy/update lambda functions.

```
cd lambda/pbuf-sr
chalice deploy
```

## Deploy Lambda Sink to Confluent Cloud 
```
confluent connect create --config submit.json
```