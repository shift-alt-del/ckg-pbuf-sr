## protobuf files

- Generates source code to `/resorces/generated` folder.

```
protoc -I ./resources --go_out=. item.proto user.proto
```