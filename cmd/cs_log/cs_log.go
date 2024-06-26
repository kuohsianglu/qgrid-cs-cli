package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/chirpstack/chirpstack/api/go/v4/integration"
	"github.com/chirpstack/chirpstack/api/go/v4/stream"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	server string
	key    string
)

func request_log() {
	flag.StringVar(&server, "server", "localhost:6379", "Redis hostname:port")
	flag.StringVar(&key, "key", "api:stream:request", "Redis Streams key to read from")
	flag.Parse()

	rdb := redis.NewClient(&redis.Options{
		Addr: server,
	})
	ctx := context.Background()
	lastID := "0"

	for {
		resp, err := rdb.XRead(ctx, &redis.XReadArgs{
			Streams: []string{key, lastID},
			Count:   10,
			Block:   0,
		}).Result()
		if err != nil {
			log.Fatal(err)
		}

		if len(resp) != 1 {
			log.Fatal("Exactly one stream response is expected")
		}

		for _, msg := range resp[0].Messages {
			lastID = msg.ID

			if b, ok := msg.Values["request"].(string); ok {
				var pl stream.ApiRequestLog
				if err := proto.Unmarshal([]byte(b), &pl); err != nil {
					log.Fatal(err)
				}

				fmt.Println("=== Request ===")
				fmt.Println(protojson.Format(&pl))
				fmt.Println("===============")
			}

		}
	}
}

func meta_log() {
	flag.StringVar(&server, "server", "localhost:6379", "Redis hostname:port")
	flag.StringVar(&key, "key", "stream:meta", "Redis Streams key to read from")
	flag.Parse()

	rdb := redis.NewClient(&redis.Options{
		Addr: server,
	})
	ctx := context.Background()

	lastID := "0"

	for {
		resp, err := rdb.XRead(ctx, &redis.XReadArgs{
			Streams: []string{key, lastID},
			Count:   10,
			Block:   0,
		}).Result()
		if err != nil {
			log.Fatal(err)
		}

		if len(resp) != 1 {
			log.Fatal("Exactly one stream response is expected")
		}

		for _, msg := range resp[0].Messages {
			lastID = msg.ID

			if b, ok := msg.Values["up"].(string); ok {
				var pl stream.UplinkMeta
				if err := proto.Unmarshal([]byte(b), &pl); err != nil {
					log.Fatal(err)
				}

				fmt.Println("=== UP ===")
				fmt.Println(protojson.Format(&pl))
				fmt.Println("==========")
			}

			if b, ok := msg.Values["down"].(string); ok {
				var pl stream.DownlinkMeta
				if err := proto.Unmarshal([]byte(b), &pl); err != nil {
					log.Fatal(err)
				}

				fmt.Println("=== DOWN ===")
				fmt.Println(protojson.Format(&pl))
				fmt.Println("============")
			}
		}
	}
}

func event_log() {
	flag.StringVar(&server, "server", "localhost:6379", "Redis hostname:port")
	flag.StringVar(&key, "key", "device:stream:event", "Redis Streams key to read from")
	flag.Parse()

	rdb := redis.NewClient(&redis.Options{
		Addr: server,
	})
	ctx := context.Background()

	lastID := "0"

	for {
		resp, err := rdb.XRead(ctx, &redis.XReadArgs{
			Streams: []string{key, lastID},
			Count:   10,
			Block:   0,
		}).Result()
		if err != nil {
			log.Fatal(err)
		}

		if len(resp) != 1 {
			log.Fatal("Exactly one stream response is expected")
		}

		for _, msg := range resp[0].Messages {
			lastID = msg.ID

			if b, ok := msg.Values["up"].(string); ok {
				var pl integration.UplinkEvent
				if err := proto.Unmarshal([]byte(b), &pl); err != nil {
					log.Fatal(err)
				}

				fmt.Println("=== UP ===")
				fmt.Println(protojson.Format(&pl))
				fmt.Println("==========")
			}

			if b, ok := msg.Values["join"].(string); ok {
				var pl integration.JoinEvent
				if err := proto.Unmarshal([]byte(b), &pl); err != nil {
					log.Fatal(err)
				}

				fmt.Println("=== JOIN ===")
				fmt.Println(protojson.Format(&pl))
				fmt.Println("============")
			}

			if b, ok := msg.Values["ack"].(string); ok {
				var pl integration.AckEvent
				if err := proto.Unmarshal([]byte(b), &pl); err != nil {
					log.Fatal(err)
				}

				fmt.Println("=== ACK ===")
				fmt.Println(protojson.Format(&pl))
				fmt.Println("===========")
			}

			if b, ok := msg.Values["txack"].(string); ok {
				var pl integration.TxAckEvent
				if err := proto.Unmarshal([]byte(b), &pl); err != nil {
					log.Fatal(err)
				}

				fmt.Println("=== TX ACK ===")
				fmt.Println(protojson.Format(&pl))
				fmt.Println("==============")
			}

			if b, ok := msg.Values["log"].(string); ok {
				var pl integration.LogEvent
				if err := proto.Unmarshal([]byte(b), &pl); err != nil {
					log.Fatal(err)
				}

				fmt.Println("=== LOG ===")
				fmt.Println(protojson.Format(&pl))
				fmt.Println("===========")
			}

			if b, ok := msg.Values["status"].(string); ok {
				var pl integration.StatusEvent
				if err := proto.Unmarshal([]byte(b), &pl); err != nil {
					log.Fatal(err)
				}

				fmt.Println("=== STATUS ===")
				fmt.Println(protojson.Format(&pl))
				fmt.Println("==============")
			}

			if b, ok := msg.Values["location"].(string); ok {
				var pl integration.LocationEvent
				if err := proto.Unmarshal([]byte(b), &pl); err != nil {
					log.Fatal(err)
				}

				fmt.Println("=== LOCATION ===")
				fmt.Println(protojson.Format(&pl))
				fmt.Println("================")
			}

			if b, ok := msg.Values["integration"].(string); ok {
				var pl integration.IntegrationEvent
				if err := proto.Unmarshal([]byte(b), &pl); err != nil {
					log.Fatal(err)
				}

				fmt.Println("=== INTEGRATION ===")
				fmt.Println(protojson.Format(&pl))
				fmt.Println("===================")
			}
		}
	}
}

func frame_log() {
	flag.StringVar(&server, "server", "localhost:6379", "Redis hostname:port")
	flag.StringVar(&key, "key", "gw:stream:frame", "Redis Streams key to read from")
	flag.Parse()

	rdb := redis.NewClient(&redis.Options{
		Addr: server,
	})
	ctx := context.Background()

	lastID := "0"

	for {
		resp, err := rdb.XRead(ctx, &redis.XReadArgs{
			Streams: []string{key, lastID},
			Count:   10,
			Block:   0,
		}).Result()
		if err != nil {
			log.Fatal(err)
		}

		if len(resp) != 1 {
			log.Fatal("Exactly one stream response is expected")
		}

		for _, msg := range resp[0].Messages {
			lastID = msg.ID

			if b, ok := msg.Values["up"].(string); ok {
				var pl stream.UplinkFrameLog
				if err := proto.Unmarshal([]byte(b), &pl); err != nil {
					log.Fatal(err)
				}

				fmt.Println("=== UP ===")
				fmt.Println(protojson.Format(&pl))
				fmt.Println("==========")
			}

			if b, ok := msg.Values["down"].(string); ok {
				var pl stream.DownlinkFrameLog
				if err := proto.Unmarshal([]byte(b), &pl); err != nil {
					log.Fatal(err)
				}

				fmt.Println("=== DOWN ===")
				fmt.Println(protojson.Format(&pl))
				fmt.Println("============")
			}
		}
	}
}

func main() {
	log_type := os.Args[1]

	switch {
	case log_type == "event":
		event_log()
	case log_type == "frame":
		frame_log()
	case log_type == "meta":
		meta_log()
	case log_type == "request":
		request_log()
	}
}
