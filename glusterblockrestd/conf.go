package main

import (
	"errors"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/BurntSushi/toml"
)

// Config maintains the glusterblockrestd configurations
type Config struct {
	GlusterBlockCLIPath string `toml:"gluster-block-cli-path"`
	Addr                string `toml:"address"`
	LogDir              string `toml:"logdir"`
	LogFile             string `toml:"logfile"`
	LogLevel            string `toml:"loglevel"`
	CertFile            string `toml:"certfile"`
	KeyFile             string `toml:"keyfile"`
	User                string `toml:"user"`
	Secret              string `toml:"secret"`
	AuthEnabled         bool   `toml:"restauth"`
	EnableTLS           bool   `toml:"enabletls"`
}

func loadConfig(confFilePath string) (*Config, error) {
	var conf Config
	b, err := ioutil.ReadFile(filepath.Clean(confFilePath))
	if err != nil {
		return nil, err
	}
	if _, err := toml.Decode(string(b), &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

func validateAddress(addr string) error {
	shost, sport, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}

	port, err := strconv.Atoi(sport)
	if err != nil {
		return err
	}

	err = errors.New("invalid address for glusterblockrestd")
	if port < 0 || port > 65353 {
		return err
	}
	if shost != "" && net.ParseIP(shost) == nil {
		return err
	}
	return nil
}
func setWithEnvVariables(options *Config) {
	// Check for user key
	env := os.Getenv("GB_REST_USER_KEY")
	if "" != env {
		options.AuthEnabled = true
		options.User = env
	}

	// Check for secret
	env = os.Getenv("GB_REST_SECRET_KEY")
	if "" != env {
		options.Secret = env
	}
}
