package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

func main() {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatalf("error connecting to nats: %v", err)
	}
	defer nc.Drain()

	jsCtx, err := nc.JetStream()
	if err != nil {
		log.Fatalf("error creating jet stream: %v", err)
	}

	var kv nats.KeyValue

	// create a Key/Value store with the bucket named "discovery" if doesn't exists
	// otherwise it creates a Key/Value store by existing bucket name
	if stream, _ := jsCtx.StreamInfo("KV_discovery"); stream == nil {

		// creating a Key/Value bucket with the name "discovery"
		kv, _ = jsCtx.CreateKeyValue(&nats.KeyValueConfig{
			Bucket: "discovery",
		})
	} else {
		kv, _ = jsCtx.KeyValue("discovery")
	}

	kv.Put("services.orders", []byte("https://localhost:8080/orders"))
	entry, _ := kv.Get("services.orders")
	fmt.Printf("%s @ %d -> %q\n", entry.Key(), entry.Revision(), string(entry.Value()))

	kv.Put("services.orders", []byte("https://localhost:8080/v1/orders"))
	entry, _ = kv.Get("services.orders")
	fmt.Printf("%s @ %d -> %q\n", entry.Key(), entry.Revision(), string(entry.Value()))

	kv.Update("services.orders", []byte("https://localhost:8080/v1/orders"), 1)
	entry, _ = kv.Get("services.orders")
	fmt.Printf("%s @ %d -> %q\n", entry.Key(), entry.Revision(), string(entry.Value()))

	name := <-jsCtx.StreamNames()
	fmt.Printf("KV stream name: %s\n", name)

	kv.Put("services.consumers", []byte("http://localhost:8080/v2/consumers"))
	entry, _ = kv.Get("services.consumers")
	fmt.Printf("%s @ %d -> %q\n", entry.Key(), entry.Revision(), string(entry.Value()))

	kv.Delete("services.consumers")
	entry, err = kv.Get("services.consumers")
	if err == nil {
		fmt.Printf("%s @ %d -> %q\n", entry.Key(), entry.Revision(), string(entry.Value()))
	}
}
