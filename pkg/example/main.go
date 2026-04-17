// Package main is a runnable example of the golog v2 API.
//
// Run with:
//
//	go run github.com/marcosstupnicki/go-log/pkg/example
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	golog "github.com/marcosstupnicki/go-log"
)

func main() {
	// 1. Create a logger. `env` is required and becomes a sticky field.
	logger, err := golog.New("local", golog.WithLevel(golog.DebugLevel))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create logger: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()

	// 2. Per-call fields via the universal Field() builder.
	logger.Info(ctx, "service started",
		golog.Field("version", "2.0.0"),
		golog.Field("port", 8080),
		golog.Field("startup", 120*time.Millisecond),
	)

	// 3. Enrich the context once (e.g. in middleware) — every subsequent
	// log call on this context includes those fields automatically.
	ctx = golog.Enrich(ctx,
		golog.Field("request_id", "req-abc-123"),
		golog.Field("path", "/api/users"),
	)

	logger.Info(ctx, "handler invoked",
		golog.Field("user_id", "u-42"),
	)

	// 4. Errors: the key you pass is ignored — zap always uses "error".
	err = errors.New("database connection refused")
	logger.Error(ctx, "dependency failure",
		golog.Field("whatever-key", err),
		golog.Field("attempt", 3),
	)
}
