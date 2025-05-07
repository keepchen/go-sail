package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMergeStandardClaims(t *testing.T) {
	t.Run("MergeStandardClaims", func(t *testing.T) {
		fields := map[string]interface{}{
			"uid": "1000",
			"acc": "account-1000",
			"jti": "AAA-BBB-CCC-DDD-EEE",
			"iss": "A Inc.",
		}
		t.Log(MergeStandardClaims(fields))

		t.Log(MergeStandardClaims(nil))
	})
}

func TestGetExpirationTime(t *testing.T) {
	t.Run("GetExpirationTime", func(t *testing.T) {
		var mp MapClaims = map[string]interface{}{
			"exp": time.Now().Add(time.Hour).Unix(),
		}
		exp, err := mp.GetExpirationTime()
		assert.NotNil(t, exp)
		assert.NoError(t, err)
	})

	t.Run("GetExpirationTime-Failure", func(t *testing.T) {
		var mp MapClaims = map[string]interface{}{}
		exp, err := mp.GetExpirationTime()
		assert.Nil(t, exp)
		assert.Error(t, err)
	})
}

func TestGetIssuedAt(t *testing.T) {
	t.Run("GetIssuedAt", func(t *testing.T) {
		var mp MapClaims = map[string]interface{}{
			"iat": time.Now().Unix(),
		}
		iat, err := mp.GetIssuedAt()
		assert.NotNil(t, iat)
		assert.NoError(t, err)
	})

	t.Run("GetIssuedAt-Failure", func(t *testing.T) {
		var mp MapClaims = map[string]interface{}{}
		iat, err := mp.GetIssuedAt()
		assert.Nil(t, iat)
		assert.Error(t, err)
	})
}

func TestGetNotBefore(t *testing.T) {
	t.Run("GetNotBefore", func(t *testing.T) {
		var mp MapClaims = map[string]interface{}{
			"nbf": time.Now().Unix(),
		}
		nbf, err := mp.GetNotBefore()
		assert.NotNil(t, nbf)
		assert.NoError(t, err)
	})

	t.Run("GetNotBefore-Failure", func(t *testing.T) {
		var mp MapClaims = map[string]interface{}{}
		nbf, err := mp.GetNotBefore()
		assert.Nil(t, nbf)
		assert.Error(t, err)
	})
}

func TestGetIssuer(t *testing.T) {
	t.Run("GetIssuer", func(t *testing.T) {
		var mp MapClaims = map[string]interface{}{
			"iss": defaultTokenIssuer,
		}
		iss, err := mp.GetIssuer()
		assert.NotNil(t, iss)
		assert.NoError(t, err)
	})

	t.Run("GetIssuer-Failure", func(t *testing.T) {
		var mp MapClaims = map[string]interface{}{}
		iss, err := mp.GetIssuer()
		assert.Equal(t, true, len(iss) == 0)
		assert.Error(t, err)
	})
}

func TestGetSubject(t *testing.T) {
	t.Run("GetSubject", func(t *testing.T) {
		var mp MapClaims = map[string]interface{}{
			"sub": defaultTokenIssuer,
		}
		sub, err := mp.GetSubject()
		assert.NotNil(t, sub)
		assert.NoError(t, err)
	})

	t.Run("GetSubject-Failure", func(t *testing.T) {
		var mp MapClaims = map[string]interface{}{}
		sub, err := mp.GetSubject()
		assert.Equal(t, true, len(sub) == 0)
		assert.Error(t, err)
	})
}

func TestGetAudience(t *testing.T) {
	t.Run("GetAudience", func(t *testing.T) {
		var mp MapClaims = map[string]interface{}{
			"aud": defaultTokenIssuer,
		}
		aud, err := mp.GetAudience()
		assert.NotNil(t, aud)
		assert.NoError(t, err)
	})

	t.Run("GetAudience-Failure", func(t *testing.T) {
		var mp MapClaims = map[string]interface{}{}
		aud, err := mp.GetAudience()
		assert.Equal(t, true, len(aud[0]) == 0)
		assert.Error(t, err)
	})
}
