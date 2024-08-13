package jwt

import "testing"

func TestMergeStandardClaims(t *testing.T) {
	fields := map[string]interface{}{
		"uid": "1000",
		"acc": "account-1000",
		"jti": "AAA-BBB-CCC-DDD-EEE",
		"iss": "A Inc.",
	}
	t.Log(MergeStandardClaims(fields))

	t.Log(MergeStandardClaims(nil))
}
