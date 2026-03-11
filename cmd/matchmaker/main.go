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
	"github.com/alecthomas/kong"
)

type Config struct {
	LogLevel      string `enum:"DEBUG,INFO,WARNNING,ERROR" default:"INFO" help:"Set log level"`
	RedisUrl      string `default:"localhost:6379" help:"Redis url"`
	MatchInterval int    `help:""`
	MatchSize     int
}

func main() {
	var cfg Config
	kong.Parse(&cfg)
}
