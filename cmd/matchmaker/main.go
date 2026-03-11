/*
Matchmaker works by providing two API, one for "queue-in" and another to poll
for the status.

The service is designed with "Active-Passive" architecture in mind,
we use an external cache as distributed lock, by acquiring the lock
the matchmaker service becomes "active".

Matchmaking logic:
When client "queue-in", they are added to a pool.
Every round, the matchmaker will group clients in the pool and update the status.
Each client will recive a "gameID" and a unique "token" that is than used to connect to the
game server.
*/
package main

import (
	"errors"
	"fmt"

	"github.com/alecthomas/kong"
)

type Config struct {
	LogLevel      string `enum:"DEBUG,INFO,WARNNING,ERROR" default:"INFO" env:"LOG_LEVEL" help:"Set log level"`
	RedisUrl      string `default:"localhost:6379" env:"REDIS_URL" help:"Redis url"`
	MatchInterval int    `default:"5" env:"MATCH_INTERVAL" help:"Sleep interval before next round of matchmaking"`
	MatchSize     int    `default:"5" env:"MATCH_SIZE" help:"Max players per game"`
}

func (c *Config) Run(ctx *Config) error {
	fmt.Printf("Current Config: %+v\n", ctx)
	return errors.New("EEEE")
}

func main() {
	var cfg Config
	ctx := kong.Parse(&cfg, kong.Description("Matchmaking service for Typephoon project"))
	ctx.FatalIfErrorf(ctx.Run())
}
