package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var urisRaw string
	var timeout time.Duration

	flag.StringVar(&urisRaw, "uris", "", "comma-separated mongodb uris")
	flag.DurationVar(&timeout, "timeout", 8*time.Second, "connect timeout")
	flag.Parse()

	if strings.TrimSpace(urisRaw) == "" {
		fmt.Fprintln(os.Stderr, "mongocheck error: -uris is required")
		os.Exit(1)
	}

	uris := splitURIs(urisRaw)
	failed := false

	for _, uri := range uris {
		if err := checkURI(uri, timeout); err != nil {
			failed = true
			fmt.Printf("FAIL %s -> %v\n", uri, err)
			continue
		}
		fmt.Printf("OK   %s\n", uri)
	}

	if failed {
		os.Exit(1)
	}
}

func splitURIs(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

func checkURI(uri string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer client.Disconnect(context.Background())

	pingCtx, pingCancel := context.WithTimeout(context.Background(), timeout)
	defer pingCancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		return fmt.Errorf("ping: %w", err)
	}
	return nil
}
