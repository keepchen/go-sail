package utils

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

type domainValidStdDeprecated struct {
	domain string
	valid  bool
}

func TestValidateDomain(t *testing.T) {
	var domains = []domainValidStdDeprecated{
		{"example.com", true},
		{"*.example.com", false},
		{"*example.com", false},
		{"x.example.com", true},
		{"a.b.example.com", true},
		{"example*.com", false},
		{"example-x.com", true},
		{"example-*.com", false},
		{"*-example.com", false},
	}

	for _, v := range domains {
		assert.Equal(t, ValidateDomain(v.domain), v.valid)
	}
}

func TestValidateDomainWithWildcard(t *testing.T) {
	var domains = []domainValidStdDeprecated{
		{"example.com", true},
		{"*.example.com", true},
		{"*example.com", false},
		{"x.example.com", true},
		{"a.b.example.com", true},
		{"example*.com", false},
		{"example-x.com", true},
		{"example-*.com", false},
		{"*-example.com", false},
	}

	for _, v := range domains {
		assert.Equal(t, ValidateDomainWithWildcard(v.domain), v.valid)
	}
}

type refererDomainValidStdDeprecated struct {
	referer string
	domain  string
	valid   bool
}

func TestRefererMatchDomain(t *testing.T) {
	var domains = []refererDomainValidStdDeprecated{
		{"https://example.com", "example.com", true},
		{"https://example.com:8443", "example.com", true},
		{"https://example.com", "*.example.com", false},
		{"https://a.example.com", "*.example.com", true},
		{"https://a.example.com:8443", "*.example.com", true},
		{"https://a.aexample.com", "*.example.com", false},
		{"https://a.b.example.com", "*.example.com", true},
		{"https://a.b.example.com:8443", "*.example.com", true},
	}
	for _, v := range domains {
		t.Log(v.referer, v.domain)
		assert.Equal(t, RefererMatchDomain(v.referer, v.domain), v.valid)
	}
}

func TestLookupCNAME(t *testing.T) {
	t.Log(LookupCNAME("cos.stardots.ink", "stardots-1251905630.cos.ap-singapore.myqcloud.com."))
}
