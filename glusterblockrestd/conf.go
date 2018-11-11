package main

import (
	"io/ioutil"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config maintains the glusterblockrestd configurations
type Config struct {
	GlusterBlockCLIPath string `toml:"gluster-block-cli-path"`
	Port                int    `toml:"port"`
	LogDir              string `toml:"logdir"`
	LogFile             string `toml:"logfile"`
	LogLevel            string `toml:"loglevel"`
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
