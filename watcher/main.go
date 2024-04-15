package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

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
		log.Fatalf("error creating JetStream: %v", err)
	}
	var kv nats.KeyValue

	if stream, _ := jsCtx.StreamInfo("KV_discovery"); stream == nil {
		kv, _ = jsCtx.CreateKeyValue(&nats.KeyValueConfig{
			Bucket: "discovery",
		})
	} else {
		kv, _ = jsCtx.KeyValue("discovery")
	}

	// key watcher for wildcard "services.*"
	w, _ := kv.Watch("services.*")
	defer w.Stop()

	for kve := range w.Updates() {
		if kve != nil {
			fmt.Printf("%s @ %d -> %q (op: %s)\n", kve.Key(), kve.Revision(), string(kve.Value()), kve.Operation())
		}
	}
	runtime.Goexit()
}
