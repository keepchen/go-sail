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
	// Validate 验证域名是否合法
	//
	// 不支持通配符
	Validate(domain string) bool
	// ValidateWithWildcard 验证域名是否合法
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
	ValidateWithWildcard(domain string) bool
	// RefererMatch refer是否匹配域名
	//
	// 提示：
	//
	// 1.referer和域名大小写不敏感
	//
	// 2.此函数支持通配符域名检测
	RefererMatch(referer, domain string) bool
	// LookupCNAME 查询域名cname记录
	//
	// # 注意
	//
	// nslookup的结果会在域名后增加一个「.」，
	// 因此此函数在对比前会检测cnameTarget参数是否以「.」
	// 结尾，如果不是，则会在结尾加上「.」，然后再进行对比。
	LookupCNAME(domain, cnameTarget string) bool
}

var _ IDomain = &domainImpl{}

// Domain 实例化domain工具类
func Domain() IDomain {
	return &domainImpl{}
}

// Validate 验证域名是否合法
//
// 不支持通配符
func (domainImpl) Validate(domain string) bool {
	ok, err := regexp.MatchString("^([a-zA-Z0-9-]+\\.)+[a-zA-Z]{2,}$", domain)
	if err != nil {
		return false
	}

	return ok
}

// ValidateWithWildcard 验证域名是否合法
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
func (domainImpl) ValidateWithWildcard(domain string) bool {
	ok, err := regexp.MatchString("^(\\*\\.)?([a-zA-Z0-9-]+\\.)+[a-zA-Z]{2,}$", domain)
	if err != nil {
		return false
	}

	return ok
}

// RefererMatch refer是否匹配域名
//
// 提示：
//
// 1.referer和域名大小写不敏感
//
// 2.此函数支持通配符域名检测
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
