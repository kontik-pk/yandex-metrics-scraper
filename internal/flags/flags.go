package flags

import (
	"flag"
	"os"
	"strconv"
	"time"
)

const (
	defaultAddr           string = "localhost:8080"
	defaultReportInterval        = 2 * time.Second
	defaultPollInterval          = 1 * time.Second
)

type Option func(params2 *params)

func WithAddr() Option {
	return func(p *params) {
		flag.StringVar(&p.FlagRunAddr, "a", defaultAddr, "address and port to run server")
		if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
			p.FlagRunAddr = envRunAddr
		}
	}
}

func WithReportInterval() Option {
	return func(p *params) {
		flag.DurationVar(&p.ReportInterval, "r", defaultReportInterval, "report interval")
		if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
			reportIntervalEnv, err := strconv.Atoi(envReportInterval)
			if err == nil {
				p.ReportInterval = time.Duration(reportIntervalEnv) * time.Second
			}
		}
	}
}

func WithPollInterval() Option {
	return func(p *params) {
		flag.DurationVar(&p.PollInterval, "p", defaultPollInterval, "poll interval")
		if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
			pollIntervalEnv, err := strconv.Atoi(envPollInterval)
			if err == nil {
				p.PollInterval = time.Duration(pollIntervalEnv) * time.Second
			}
		}
	}
}

func Init(opts ...Option) *params {
	p := &params{}
	for _, opt := range opts {
		opt(p)
	}
	flag.Parse()
	return p
}

type params struct {
	FlagRunAddr    string
	ReportInterval time.Duration
	PollInterval   time.Duration
}
