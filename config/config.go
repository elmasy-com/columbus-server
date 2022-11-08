package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type conf struct {
	MongoURI       string   `yaml:"MongoURI"`
	Address        string   `yaml:"Address"`
	TrustedProxies []string `yaml:"TrustedProxies"`
	BLacklistSec   int      `yaml:"BlacklistSec"`
	SSLCert        string   `yaml:"SSLCert"`
	SSLKey         string   `yaml:"SSLKey"`
	LogErrorOnly   bool     `yaml:"LogErrorOnly"`
}

var (
	MongoURI       string   // MongoDB connection string
	Address        string   // Address to listen on
	TrustedProxies []string // A list of trusted proxies
	BlacklistTime  time.Duration
	SSLCert        string
	SSLKey         string
	LogErrorOnly   bool
)

// Parse parses the config file in path and gill the global variables.
func Parse(path string) error {

	c := conf{}

	out, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %s", path, err)
	}

	err = yaml.Unmarshal(out, &c)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %s", err)
	}

	if c.MongoURI == "" {
		return fmt.Errorf("MongoURI is empty")
	}
	MongoURI = c.MongoURI

	if c.Address == "" {
		c.Address = ":8080"
	}
	Address = c.Address

	TrustedProxies = c.TrustedProxies

	if c.BLacklistSec == 0 {
		c.BLacklistSec = 60
	}

	BlacklistTime = time.Duration(c.BLacklistSec) * time.Second

	SSLCert = c.SSLCert
	SSLKey = c.SSLKey

	LogErrorOnly = c.LogErrorOnly

	return nil
}
