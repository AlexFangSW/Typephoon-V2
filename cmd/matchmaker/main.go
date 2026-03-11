/*
Matchmaker works by providing two API, one for "queue-in" and another to poll
for the status.

The service is designed with "Active-Passive" architecture in mind,
we use an external cache as distributed lock, by acquiring the lock
the matchmaker service becomes "active".

Matchmaking logic:
When client "queue-in", they are added to a pool.
Every round, the matchmaker will group clients in the pool and update the status.
Each client will recive a "gameID"" that is used to connect to the game server.
*/
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	kongyaml "github.com/alecthomas/kong-yaml"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DumpConfig bool            `name:"dump_config" yaml:"-" help:"Dump current config as YAML and exit"`
	ConfigFile kong.ConfigFlag `name:"config_file" yaml:"-" placeholder:"./matchmaker.yaml" help:"Path to YAML config file"`

	LogLevel      string        `name:"log_level" yaml:"log_level" enum:"DEBUG,INFO,WARNNING,ERROR" default:"INFO" env:"LOG_LEVEL" help:"Set log level"`
	ListenAddress string        `name:"listen_address" yaml:"listen_address" default:":8080" env:"LISTEN_ADDRESS" help:"Server listen address"`
	RedisUrl      string        `name:"redis_url" yaml:"redis_url" default:"localhost:6379" env:"REDIS_URL" help:"Redis url"`
	MatchInterval time.Duration `name:"match_interval" yaml:"match_interval" default:"5s" env:"MATCH_INTERVAL" help:"Sleep interval before next round of matchmaking"`
	MatchSize     int           `name:"match_size" yaml:"match_size" default:"5" env:"MATCH_SIZE" help:"Max players per game"`
}

func (c *Config) Run() error {
	// Print out current config
	if c.DumpConfig {
		out, err := yaml.Marshal(c)
		if err != nil {
			return err
		}
		fmt.Println(string(out))
		return nil
	}

	// Set logger
	logLevel := slog.LevelInfo
	switch c.LogLevel {
	case slog.LevelDebug.String():
		logLevel = slog.LevelDebug
	case slog.LevelInfo.String():
		logLevel = slog.LevelInfo
	case slog.LevelWarn.String():
		logLevel = slog.LevelWarn
	case slog.LevelError.String():
		logLevel = slog.LevelError
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	}))
	slog.SetDefault(logger)

	// Start main workflow
	service := NewMatchmakingService(
		c.ListenAddress,
		c.RedisUrl,
		c.MatchInterval,
		c.MatchSize,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	serviceErr := make(chan error, 1)
	go func() {
		if err := service.Start(ctx); !errors.Is(err, http.ErrServerClosed) {
			serviceErr <- err
		}
	}()

	// Handle shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	select {
	case err := <-serviceErr:
		return err
	case <-stop:
		tctx, tcancel := context.WithTimeout(ctx, 30*time.Second)
		defer tcancel()
		return service.Stop(tctx)
	}
}

type MatchmakingService struct {
	server *http.Server

	redisUrl      string
	matchInterval time.Duration
	matchSize     int
}

func NewMatchmakingService(
	listenAddr string,
	redisUrl string,
	matchInterval time.Duration,
	matchSize int,
) *MatchmakingService {
	ms := &MatchmakingService{
		redisUrl:      redisUrl,
		matchInterval: matchInterval,
		matchSize:     matchSize,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/matchmaker/queue-int", LogMiddleware(ms.queueIn))
	mux.HandleFunc("GET /api/v1/matchmaker/status", LogMiddleware(ms.status))
	mux.HandleFunc("GET /health", LogMiddleware(ms.health))

	server := &http.Server{
		Addr:    listenAddr,
		Handler: mux,
	}
	ms.server = server

	return ms
}

// Add client to wait pool, hmm... how do we implement this pool ???
func (ms *MatchmakingService) queueIn(w http.ResponseWriter, r *http.Request) {

}

// Check client status
func (ms *MatchmakingService) status(w http.ResponseWriter, r *http.Request) {

}

func (ms *MatchmakingService) health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

// Matchmaking background worker
func (ms *MatchmakingService) worker() {

}

func (ms *MatchmakingService) Start(ctx context.Context) error {
	slog.Info("Start matchmaking service", "listen", ms.server.Addr)
	// TODO: Start background worker
	return ms.server.ListenAndServe()
}

func (ms *MatchmakingService) Stop(ctx context.Context) error {
	return ms.server.Shutdown(ctx)
}

// Custom response writer that has extra info
type respWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *respWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request", "method", r.Method, "path", r.URL.Path, "headers", r.Header)
		start := time.Now()

		rw := &respWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next(rw, r)

		totalTime := time.Since(start)
		slog.Info("Response", "method", r.Method, "path", r.URL.Path, "headers", r.Header, "status", rw.statusCode, "duration", totalTime)
	}
}

func main() {
	var cfg Config
	ctx := kong.Parse(&cfg,
		kong.Description("Matchmaking service for Typephoon project"),
		kong.Configuration(kongyaml.Loader, "./matchmaker.yaml"))
	ctx.FatalIfErrorf(ctx.Run()) // This is basically Config.Run()
}
