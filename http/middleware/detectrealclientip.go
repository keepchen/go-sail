package middleware

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

// DetectRealClientIP 获取客户端真实ip
//
// 从请求头中读取 X-Forwarded-For 字段进行处理
//
// # 提示
//
// 为了更加安全准确的获取客户端真实ip，建议使用 gin.Engine.SetTrustedProxies() 方法设置授信的代理ip
func DetectRealClientIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		//读取请求头中的X-Forwarded-For字段
		xForwardedFor := c.GetHeader("X-Forwarded-For")

		if len(xForwardedFor) == 0 {
			c.Set("realClientIp", c.ClientIP()) // 回退到默认方法
			c.Next()
			return
		}

		ips := strings.Split(xForwardedFor, ",")
		//取最后一个IP
		realIP := strings.TrimSpace(ips[len(ips)-1])
		ip := net.ParseIP(realIP)
		if ip == nil {
			c.Set("realClientIp", "")
			c.Next()
			return
		}
		c.Set("realClientIp", ip)
		c.Next()
	}
}
