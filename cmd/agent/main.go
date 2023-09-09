package main

import (
	"context"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/agent/runner"

	"github.com/kontik-pk/yandex-metrics-scraper/internal/flags"
)

func main() {
	params := flags.Init(
		flags.WithConfig(),
		flags.WithPollInterval(),
		flags.WithReportInterval(),
		flags.WithAddr(),
		flags.WithKey(),
		flags.WithRateLimit(),
		flags.WithTLSKeyPath(),
		flags.WithGrpcAddr(),
	)
	ctx := context.Background()
	runner := runner.New(params)

	runner.Run(ctx)
}
