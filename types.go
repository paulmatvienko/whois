package whois

import (
	"paulmatvienko/whois/internal/resolver"
	"time"
)

type Options struct {
	whoisServersFilePath string
	server               *resolver.ZoneConfig
	timeout              time.Duration
	follow               int
	maxResponse          int64
}

type Result struct {
	RawData        string
	Domain         string
	ReferralServer string
}
