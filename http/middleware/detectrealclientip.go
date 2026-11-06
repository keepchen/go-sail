package middleware

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

// DetectRealClientIP retrieves the client's real IP address
//
// Read and process the X-Forwarded-For request header
//
// # Tip
//
// For safer and more accurate client IP detection, use gin.Engine.SetTrustedProxies() to configure trusted proxy IP addresses
func DetectRealClientIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Read the X-Forwarded-For request header
		xForwardedFor := c.GetHeader("X-Forwarded-For")

		if len(xForwardedFor) == 0 {
			c.Set("realClientIp", c.ClientIP()) // Fall back to the default method
			c.Next()
			return
		}

		ips := strings.Split(xForwardedFor, ",")
		//Use the last IP address
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
