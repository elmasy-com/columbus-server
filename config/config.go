package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type conf struct {
	MongoURI       string   `yaml:"MongoURI"`
	Address        string   `yaml:"Address"`
	TrustedProxies []string `yaml:"TrustedProxies"`
	SSLCert        string   `yaml:"SSLCert"`
	SSLKey         string   `yaml:"SSLKey"`
	LogErrorOnly   bool     `yaml:"LogErrorOnly"`
	EnableStatAPI  bool     `yaml:"EnableStatAPI"`
	StatAPIWait    int      `yaml:"StatAPIWait"`
}

var (
	MongoURI       string   // MongoDB connection string
	Address        string   // Address to listen on
	TrustedProxies []string // A list of trusted proxies
	SSLCert        string
	SSLKey         string
	LogErrorOnly   bool
	EnableStatAPI  bool
	StatAPIWait    int
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

	SSLCert = c.SSLCert
	SSLKey = c.SSLKey

	LogErrorOnly = c.LogErrorOnly
	EnableStatAPI = c.EnableStatAPI

	if c.StatAPIWait < 0 {
		return fmt.Errorf("StatAPIWait is negative")
	}
	if c.StatAPIWait == 0 {
		c.StatAPIWait = 1440
	}

	StatAPIWait = c.StatAPIWait

	return nil
}
