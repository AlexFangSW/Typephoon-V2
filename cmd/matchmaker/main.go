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

	"github.com/AlexFangSW/Typephoon-V2/types"
	"github.com/alecthomas/kong"
	kongyaml "github.com/alecthomas/kong-yaml"
	nats "github.com/nats-io/nats.go"
	redis "github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
)

var (
	leaderKey = "matchmaker:leader"
)

type Config struct {
	DumpConfig bool            `name:"dump_config" yaml:"-" help:"Dump current config as YAML and exit"`
	ConfigFile kong.ConfigFlag `name:"config_file" yaml:"-" placeholder:"./matchmaker.yaml" help:"Path to YAML config file"`

	LogLevel      string        `name:"log_level" yaml:"log_level" enum:"DEBUG,INFO,WARNNING,ERROR" default:"INFO" env:"LOG_LEVEL" help:"Set log level"`
	NATSURL       string        `name:"nats_url" yaml:"nats_url" default:"nats://127.0.0.1:4222" env:"NATS_URL" help:"NATS URL"`
	RedisURL      string        `name:"redis_url" yaml:"redis_url" default:"localhost:6379" env:"REDIS_URL" help:"Redis URL"`
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

	// Connections and clients
	natsConn, err := nats.Connect(c.NATSURL)
	if err != nil {
		return fmt.Errorf("nats connect: %w", err)
	}
	defer natsConn.Drain()

	redisClient := redis.NewClient(&redis.Options{
		Addr: c.RedisURL,
	})
	defer redisClient.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start main workflow
	service := NewMatchmakingService(
		natsConn,
		redisClient,
		c.MatchInterval,
		c.MatchSize,
	)

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
	natsConn    *nats.Conn
	redisClient *redis.Client

	matchInterval time.Duration
	matchSize     int

	events chan types.EventEnvelop
}

func NewMatchmakingService(
	natsConn *nats.Conn,
	redisClient *redis.Client,
	matchInterval time.Duration,
	matchSize int,
) *MatchmakingService {
	return &MatchmakingService{
		natsConn:      natsConn,
		redisClient:   redisClient,
		matchInterval: matchInterval,
		matchSize:     matchSize,
		events:        make(chan types.EventEnvelop),
	}
}

// Matchmaking background worker
//
// Consume channel, group by players of N under X timeout.
// Provision the game when N players are found or when timeout
// happens.
// Timeout starts when the first player for the next group comes
// in.
func (ms *MatchmakingService) worker() {
	// TODO: We consume events and keep a map of players and adjust according to the recived
	// event type (join, leave)
}

func (ms *MatchmakingService) Start(ctx context.Context) error {
	// TODO:
	// - Subscribe to subject: `match.join`
	// - Try to obtain leader lock, only start after lock is aquired,
	// keep refreshing lock, if lock is lost, return error

	return nil
}

func (ms *MatchmakingService) Stop(ctx context.Context) error {
	// TODO
	return nil
}

func main() {
	var cfg Config
	ctx := kong.Parse(&cfg,
		kong.Description("Matchmaking service for Typephoon project"),
		kong.Configuration(kongyaml.Loader, "./matchmaker.yaml"))
	ctx.FatalIfErrorf(ctx.Run()) // This is basically Config.Run()
}
