package main

import (
	"context"
	"fmt"

	"github.com/kontik-pk/yandex-metrics-scraper/internal/flags"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/runner/server"
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
	)
	fmt.Println(params)
	ctx := context.Background()
	serverRunner := server.New(params)

	serverRunner.Run(ctx)
}
