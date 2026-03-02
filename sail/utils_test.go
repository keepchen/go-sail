package sail

import (
	"testing"

	"github.com/keepchen/go-sail/v3/lib/redis"

	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {
	t.Run("Utils", func(t *testing.T) {
		t.Log(Utils())
	})
}

func TestUtilsFunc(t *testing.T) {
	t.Run("UtilsFunc-Aes", func(t *testing.T) {
		t.Log(Utils().Aes())
	})
	t.Run("UtilsFunc-Base64", func(t *testing.T) {
		t.Log(Utils().Base64())
	})
	t.Run("UtilsFunc-Cert", func(t *testing.T) {
		t.Log(Utils().Cert())
	})
	t.Run("UtilsFunc-Crc32", func(t *testing.T) {
		t.Log(Utils().Crc32())
	})
	t.Run("UtilsFunc-Crc64", func(t *testing.T) {
		t.Log(Utils().Crc64())
	})
	t.Run("UtilsFunc-Datetime", func(t *testing.T) {
		t.Log(Utils().Datetime())
	})
	t.Run("UtilsFunc-Domain", func(t *testing.T) {
		t.Log(Utils().Domain())
	})
	t.Run("UtilsFunc-File", func(t *testing.T) {
		t.Log(Utils().File())
	})
	t.Run("UtilsFunc-Gzip", func(t *testing.T) {
		t.Log(Utils().Gzip())
	})
	t.Run("UtilsFunc-HttpClient", func(t *testing.T) {
		t.Log(Utils().HttpClient())
	})
	t.Run("UtilsFunc-IP", func(t *testing.T) {
		t.Log(Utils().IP())
	})
	t.Run("UtilsFunc-MD5", func(t *testing.T) {
		t.Log(Utils().MD5())
	})
	t.Run("UtilsFunc-Number", func(t *testing.T) {
		t.Log(Utils().Number())
	})
	t.Run("UtilsFunc-RedisLocker", func(t *testing.T) {
		if redis.GetInstance() == nil && redis.GetClusterInstance() == nil {
			assert.Panics(t, func() {
				t.Log(Utils().RedisLocker())
			})
		}
		if redis.GetInstance() != nil {
			t.Log(Utils().RedisLocker(redis.GetInstance()))
		}
		if redis.GetClusterInstance() != nil {
			t.Log(Utils().RedisLocker(redis.GetClusterInstance()))
		}
	})
	t.Run("UtilsFunc-RSA", func(t *testing.T) {
		t.Log(Utils().RSA())
	})
	t.Run("UtilsFunc-Signal", func(t *testing.T) {
		t.Log(Utils().Signal())
	})
	t.Run("UtilsFunc-SM4", func(t *testing.T) {
		t.Log(Utils().SM4())
	})
	t.Run("UtilsFunc-String", func(t *testing.T) {
		t.Log(Utils().String())
	})
	t.Run("UtilsFunc-Swagger", func(t *testing.T) {
		t.Log(Utils().Swagger())
	})
	t.Run("UtilsFunc-Validator", func(t *testing.T) {
		t.Log(Utils().Validator())
	})
	t.Run("UtilsFunc-Version", func(t *testing.T) {
		t.Log(Utils().Version())
	})
	t.Run("UtilsFunc-WebPush", func(t *testing.T) {
		t.Log(Utils().WebPush())
	})
}
