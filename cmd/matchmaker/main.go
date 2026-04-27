package main

import (
	"context"
	"encoding/json"
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

type MatchmakingService struct {
	natsConn    *nats.Conn
	redisClient *redis.Client

	matchInterval time.Duration
	matchSize     int
	leaderLockTTL time.Duration

	subject types.Subject
	events  chan *nats.Msg
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
		subject:       types.Subjects.MatchJoin,
		events:        make(chan *nats.Msg),
	}
}

type playerInfo struct {
	msg    *nats.Msg
	userID string
}

// Matchmaking background worker
//
// Consume channel, group by players of N under X timeout.
// Provision the game when N players are found or when timeout
// happens.
// Timeout starts when the first player for the next group comes
// in.
func (ms *MatchmakingService) worker(ctx context.Context) {
	var playerPool = make(map[string]playerInfo)
	var timer *time.Timer
	var timerChan <-chan time.Time // Select will ignore nil channel

	for {
		select {
		case msg := <-ms.events:
			// Parse message
			var msgBody types.EventEnvelop
			if err := json.Unmarshal(msg.Data, &msgBody); err != nil {
				slog.Warn("bad message", "msg", msg.Data, "header", msg.Header, "error", err)

				// Reply with error
				payload, jsonErr := json.Marshal(types.ErrorPayload{
					Error: fmt.Sprintf("failed to unmarshal payload, error: %s", err),
				})
				if jsonErr != nil {
					slog.Warn("failed to marshal error msg payload")
					continue
				}

				data, jsonErr := json.Marshal(types.EventEnvelop{
					Type:    types.Events.Error,
					Payload: payload,
				})
				if jsonErr != nil {
					slog.Warn("failed to marshal error msg")
					continue
				}

				ms.natsConn.Publish(msg.Reply, data)
				continue
			}
			userID := msg.Header.Get(string(types.Headers.UserID))
			if userID == "" {
				slog.Warn("missing user ID")

				// Reply with error
				payload, jsonErr := json.Marshal(types.ErrorPayload{
					Error: "missing user ID",
				})
				if jsonErr != nil {
					slog.Warn("failed to marshal error msg payload")
					continue
				}

				data, jsonErr := json.Marshal(types.EventEnvelop{
					Type:    types.Events.Error,
					Payload: payload,
				})
				if jsonErr != nil {
					slog.Warn("failed to marshal error msg")
					continue
				}

				ms.natsConn.Publish(msg.Reply, data)
				continue
			}

			switch msgBody.Type {
			case types.Events.MatchJoin:
				if timer == nil {
					timer = time.NewTimer(ms.matchInterval)
					timerChan = timer.C
				}

				player := playerInfo{
					msg:    msg,
					userID: msg.Header.Get(string(types.Headers.UserID)),
				}
				playerPool[player.userID] = player

				// Check current size
				if len(playerPool) < ms.matchSize {
					continue
				}
				// TODO: Provision a game
				// TODO: Reply

				// Reset
				clear(playerPool)
				timer.Stop()
				timerChan = nil
			case types.Events.MatchLeave:
				delete(playerPool, userID)
				if len(playerPool) == 0 {
					// Reset
					timer.Stop()
					timerChan = nil
				}
			}

		case <-timerChan:
			// TODO: Provision a game
			// TODO: Reply

			// Reset
			timer.Stop()
			timerChan = nil
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

	sub, err := ms.natsConn.Subscribe(string(ms.subject), func(msg *nats.Msg) {
		ms.events <- msg
	})
	if err != nil {
		return fmt.Errorf("nats subscribe: %w", err)
	}
	go ms.worker(ctx)

	defer func(sub *nats.Subscription) {
		sub.Unsubscribe()
		sub.Drain()
	}(sub)

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
