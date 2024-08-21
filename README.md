### Instructions

```
docker compose up
```

### Request examples

```
curl -X POST http://localhost:8080/bikes/assign -H "Content-Type: application/json" -d '{"user_uuid":"d0ab33d7-8fcc-463d-bade-fefd53b77a96"}' | jq
curl http://localhost:8080/bikes/available | jq
curl http://localhost:8080/bikes | jq
```

### Run unit tests

```
go test ./...
```
