package mktls

import (
	"testing"
)

func TestMktls(t *testing.T) {
	k := (&TlsKey{}).GenerateKey()
	c := k.GenerateCertificate()

	t.Logf("Key: %s", k)
	t.Logf("Cert: %s", c)

	// TODO: test that they match
}
