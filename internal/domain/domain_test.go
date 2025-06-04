package domain

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantName string
		wantTLD  string
		wantErr  error
	}{
		{
			name:     "simple domain",
			input:    "example.com",
			wantName: "example",
			wantTLD:  "com",
		},
		{
			name:     "subdomain",
			input:    "sub.example.co.uk",
			wantName: "example",
			wantTLD:  "co.uk",
		},
		{
			name:     "IDN domain",
			input:    "пример.рф",
			wantName: "xn--e1afmkfd",
			wantTLD:  "xn--p1ai",
		},
		{
			name:     "second level TLD",
			input:    "example.academy",
			wantName: "example",
			wantTLD:  "academy",
		},
		{
			name:     "no subdomains",
			input:    "example.com",
			wantName: "example",
			wantTLD:  "com",
		},
		{
			name:  "empty string",
			input: "",
		},
		{
			name:  "whitespace only",
			input: "  ",
		},
		{
			name:  "invalid domain",
			input: "example..com",
		},
		{
			name:  "dot only",
			input: ".",
		},
		{
			name:     "just TLD",
			input:    "com",
			wantName: "",
			wantTLD:  "com",
		},
		{
			name:     "multi-level subdomains",
			input:    "a.b.c.d.example.io",
			wantName: "example",
			wantTLD:  "io",
		},
		{
			name:     "private TLD",
			input:    "example.blogspot.com",
			wantName: "example",
			wantTLD:  "blogspot.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if err != nil {
				return
			}

			if got.Name != tt.wantName {
				t.Errorf("Parse() name = %v, want %v", got.Name, tt.wantName)
			}

			if got.TLD != tt.wantTLD {
				t.Errorf("Parse() TLD = %v, want %v", got.TLD, tt.wantTLD)
			}

			// Проверяем, что String() возвращает корректное значение
			expectedString := tt.wantName
			if expectedString != "" {
				expectedString += "." + tt.wantTLD
			} else {
				expectedString = tt.wantTLD
			}

			if got.String() != expectedString {
				t.Errorf("String() = %v, want %v", got.String(), expectedString)
			}

			// Проверяем, что Raw() возвращает оригинальную строку
			if got.Raw != tt.input {
				t.Errorf("Raw() = %v, want %v", got.Raw, tt.input)
			}
		})
	}
}

func TestTLDProperties(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantICANN bool
	}{
		{
			name:      "ICANN TLD",
			input:     "example.com",
			wantICANN: true,
		},
		{
			name:      "private TLD",
			input:     "example.blogspot.com",
			wantICANN: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() failed: %v", err)
			}

			if got.ICANN != tt.wantICANN {
				t.Errorf("IsICANN() = %v, want %v", got.ICANN, tt.wantICANN)
			}

			if got.IsCustom == tt.wantICANN {
				t.Errorf("IsCustom() = %v, want %v", got.IsCustom, !tt.wantICANN)
			}
		})
	}
}

func TestConcurrentParsing(t *testing.T) {
	const numWorkers = 100
	results := make(chan *Domain, numWorkers)
	errs := make(chan error, numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func() {
			d, err := Parse("example.com")
			if err != nil {
				errs <- err
				return
			}
			results <- d
		}()
	}

	for i := 0; i < numWorkers; i++ {
		select {
		case d := <-results:
			if d.Name != "example" || d.TLD != "com" {
				t.Errorf("Unexpected result: %+v", d)
			}
		case err := <-errs:
			t.Errorf("Unexpected error: %v", err)
		}
	}
}
