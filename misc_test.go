package onens

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormaliseDomain(t *testing.T) {
	tests := []struct {
		input  string
		output string
		err    error
	}{
		{"", "", nil},
		{".", ".", nil},
		{"country", "country", nil},
		{"COUNTRY", "country", nil},
		{".country", ".country", nil},
		{".country.", ".country.", nil},
		{"1ns.country", "1ns.country", nil},
		{".1ns.country", ".1ns.country", nil},
		{"subdomain.1ns.country", "subdomain.1ns.country", nil},
		{"*.1ns.country", "*.1ns.country", nil},
		{"omg.thetoken.country", "omg.thetoken.country", nil},
		{"_underscore.thetoken.country", "_underscore.thetoken.country", nil},
		{"點看.country", "點看.country", nil},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := NormaliseDomain(tt.input)
			if tt.err != nil {
				if err == nil {
					t.Fatalf("missing expected error")
				}
				if tt.err.Error() != err.Error() {
					t.Errorf("unexpected error value %v", err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error %v", err)
				}
				if tt.output != result {
					t.Errorf("%v => %v (expected %v)", tt.input, result, tt.output)
				}
			}
		})
	}
}

func TestNormaliseDomainStrict(t *testing.T) {
	tests := []struct {
		input  string
		output string
		err    error
	}{
		{"", "", nil},
		{".", ".", nil},
		{"country", "country", nil},
		{"COUNTRY", "country", nil},
		{".country", ".country", nil},
		{".country.", ".country.", nil},
		{"1ns.country", "1ns.country", nil},
		{".1ns.country", ".1ns.country", nil},
		{"subdomain.1ns.country", "subdomain.1ns.country", nil},
		{"*.1ns.country", "*.1ns.country", nil},
		{"omg.thetoken.country", "omg.thetoken.country", nil},
		{"_underscore.thetoken.country", "", errors.New("idna: disallowed rune U+005F")},
		{"點看.country", "點看.country", nil},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := NormaliseDomainStrict(tt.input)
			if tt.err != nil {
				if err == nil {
					t.Fatalf("missing expected error")
				}
				if tt.err.Error() != err.Error() {
					t.Errorf("unexpected error value %v", err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error %v", err)
				}
				if tt.output != result {
					t.Errorf("%v => %v (expected %v)", tt.input, result, tt.output)
				}
			}
		})
	}
}

func TestTld(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"", ""},
		{".", ""},
		{"country", "country"},
		{"COUNTRY", "country"},
		{".country", "country"},
		{"1ns.country", "country"},
		{".1ns.country", "country"},
		{"subdomain.1ns.country", "country"},
	}

	for _, tt := range tests {
		result := Tld(tt.input)
		if tt.output != result {
			t.Errorf("Failure: %v => %v (expected %v)\n", tt.input, result, tt.output)
		}
	}
}

func TestDomainPart(t *testing.T) {
	tests := []struct {
		input  string
		part   int
		output string
		err    bool
	}{
		{"", 1, "", false},
		{"", 2, "", true},
		{"", -1, "", false},
		{"", -2, "", true},
		{".", 1, "", false},
		{".", 2, "", false},
		{".", 3, "", true},
		{".", -1, "", false},
		{".", -2, "", false},
		{".", -3, "", true},
		{"COUNTRY", 1, "country", false},
		{"COUNTRY", 2, "", true},
		{"COUNTRY", -1, "country", false},
		{"COUNTRY", -2, "", true},
		{".COUNTRY", 1, "", false},
		{".COUNTRY", 2, "country", false},
		{".COUNTRY", 3, "", true},
		{".COUNTRY", -1, "country", false},
		{".COUNTRY", -2, "", false},
		{".COUNTRY", -3, "", true},
		{"1ns.country", 1, "1ns", false},
		{"1ns.country", 2, "country", false},
		{"1ns.country", 3, "", true},
		{"1ns.country", -1, "country", false},
		{"1ns.country", -2, "1ns", false},
		{"1ns.country", -3, "", true},
		{".1ns.country", 1, "", false},
		{".1ns.country", 2, "1ns", false},
		{".1ns.country", 3, "country", false},
		{".1ns.country", 4, "", true},
		{".1ns.country", -1, "country", false},
		{".1ns.country", -2, "1ns", false},
		{".1ns.country", -3, "", false},
		{".1ns.country", -4, "", true},
		{"subdomain.1ns.country", 1, "subdomain", false},
		{"subdomain.1ns.country", 2, "1ns", false},
		{"subdomain.1ns.country", 3, "country", false},
		{"subdomain.1ns.country", 4, "", true},
		{"subdomain.1ns.country", -1, "country", false},
		{"subdomain.1ns.country", -2, "1ns", false},
		{"subdomain.1ns.country", -3, "subdomain", false},
		{"subdomain.1ns.country", -4, "", true},
		{"a.b.c", 1, "a", false},
		{"a.b.c", 2, "b", false},
		{"a.b.c", 3, "c", false},
		{"a.b.c", 4, "", true},
		{"a.b.c", -1, "c", false},
		{"a.b.c", -2, "b", false},
		{"a.b.c", -3, "a", false},
		{"a.b.c", -4, "", true},
	}

	for _, tt := range tests {
		result, err := DomainPart(tt.input, tt.part)
		if err != nil && !tt.err {
			t.Errorf("Failure: %v, %v => error (unexpected)\n", tt.input, tt.part)
		}
		if err == nil && tt.err {
			t.Errorf("Failure: %v, %v => no error (unexpected)\n", tt.input, tt.part)
		}
		if tt.output != result {
			t.Errorf("Failure: %v, %v => %v (expected %v)\n", tt.input, tt.part, result, tt.output)
		}
	}
}

func TestUnqualifiedName(t *testing.T) {
	tests := []struct {
		domain string
		root   string
		name   string
		err    error
	}{
		{
			domain: "",
			root:   "",
			name:   "",
		},
		{
			domain: "1ns.country",
			root:   "country",
			name:   "1ns",
		},
	}

	for i, test := range tests {
		name, err := UnqualifiedName(test.domain, test.root)
		if test.err != nil {
			assert.Equal(t, test.err, err, fmt.Sprintf("Incorrect error at test %d", i))
		} else {
			require.Nil(t, err, fmt.Sprintf("Unexpected error at test %d", i))
			assert.Equal(t, test.name, name, fmt.Sprintf("Incorrect result at test %d", i))
		}
	}
}
