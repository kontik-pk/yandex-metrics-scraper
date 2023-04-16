package flags

import (
	"flag"
	"os"
	"strconv"
)

const (
	defaultAddr           string = "localhost:8080"
	defaultReportInterval int    = 10
	defaultPollInterval   int    = 2
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
		flag.IntVar(&p.ReportInterval, "r", defaultReportInterval, "report interval")
		if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
			reportIntervalEnv, err := strconv.Atoi(envReportInterval)
			if err == nil {
				p.ReportInterval = reportIntervalEnv
			}
		}
	}
}

func WithPollInterval() Option {
	return func(p *params) {
		flag.IntVar(&p.PollInterval, "p", defaultPollInterval, "poll interval")
		if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
			pollIntervalEnv, err := strconv.Atoi(envPollInterval)
			if err == nil {
				p.PollInterval = pollIntervalEnv
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
	ReportInterval int
	PollInterval   int
}
