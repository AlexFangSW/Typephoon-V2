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

	"github.com/AlexFangSW/Typephoon-V2/subjects"
	"github.com/alecthomas/kong"
	kongyaml "github.com/alecthomas/kong-yaml"
	"github.com/bsm/redislock"
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
	LeaderLockTTL time.Duration `name:"leader_lock_ttl" yaml:"leader_lock_ttl" default:"5s" env:"LEADER_LOCK_TTL" help:"Leader lock TTL"`
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
		c.LeaderLockTTL,
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

type user struct {
	ID   string
	Name string
	Msg  *nats.Msg
}

type MatchmakingService struct {
	natsConn    *nats.Conn
	redisClient *redis.Client

	matchInterval time.Duration
	matchSize     int
	leaderLockTTL time.Duration

	joinMsg  chan *nats.Msg
	leaveMsg chan *nats.Msg

	pool      map[string]user
	timer     *time.Timer
	timerChan <-chan time.Time
}

func NewMatchmakingService(
	natsConn *nats.Conn,
	redisClient *redis.Client,
	matchInterval time.Duration,
	matchSize int,
	leaderLockTTL time.Duration,
) *MatchmakingService {
	return &MatchmakingService{
		natsConn:      natsConn,
		redisClient:   redisClient,
		matchInterval: matchInterval,
		matchSize:     matchSize,
		leaderLockTTL: leaderLockTTL,
		joinMsg:       make(chan *nats.Msg),
		leaveMsg:      make(chan *nats.Msg),
		pool:          make(map[string]user),
	}
}

func (ms *MatchmakingService) handleJoin(ctx context.Context, msg *nats.Msg) error {
	return nil
}

func (ms *MatchmakingService) handleLeave(ctx context.Context, msg *nats.Msg) error {
	return nil
}

func (ms *MatchmakingService) reset() {
	clear(ms.pool)
	ms.timer.Stop()
	ms.timerChan = nil
}

func (ms *MatchmakingService) provisionGame(ctx context.Context) {
}

// Matchmaking background worker
//
// Consume channel, group by players of N under X timeout.
// Provision the game when N players are found or when timeout
// happens.
// Timeout starts when the first player for the next group comes
// in.
func (ms *MatchmakingService) worker(ctx context.Context) {
	for {
		select {
		case msg := <-ms.joinMsg:
			if err := ms.handleJoin(ctx, msg); err != nil {
				// depending on the error, resp to current msg or all players
			}
		case msg := <-ms.leaveMsg:
			if err := ms.handleLeave(ctx, msg); err != nil {
				// depending on the error, resp to current msg or all players
			}
		case <-ms.timerChan:
			ms.provisionGame(ctx)
		case <-ctx.Done():
			return
		}
	}
}

// Try to aquire leader lock
// A total of two events will be sent in its lifetime
// - true: when we aquire the lock
// - false: when we lose the lock
func (ms *MatchmakingService) aquireLeader(ctx context.Context, resultChan chan<- bool) {
	locker := redislock.New(ms.redisClient)
	aquired := false
	var lock *redislock.Lock

	for {
		select {
		case <-ctx.Done():
			return

		default:
			time.Sleep(ms.leaderLockTTL / 2)
			slog.Debug("leader lock state", "aquired", aquired)

			if !aquired {
				var err error
				lock, err = locker.Obtain(ctx, leaderKey, ms.leaderLockTTL, &redislock.Options{})
				if err != nil {
					continue
				}
				defer lock.Release(ctx)
				aquired = true
				resultChan <- true
			}

			if err := lock.Refresh(ctx, ms.leaderLockTTL, &redislock.Options{}); err != nil {
				resultChan <- false
				return
			}
		}
	}
}

func (ms *MatchmakingService) Start(ctx context.Context) error {
	leaderChan := make(chan bool)
	go ms.aquireLeader(ctx, leaderChan)
	<-leaderChan

	subJoin, err := ms.natsConn.ChanSubscribe(string(subjects.MatchJoin), ms.joinMsg)
	if err != nil {
		return fmt.Errorf("chan subscribe match join: %w", err)
	}
	defer func(sub *nats.Subscription) {
		sub.Unsubscribe()
		sub.Drain()
	}(subJoin)

	subLeave, err := ms.natsConn.ChanSubscribe(string(subjects.MatchLeave), ms.leaveMsg)
	if err != nil {
		return fmt.Errorf("chan subscribe match join: %w", err)
	}
	defer func(sub *nats.Subscription) {
		sub.Unsubscribe()
		sub.Drain()
	}(subLeave)

	go ms.worker(ctx)

	for {
		select {
		case aquired := <-leaderChan:
			if !aquired {
				return errors.New("lost the leader lock")
			}
		case <-ctx.Done():
			return nil
		}
	}
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
