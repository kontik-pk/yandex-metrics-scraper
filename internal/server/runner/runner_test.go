package runner

import (
	"context"
	"fmt"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/agent/collector"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/kontik-pk/yandex-metrics-scraper/internal/flags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestRunner_Run(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		mockedSaver := newMockSaver(t)
		mockedSaver.On("Restore", mock.Anything).Return([]collector.StoredMetric{}, nil)
		mockedSaver.On("Save", mock.Anything, mock.AnythingOfType("[]collector.StoredMetric")).Return(nil)

		mockedAppServer := newMockHttpServer(t)
		mockedAppServer.On("ListenAndServe").Return(nil)
		mockedPprofServer := newMockHttpServer(t)
		mockedPprofServer.On("ListenAndServe").Return(nil)

		logger, err := zap.NewDevelopment()
		if err != nil {
			fmt.Println("error while creating logger, exit")
			return
		}
		defer logger.Sync()
		log := logger.Sugar()
		r := Runner{
			saver:           mockedSaver,
			metricsInterval: 1,
			isRestore:       true,
			storeInterval:   1,
			tlsKey:          "",
			appSrv:          mockedAppServer,
			pprofSrv:        mockedPprofServer,
			logger:          log,
		}
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		ctx, cancel := context.WithTimeout(ctx, time.Second*3)
		defer cancel()
		go r.Run(ctx)
		<-ctx.Done()
	})
	t.Run("positive: signals", func(t *testing.T) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		mockedSaver := newMockSaver(t)
		mockedSaver.On("Restore", mock.Anything).Return([]collector.StoredMetric{}, nil)
		mockedSaver.On("Save", mock.Anything, mock.AnythingOfType("[]collector.StoredMetric")).Return(nil)

		mockedAppServer := newMockHttpServer(t)
		mockedAppServer.On("ListenAndServe").Return(nil)
		mockedAppServer.On("Shutdown", ctx).Return(nil).Maybe()
		mockedPprofServer := newMockHttpServer(t)
		mockedPprofServer.On("ListenAndServe").Return(nil)

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

		logger, err := zap.NewDevelopment()
		if err != nil {
			fmt.Println("error while creating logger, exit")
			return
		}
		defer logger.Sync()
		r := Runner{
			saver:           mockedSaver,
			metricsInterval: 1,
			isRestore:       true,
			storeInterval:   1,
			tlsKey:          "",
			appSrv:          mockedAppServer,
			pprofSrv:        mockedPprofServer,
			signals:         sigs,
			logger:          logger.Sugar(),
		}
		go r.Run(ctx)
		time.Sleep(3 * time.Second)
		r.signals <- syscall.SIGTERM
	})
	t.Run("positive: with grpc", func(t *testing.T) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		mockedSaver := newMockSaver(t)
		mockedSaver.On("Restore", mock.Anything).Return([]collector.StoredMetric{}, nil)
		mockedSaver.On("Save", mock.Anything, mock.AnythingOfType("[]collector.StoredMetric")).Return(nil)

		mockedAppServer := newMockHttpServer(t)
		mockedAppServer.On("ListenAndServe").Return(nil)
		mockedAppServer.On("Shutdown", ctx).Return(nil).Maybe()

		mockedPprofServer := newMockHttpServer(t)
		mockedPprofServer.On("ListenAndServe").Return(nil)

		mockedGrpc := newMockGrpcServer(t)
		mockedListener := newMockListener(t)
		mockedListener.On("Close").Return(nil)
		mockedGrpc.On("Serve", mockedListener).Return(nil)
		mockedGrpc.On("GracefulStop").Return(nil)

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

		logger, err := zap.NewDevelopment()
		if err != nil {
			fmt.Println("error while creating logger, exit")
			return
		}
		defer logger.Sync()
		r := Runner{
			saver:           mockedSaver,
			metricsInterval: 1,
			isRestore:       true,
			storeInterval:   1,
			tlsKey:          "",
			appSrv:          mockedAppServer,
			pprofSrv:        mockedPprofServer,
			signals:         sigs,
			logger:          logger.Sugar(),
			grpcServer:      mockedGrpc,
			listener:        mockedListener,
		}
		go r.Run(ctx)
		time.Sleep(3 * time.Second)
		r.signals <- syscall.SIGTERM
	})
}

func TestNew(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		params := flags.Params{
			FlagRunAddr:     "http://127.0.0.1:8080",
			DatabaseAddress: "",
			ReportInterval:  1,
			PollInterval:    1,
			StoreInterval:   1,
			FileStoragePath: "/tmp/file.json",
			Restore:         true,
			Key:             "key",
			RateLimit:       10,
			CryptoKeyPath:   "",
		}
		r := New(&params)
		s, _ := initSaver(&params)
		expected := Runner{
			saver:           s,
			metricsInterval: 1,
			isRestore:       true,
			storeInterval:   1,
			tlsKey:          "",
		}
		assert.Equal(t, r.saver, expected.saver)
		assert.Equal(t, r.metricsInterval, expected.metricsInterval)
		assert.Equal(t, r.isRestore, expected.isRestore)
		assert.Equal(t, r.tlsKey, expected.tlsKey)
		assert.Equal(t, r.storeInterval, expected.storeInterval)
	})
}
