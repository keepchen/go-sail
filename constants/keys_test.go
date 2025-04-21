package constants

import "testing"

func TestPrintKeys(t *testing.T) {
	t.Log(PublicKeyBeginStr)
	t.Log(PublicKeyEndStr)
	t.Log(PrivateKeyBeginStr)
	t.Log(PrivateKeyEndStr)
}
