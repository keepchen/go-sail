package middleware

import (
	"net/http"
	"testing"
)

func TestWithCors(t *testing.T) {
	t.Run("WithCors-WithOrigin", func(t *testing.T) {
		c, _ := createTestContextAndEngine()

		req, _ := http.NewRequest("GET", "/test?name=foo", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", "127.0.0.1")

		WithCors(nil)(c)
	})

	t.Run("WithCors-WithoutOrigin", func(t *testing.T) {
		c, _ := createTestContextAndEngine()

		req, _ := http.NewRequest("GET", "/test?name=foo", nil)
		req.Header.Set("Content-Type", "application/json")

		WithCors(nil)(c)
	})

	t.Run("WithCors-WithHeaders", func(t *testing.T) {
		c, _ := createTestContextAndEngine()

		req, _ := http.NewRequest("GET", "/test?name=foo", nil)
		req.Header.Set("Content-Type", "application/json")

		headers := map[string]string{
			"Access-Control-Allow-Credentials": "true",
		}

		WithCors(headers)(c)
	})
}

func TestWithCorsOnlyOptions(t *testing.T) {
	t.Run("WithCorsOnlyOptions-WithOrigin", func(t *testing.T) {
		c, _ := createTestContextAndEngine()

		req, _ := http.NewRequest("GET", "/test?name=foo", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", "127.0.0.1")

		WithCorsOnlyOptions(nil)(c)
	})

	t.Run("WithCorsOnlyOptions-WithoutOrigin", func(t *testing.T) {
		c, _ := createTestContextAndEngine()

		req, _ := http.NewRequest("GET", "/test?name=foo", nil)
		req.Header.Set("Content-Type", "application/json")

		WithCorsOnlyOptions(nil)(c)
	})

	t.Run("WithCorsOnlyOptions-WithHeaders", func(t *testing.T) {
		c, _ := createTestContextAndEngine()

		req, _ := http.NewRequest("GET", "/test?name=foo", nil)
		req.Header.Set("Content-Type", "application/json")

		headers := map[string]string{
			"Access-Control-Allow-Credentials": "true",
		}

		WithCorsOnlyOptions(headers)(c)
	})

	t.Run("WithCorsOnlyOptions-WithOptions", func(t *testing.T) {
		c, _ := createTestContextAndEngine()

		req, _ := http.NewRequest("OPTIONS", "/test?name=foo", nil)
		req.Header.Set("Content-Type", "application/json")

		headers := map[string]string{
			"Access-Control-Allow-Credentials": "true",
		}

		WithCorsOnlyOptions(headers)(c)
	})
}
