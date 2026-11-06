package utils

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
)

type domainImpl struct {
}

type IDomain interface {
	// Validate checks whether a domain name is valid
	//
	// Wildcards are not supported
	Validate(domain string) bool
	// ValidateWithWildcard checks whether a domain name is valid
	//
	// Includes wildcard rule validation
	//
	// eg.
	//
	// example.com    -> ✔
	//
	// x.example.com  -> ✔
	//
	// *.example.com  -> ✔
	//
	// *example.com   -> ✖
	//
	// example*.com   -> ✖
	ValidateWithWildcard(domain string) bool
	// RefererMatch checks whether the referer matches a domain
	//
	// Tips:
	//
	// 1.Referer and domain comparisons are case-insensitive
	//
	// 2.This function supports wildcard domain matching
	RefererMatch(referer, domain string) bool
	// LookupCNAME queries a domain's CNAME record
	//
	// # Note
	//
	// nslookup appends a "." to domain names,
	// so this function checks whether cnameTarget ends with "."
	// before comparing; if not, it appends "." and then compares.
	LookupCNAME(domain, cnameTarget string) bool
}

var doi IDomain = &domainImpl{}

// Domain instantiates domain utilities
func Domain() IDomain {
	return doi
}

// Validate checks whether a domain name is valid
//
// Wildcards are not supported
func (domainImpl) Validate(domain string) bool {
	ok, err := regexp.MatchString("^([a-zA-Z0-9-]+\\.)+[a-zA-Z]{2,}$", domain)

	return ok && err == nil
}

// ValidateWithWildcard checks whether a domain name is valid
//
// Includes wildcard rule validation
//
// eg.
//
// example.com    -> ✔
//
// x.example.com  -> ✔
//
// *.example.com  -> ✔
//
// *example.com   -> ✖
//
// example*.com   -> ✖
func (domainImpl) ValidateWithWildcard(domain string) bool {
	ok, err := regexp.MatchString("^(\\*\\.)?([a-zA-Z0-9-]+\\.)+[a-zA-Z]{2,}$", domain)

	return ok && err == nil
}

// RefererMatch checks whether the referer matches a domain
//
// Tips:
//
// 1.Referer and domain comparisons are case-insensitive
//
// 2.This function supports wildcard domain matching
func (domainImpl) RefererMatch(referer, domain string) bool {
	if domain == "*" {
		return true
	}
	u, err := url.Parse(referer)
	if err != nil {
		return false
	}
	var host string
	if strings.Contains(u.Host, ":") {
		h, _, _ := net.SplitHostPort(u.Host)
		host = h
	} else {
		host = u.Host
	}
	//If the domain has no asterisk, compare the two directly
	if !strings.HasPrefix(domain, "*") {
		return strings.EqualFold(host, domain)
	}
	//Remove the asterisk and check for containment
	clearDomain := strings.Replace(domain, "*", "", 1)

	// host=a.b.com
	// clearDomain=.b.com
	return strings.HasSuffix(strings.ToLower(host), strings.ToLower(clearDomain))
}

// LookupCNAME queries a domain's CNAME record
//
// # Note
//
// nslookup appends a "." to domain names,
// so this function checks whether cnameTarget ends with "."
// before comparing; if not, it appends "." and then compares.
func (domainImpl) LookupCNAME(domain, cnameTarget string) bool {
	if !strings.HasSuffix(cnameTarget, ".") {
		cnameTarget = fmt.Sprintf("%s.", cnameTarget)
	}
	cname, err := net.LookupCNAME(domain)
	if err != nil {
		fmt.Printf("LookupCNAME error: %v\n", err)
		return false
	}

	return cname == cnameTarget
}
