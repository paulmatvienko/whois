package whois

import (
	"context"
	"fmt"
	"paulmatvienko/whois/internal/domain"
	"paulmatvienko/whois/internal/resolver"
)

type Whois struct {
	resolver *resolver.Resolver
	provider *resolver.Provider
	opts     *Options
}

// New initializes a new Whois client with provided options.
// Defaults are applied for missing timeout, follow depth and max response size.
func New(opts Options) (*Whois, error) {
	if opts.timeout == 0 {
		opts.timeout = DefaultTimeout
	}
	if opts.maxResponse == 0 {
		opts.maxResponse = DefaultMaxResponse
	}
	if opts.follow == 0 {
		opts.follow = DefaultFollow
	}
	if opts.whoisServersFilePath == "" {
		opts.whoisServersFilePath = DefaultWhoisServersFilePath
	}

	r := resolver.NewResolver(opts.timeout, opts.maxResponse, opts.follow, QueryPort)
	p, err := resolver.ParseConfigFromFile(opts.whoisServersFilePath)
	if err != nil {
		return nil, err
	}

	return &Whois{
		resolver: r,
		provider: p,
		opts:     &opts,
	}, nil
}

// Lookup performs a WHOIS lookup for the given domain.
// It automatically follows referral servers (up to opts.follow).
func (w *Whois) Lookup(ctx context.Context, input string) (*Result, *Error) {
	if !domain.IsValid(input) {
		return nil, NewError(input, "", ErrInvalidDomain, fmt.Errorf("invalid domain: %q", input))
	}

	domainName, err := domain.Parse(input)
	if err != nil {
		return nil, NewError(input, "", ErrInvalidDomain, err)
	}

	// Select initial WHOIS server
	server := w.opts.server
	if server == nil {
		server, err = w.provider.GetServer(domainName)
		if err != nil {
			return nil, NewError(domainName.String(), "", ErrServerNotFound, err)
		}
	}

	// Follow a referral chain
	raw, finalHost, err := w.resolver.ResolveWithReferrals(ctx, domainName, server)
	if err != nil {
		return nil, NewError(domainName.String(), finalHost, ErrQuery, err)
	}

	return &Result{
		RawData:        raw,
		Domain:         domainName.String(),
		ReferralServer: finalHost,
	}, nil
}
