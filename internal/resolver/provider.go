package resolver

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"paulmatvienko/whois/internal/domain"
)

// Standardized provider errors
var (
	ErrConfigReadFailed  = errors.New("resolver: [provider] config read failed")
	ErrConfigParseFailed = errors.New("resolver: [provider] config parse failed")
	ErrEmptyZoneConfig   = errors.New("resolver: [provider] empty zone config")
	ErrTLDNotSupported   = errors.New("resolver: [provider] TLD not supported")
	ErrInvalidDomain     = errors.New("resolver: [provider] invalid domain")
)

// ZoneConfig holds WHOIS server configuration for a specific TLD or IP
type ZoneConfig struct {
	Host  string // WHOIS server host
	Query string // custom query format (optional)
}

// Provider represents a resolver config provider for TLD and IP WHOIS servers
type Provider struct {
	tldServers map[string]ZoneConfig // map of TLD → WHOIS server config
	ipServer   *ZoneConfig           // fallback WHOIS server for IP addresses
}

// NewProvider initializes a new provider with given zone and IP config
func NewProvider(tldServers map[string]ZoneConfig, ipServer *ZoneConfig) *Provider {
	return &Provider{
		tldServers: tldServers,
		ipServer:   ipServer,
	}
}

// ParseConfigFromFile loads and parses a JSON WHOIS config from the specified path
func ParseConfigFromFile(path string) (*Provider, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigReadFailed, err)
	}

	// raw is used to unmarshal the JSON config into temporary structures
	var raw struct {
		Zones map[string]struct {
			Host  string `json:"host"`
			Query string `json:"query"`
		} `json:"zones"`
		IP struct {
			Host  string `json:"host"`
			Query string `json:"query"`
		} `json:"ip"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigParseFailed, err)
	}

	// Transform raw config into internal representation
	tldServers := make(map[string]ZoneConfig)
	for tld, cfg := range raw.Zones {
		if len(cfg.Host) > 0 {
			tldServers[tld] = ZoneConfig{
				Host:  cfg.Host,
				Query: cfg.Query,
			}
		}
	}

	if len(tldServers) == 0 {
		return nil, ErrEmptyZoneConfig
	}

	var ipServer *ZoneConfig
	if len(raw.IP.Host) > 0 {
		ipServer = &ZoneConfig{
			Host:  raw.IP.Host,
			Query: raw.IP.Query,
		}
	}

	return NewProvider(tldServers, ipServer), nil
}

// GetServer resolves a WHOIS server config based on the domain's TLD
func (p *Provider) GetServer(domain *domain.Domain) (*ZoneConfig, error) {
	if domain.TLD == "" {
		return nil, ErrInvalidDomain
	}

	server, ok := p.tldServers[domain.TLD]
	if !ok {
		return nil, ErrTLDNotSupported
	}

	return &server, nil
}
