package utils

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

type domainValidStd struct {
	domain string
	valid  bool
}

func TestDomainImplValidate(t *testing.T) {
	var domains = []domainValidStd{
		{"example.com", true},
		{"*.example.com", false},
		{"*example.com", false},
		{"x.example.com", true},
		{"a.b.example.com", true},
		{"example*.com", false},
		{"example-x.com", true},
		{"example-*.com", false},
		{"*-example.com", false},
		{"-", false},
	}

	for _, v := range domains {
		assert.Equal(t, Domain().Validate(v.domain), v.valid)
	}
}

func TestDomainImplValidateWithWildcard(t *testing.T) {
	var domains = []domainValidStd{
		{"example.com", true},
		{"*.example.com", true},
		{"*example.com", false},
		{"x.example.com", true},
		{"a.b.example.com", true},
		{"example*.com", false},
		{"example-x.com", true},
		{"example-*.com", false},
		{"*-example.com", false},
		{"-", false},
	}

	for _, v := range domains {
		assert.Equal(t, Domain().ValidateWithWildcard(v.domain), v.valid)
	}
}

type refererDomainValidStd struct {
	referer string
	domain  string
	valid   bool
}

func TestDomainImplRefererMatch(t *testing.T) {
	var domains = []refererDomainValidStd{
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
		assert.Equal(t, Domain().RefererMatch(v.referer, v.domain), v.valid)
	}
}

func TestDomainImplLookupCNAME(t *testing.T) {
	t.Run("DomainImplLookupCNAME", func(t *testing.T) {
		t.Log(Domain().LookupCNAME("cos.stardots.ink", "stardots-1251905630.cos.ap-singapore.myqcloud.com."))
	})

	t.Run("DomainImplLookupCNAME-Error", func(t *testing.T) {
		t.Log(Domain().LookupCNAME("-", "-"))
	})
}
