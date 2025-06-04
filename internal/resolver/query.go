package resolver

import (
	"strings"
)

// buildQuery constructs the WHOIS query string to send to the server.
func buildQuery(server *ZoneConfig, domain string) string {
	query := server.Query
	if query == "" {
		query = "$addr\r\n"
	}
	return strings.ReplaceAll(query, "$addr", domain)
}
