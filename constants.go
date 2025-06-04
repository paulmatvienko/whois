package whois

import (
	"time"
)

const (
	DefaultTimeout              = 30 * time.Second
	DefaultFollow               = 10
	DefaultMaxResponse          = 512 * 1024
	QueryPort                   = 43
	DefaultWhoisServersFilePath = "./servers.json"
)
