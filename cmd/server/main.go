package main

import (
	"context"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/flags"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/runner"
)

func main() {
	params := flags.Init(
		flags.WithConfig(),
		flags.WithAddr(),
		flags.WithStoreInterval(),
		flags.WithFileStoragePath(),
		flags.WithRestore(),
		flags.WithDatabase(),
		flags.WithKey(),
		flags.WithTLSKeyPath(),
	)

	ctx := context.Background()
	serverRunner := runner.New(params)
	serverRunner.Run(ctx)
}
