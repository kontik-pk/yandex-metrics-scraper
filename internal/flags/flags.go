package flags

import (
	"flag"
	"os"
	"strconv"
)

const (
	defaultAddr            string = "localhost:8080"
	defaultReportInterval  int    = 10
	defaultPollInterval    int    = 2
	defaultStoreInterval   int    = 300
	defaultFileStoragePath string = "/tmp/metrics-db.json"
	defaultRestore         bool   = true
)

type Option func(params *Params)

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
	ReportInterval  int
	PollInterval    int
	StoreInterval   int
	FileStoragePath string
	Restore         bool
}
