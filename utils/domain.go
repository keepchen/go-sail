package utils

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
)

// ValidateDomain 验证域名是否合法
//
// 不支持通配符
func ValidateDomain(domain string) bool {
	ok, err := regexp.MatchString("^([a-zA-Z0-9-]+\\.)+[a-zA-Z]{2,}$", domain)
	if err != nil {
		return false
	}

	return ok
}

// ValidateDomainWithWildcard 验证域名是否合法
//
// 包含通配符规则验证
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
func ValidateDomainWithWildcard(domain string) bool {
	ok, err := regexp.MatchString("^(\\*\\.)?([a-zA-Z0-9-]+\\.)+[a-zA-Z]{2,}$", domain)
	if err != nil {
		return false
	}

	return ok
}

// RefererMatchDomain refer是否匹配域名
//
// 提示：
//
// 1.referer和域名大小写不敏感
//
// 2.此函数支持通配符域名检测
func RefererMatchDomain(referer, domain string) bool {
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
	//域名不包含星号，直接检测两者是否相等
	if !strings.HasPrefix(domain, "*") {
		return strings.EqualFold(host, domain)
	}
	//去除掉星号，检测是否包含
	clearDomain := strings.Replace(domain, "*", "", 1)

	// host=a.b.com
	// clearDomain=.b.com
	return strings.HasSuffix(strings.ToLower(host), strings.ToLower(clearDomain))
}

// LookupCNAME 查询域名cname记录
//
// # 注意
//
// nslookup的结果会在域名后增加一个「.」，
// 因此此函数在对比前会检测cnameTarget参数是否以「.」
// 结尾，如果不是，则会在结尾加上「.」，然后再进行对比。
func LookupCNAME(domain, cnameTarget string) bool {
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
