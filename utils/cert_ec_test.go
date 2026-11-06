package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReportKeyWhetherMatchECPrivateKey(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	template := &x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	require.NoError(t, err)
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	otherKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	for _, tc := range []struct {
		name  string
		key   *ecdsa.PrivateKey
		match bool
	}{
		{name: "matching", key: key, match: true},
		{name: "mismatched", key: otherKey, match: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			keyDER, err := x509.MarshalECPrivateKey(tc.key)
			require.NoError(t, err)
			keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
			match, err := Cert().ReportKeyWhetherMatch(certPEM, keyPEM)
			if tc.match {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
			require.Equal(t, tc.match, match)
		})
	}
}
