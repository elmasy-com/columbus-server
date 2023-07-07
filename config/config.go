package config

import (
	"fmt"
	"os"
	"time"

	"github.com/elmasy-com/elnet/dns"
	"gopkg.in/yaml.v3"
)

type conf struct {
	MongoURI       string   `yaml:"MongoURI"`
	Address        string   `yaml:"Address"`
	TrustedProxies []string `yaml:"TrustedProxies"`
	SSLCert        string   `yaml:"SSLCert"`
	SSLKey         string   `yaml:"SSLKey"`
	LogErrorOnly   bool     `yaml:"LogErrorOnly"`
	DNSServers     []string `yaml:"DNSServers"`
	DNSPort        string   `yaml:"DNSPort"`
	DNSProtocol    string   `yaml:"DNSProtocol"`
	DNSWorker      int      `yaml:"DNSWorker"`
}

var (
	MongoURI       string   // MongoDB connection string
	Address        string   // Address to listen on
	TrustedProxies []string // A list of trusted proxies
	SSLCert        string
	SSLKey         string
	LogErrorOnly   bool
	DNSServers     []string
	DNSPort        string
	DNSProtocol    string
	DNSWorker      int
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

	if c.DNSPort == "" {
		c.DNSPort = "53"
	}

	if len(c.DNSServers) > 0 {
		dns.UpdateConf(c.DNSServers, c.DNSPort)
	}

	DNSServers = c.DNSServers
	DNSPort = c.DNSPort

	if c.DNSProtocol == "" {
		c.DNSProtocol = "udp"
	}

	err = dns.UpdateClient(c.DNSProtocol, 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed to update DNS client: %w", err)
	}

	DNSProtocol = c.DNSProtocol

	if c.DNSWorker == 0 {
		c.DNSWorker = 1
	}

	DNSWorker = c.DNSWorker

	return nil
}
