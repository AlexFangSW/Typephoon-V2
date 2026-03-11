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
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	kongyaml "github.com/alecthomas/kong-yaml"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DumpConfig bool            `name:"dump_config" yaml:"-" help:"Dump current config as YAML and exit"`
	ConfigFile kong.ConfigFlag `name:"config_file" yaml:"-" default:"./matchmaker.yaml" help:"Path to YAML config file"`

	LogLevel      string `name:"log_level" yaml:"log_level" enum:"DEBUG,INFO,WARNNING,ERROR" default:"INFO" env:"LOG_LEVEL" help:"Set log level"`
	RedisUrl      string `name:"redis_url" yaml:"redis_url" default:"localhost:6379" env:"REDIS_URL" help:"Redis url"`
	MatchInterval int    `name:"match_interval" yaml:"match_interval" default:"5" env:"MATCH_INTERVAL" help:"Sleep interval before next round of matchmaking"`
	MatchSize     int    `name:"match_size" yaml:"match_size" default:"5" env:"MATCH_SIZE" help:"Max players per game"`
}

func (c *Config) Run() error {
	if c.DumpConfig {
		out, err := yaml.Marshal(c)
		if err != nil {
			return err
		}
		fmt.Fprint(os.Stdout, string(out))
		return nil
	}

	return nil
}

func main() {
	var cfg Config
	ctx := kong.Parse(&cfg,
		kong.Description("Matchmaking service for Typephoon project"),
		kong.Configuration(kongyaml.Loader, "./matchmaker.yaml"))
	ctx.FatalIfErrorf(cfg.Run())
}
