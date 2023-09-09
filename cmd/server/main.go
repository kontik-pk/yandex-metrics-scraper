package main

import (
	"context"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/server/runner"

	"github.com/kontik-pk/yandex-metrics-scraper/internal/flags"
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
		flags.WithTrustedSubnet(),
		flags.WithGrpc(),
		flags.WithGrpcAddr(),
	)
	ctx := context.Background()
	serverRunner := runner.New(params)

	serverRunner.Run(ctx)
}
