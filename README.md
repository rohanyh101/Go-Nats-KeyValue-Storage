# Basic NATS Key-Value Store Implementation in Golang
This project implements a basic key-value store using the NATS messaging system in Golang. The key-value store is designed to be a service discovery mechanism for microservices, allowing them to register and discover each other's endpoints.

## The key-value store consists of two main components:
1. KV: The KV component handles incoming requests to put and delete key-value pairs. It publishes these updates to a NATS stream, which is consumed by the Watcher component.
2. Watcher: The Watcher component subscribes to the NATS stream and processes updates to the key-value store. It prints the updated key-value pairs to the console.

## Run Instructions,
1. Start the Watcher component by running the `watcher/main`.
2. Start the KV component by running the `kv/main.go` in a separate terminal window.

- Command: #1,
```powershell
  go run watcher/main.go
```
- Response: #1,
```powershell
services.orders @ 1 -> "https://localhost:8080/orders" (op: KeyValuePutOp)
services.orders @ 2 -> "https://localhost:8080/v1/orders" (op: KeyValuePutOp)
services.consumers @ 3 -> "http://localhost:8080/v2/consumers" (op: KeyValuePutOp)
services.consumers @ 4 -> "" (op: KeyValueDeleteOp)
```

- Command: #2,
```powershell
  go run kv/main.go
```
- Response: #2,
```powershell
services.orders @ 1 -> "https://localhost:8080/orders"
services.orders @ 2 -> "https://localhost:8080/v1/orders"
services.orders @ 2 -> "https://localhost:8080/v1/orders"
KV stream name: KV_discovery
services.consumers @ 3 -> "http://localhost:8080/v2/consumers"
```

