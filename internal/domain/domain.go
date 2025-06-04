package domain

import (
	"errors"
	"golang.org/x/net/idna"
	"golang.org/x/net/publicsuffix"
	"strings"
)

// Standardized domain errors
var (
	ErrPublicSuffix  = errors.New("domain: failed to determine public suffix")
	ErrEmptyDomain   = errors.New("domain: empty domain")
	ErrInvalidDomain = errors.New("domain: invalid domain")
)

// Domain represents a parsed domain name
type Domain struct {
	Raw      string // original input string
	Name     string // base domain name (without subdomains and TLD)
	TLD      string // public suffix (TLD)
	ICANN    bool   // whether the TLD is managed by ICANN
	IsCustom bool   // whether the TLD is private/custom
}

// Parse parses a domain string and returns a Domain struct
func Parse(raw string) (*Domain, error) {
	// Normalize input
	norm := strings.ToLower(strings.TrimSpace(raw))
	if norm == "" {
		return nil, ErrEmptyDomain
	}

	// Convert to ASCII (Punycode for IDN)
	ascii, err := idna.Lookup.ToASCII(norm)
	if err != nil {
		return nil, ErrInvalidDomain
	}

	// Basic syntax validation
	if !IsValid(ascii) {
		return nil, ErrInvalidDomain
	}

	// Determine public suffix
	tld, icann := publicsuffix.PublicSuffix(ascii)
	if tld == "" {
		return nil, ErrPublicSuffix
	}

	// Determine base domain (eTLD+1)
	baseDomain, err := publicsuffix.EffectiveTLDPlusOne(ascii)
	if err != nil {
		return nil, ErrPublicSuffix
	}

	// Extract domain name (part before eTLD)
	name := strings.TrimSuffix(baseDomain, "."+tld)

	return &Domain{
		Raw:      raw,
		Name:     name,
		TLD:      tld,
		ICANN:    icann,
		IsCustom: !icann,
	}, nil
}

// String returns the normalized string representation of the domain
func (d *Domain) String() string {
	if d.Name == "" {
		return d.TLD
	}
	return d.Name + "." + d.TLD
}
