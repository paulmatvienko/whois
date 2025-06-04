package domain

import (
	"golang.org/x/net/idna"
	"strings"
)

func IsValid(domain string) bool {
	domain = strings.TrimSpace(domain)
	if domain == "" || len(domain) > 253 {
		return false
	}

	ascii, err := idna.ToASCII(domain)
	if err != nil || len(ascii) > 253 {
		return false
	}

	labels := strings.Split(ascii, ".")
	if len(labels) < 2 {
		return false // must contain at least one dot
	}

	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 {
			return false
		}

		if label[0] == '-' || label[len(label)-1] == '-' {
			return false
		}

		for i := 0; i < len(label); i++ {
			ch := label[i]
			if !((ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '-') {
				return false
			}
		}
	}

	return true
}
