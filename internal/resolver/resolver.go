package resolver

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"paulmatvienko/whois/internal/domain"
	"regexp"
	"strings"
	"time"
)

// Standardized resolver errors
var (
	ErrConnectionFailed   = errors.New("resolver: failed to connect")
	ErrSetDeadlineFailed  = errors.New("resolver: failed to set deadline")
	ErrWriteQueryFailed   = errors.New("resolver: failed to write query")
	ErrReadResponseFailed = errors.New("resolver: failed to read response")
	ErrRequestTimeout     = errors.New("resolver: request timeout")
	ErrWhoisQueryFailed   = errors.New("resolver: whois query failed")
	ErrInfinityLoop       = errors.New("resolver: referral loop detected or max follow exceeded")
)

// Resolver handles WHOIS query sending and response reading
type Resolver struct {
	timeout     time.Duration // Timeout for connection, read and write operations
	port        int           // WHOIS server port, typically 43
	maxResponse int64         // Maximum allowed response size in bytes (to prevent DoS)
	maxFollow   int           // Maximum resolve iterations in resolve loop
}

// NewResolver creates a new Resolver instance with given parameters
func NewResolver(timeout time.Duration, maxResponse int64, maxFollow int, port int) *Resolver {
	return &Resolver{
		timeout:     timeout,
		port:        port,
		maxFollow:   maxFollow,
		maxResponse: maxResponse,
	}
}

// ResolveWithReferrals performs a WHOIS query to the specified server for the given domain,
// automatically following referral servers up to a maximum depth (r.maxFollow).
func (r *Resolver) ResolveWithReferrals(ctx context.Context, domainName *domain.Domain, server *ZoneConfig) (string, string, error) {
	visited := make(map[string]struct{})

	for i := 0; i <= r.maxFollow; i++ {
		if _, seen := visited[server.Host]; seen {
			break
		}
		visited[server.Host] = struct{}{}

		raw, err := r.Resolve(ctx, domainName, server)
		if err != nil {
			return "", server.Host, fmt.Errorf("%w: %v", ErrWhoisQueryFailed, err)
		}

		ref := r.FindReferralServer(raw)
		if ref == "" || strings.EqualFold(ref, server.Host) {
			return raw, server.Host, nil
		}

		server = &ZoneConfig{Host: ref}
	}

	return "", server.Host, ErrInfinityLoop
}

// Resolve performs a WHOIS query to the specified server for the given domain.
// The context allows cancellation and timeout control.
func (r *Resolver) Resolve(ctx context.Context, domain *domain.Domain, server *ZoneConfig) (string, error) {
	dialer := &net.Dialer{Timeout: r.timeout}
	addr := net.JoinHostPort(server.Host, fmt.Sprintf("%d", r.port))

	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrConnectionFailed, err)
	}
	defer conn.Close()

	// Set an overall deadline for read/write operations
	deadline := time.Now().Add(r.timeout)
	err = conn.SetDeadline(deadline)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrSetDeadlineFailed, err)
	}

	query := buildQuery(server, domain.String())
	if _, err := conn.Write([]byte(query)); err != nil {
		return "", fmt.Errorf("%w: %v", ErrWriteQueryFailed, err)
	}

	data, err := io.ReadAll(io.LimitReader(conn, r.maxResponse))
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return "", fmt.Errorf("%w: %v", ErrRequestTimeout, err)
		}
		return "", fmt.Errorf("%w: %v", ErrReadResponseFailed, err)
	}

	return string(data), nil
}

func (r *Resolver) FindReferralServer(data string) string {
	re := regexp.MustCompile(`(?mi)^(Whois Server|ReferralServer):\s*(?:whois://)?([^\s]+)\s*$`)
	matches := re.FindStringSubmatch(data)
	if len(matches) >= 3 {
		return strings.TrimSpace(matches[2])
	}
	return ""
}
