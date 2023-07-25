package flags

import (
	"flag"
	"os"
	"strconv"
)

const (
	defaultAddr            string = "127.0.0.1:8080"
	defaultReportInterval  int    = 5
	defaultPollInterval    int    = 1
	defaultStoreInterval   int    = 15
	defaultFileStoragePath string = "/tmp/short-url-db.json"
	defaultRestore         bool   = true
)

type Option func(params *Params)

func WithRateLimit() Option {
	return func(p *Params) {
		flag.IntVar(&p.RateLimit, "l", 1, "max requests to send on server")
		if envKey := os.Getenv("RATE_LIMIT"); envKey != "" {
			p.Key = envKey
		}
	}
}

func WithKey() Option {
	return func(p *Params) {
		flag.StringVar(&p.Key, "k", "", "key for using hash subscription")
		if envKey := os.Getenv("KEY"); envKey != "" {
			p.Key = envKey
		}
	}
}

func WithDatabase() Option {
	return func(p *Params) {
		result := ""
		flag.StringVar(&result, "d", "", "connection string for db")
		if envDBAddr := os.Getenv("DATABASE_DSN"); envDBAddr != "" {
			result = envDBAddr
		}
		p.DatabaseAddress = result
	}
}

func WithAddr() Option {
	return func(p *Params) {
		flag.StringVar(&p.FlagRunAddr, "a", defaultAddr, "address and port to run server")
		if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
			p.FlagRunAddr = envRunAddr
		}
	}
}

func WithReportInterval() Option {
	return func(p *Params) {
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
	return func(p *Params) {
		flag.IntVar(&p.PollInterval, "p", defaultPollInterval, "poll interval")
		if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
			pollIntervalEnv, err := strconv.Atoi(envPollInterval)
			if err == nil {
				p.PollInterval = pollIntervalEnv
			}
		}
	}
}

func WithStoreInterval() Option {
	return func(p *Params) {
		flag.IntVar(&p.StoreInterval, "i", defaultStoreInterval, "store interval in seconds")
		if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
			storeIntervalEnv, err := strconv.Atoi(envStoreInterval)
			if err == nil {
				p.StoreInterval = storeIntervalEnv
			}
		}
	}
}

func WithFileStoragePath() Option {
	return func(p *Params) {
		flag.StringVar(&p.FileStoragePath, "f", defaultFileStoragePath, "file name for metrics collection")
		if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
			fileStoragePath, err := strconv.Atoi(envFileStoragePath)
			if err == nil {
				p.StoreInterval = fileStoragePath
			}
		}
	}
}

func WithRestore() Option {
	return func(p *Params) {
		flag.BoolVar(&p.Restore, "r", defaultRestore, "restore data from file")
		if envRestore := os.Getenv("RESTORE"); envRestore != "" {
			restore, err := strconv.Atoi(envRestore)
			if err == nil {
				p.StoreInterval = restore
			}
		}
	}
}

func Init(opts ...Option) *Params {
	p := &Params{}
	for _, opt := range opts {
		opt(p)
	}
	flag.Parse()
	return p
}

type Params struct {
	FlagRunAddr     string
	DatabaseAddress string
	ReportInterval  int
	PollInterval    int
	StoreInterval   int
	FileStoragePath string
	Restore         bool
	Key             string
	RateLimit       int
}
