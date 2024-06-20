package mktls

import (
	"fmt"
	"bytes"
	"net"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"time"
	"math/big"
	"strings"
	"encoding/pem"
)

type TlsKey struct {
	Pk *rsa.PrivateKey
	PemBytes bytes.Buffer
}

type TlsCert struct {
	PemBytes bytes.Buffer
}

func StringFromBuffer(re *bytes.Buffer) string {
	return string(re.Bytes())
}

func (re *TlsKey) String() string {
	return StringFromBuffer(&re.PemBytes)
}

func (re *TlsKey) GenerateKey() *TlsKey {

	b := make([]byte, 16)
	n, e := rand.Read(b)
	fmt.Println(n, e, b)

	var err error
	re.Pk, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err.Error())
	}

	pbytes, err := x509.MarshalPKCS8PrivateKey(re.Pk)
	if err != nil {
		panic(err.Error())
	}

	blk := &pem.Block{Type: "PRIVATE KEY", Bytes: pbytes}
	err = pem.Encode(&re.PemBytes, blk)
	if err != nil {
		panic(err.Error())
	}

	return re
}

func (re *TlsKey) GenerateCertificate() *TlsCert {
	c := &TlsCert{}

	// assume RSA
	keyUsage := x509.KeyUsageDigitalSignature
	keyUsage |= x509.KeyUsageKeyEncipherment

	var notBefore time.Time
	notBefore = time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		panic(err.Error())
	}

	pkixName := pkix.Name{ Organization: []string{"Acme Co"} }
	eu := []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkixName,
		NotBefore: notBefore,
		NotAfter: notAfter,
		KeyUsage: keyUsage,
		ExtKeyUsage: eu,
		BasicConstraintsValid: true,
	}

	hosts := strings.Split("localhost, 127.0.0.1", ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	derBytes, err := x509.CreateCertificate(
		rand.Reader,
		&template,
		&template,
		&re.Pk.PublicKey,
		re.Pk,
	)
	if err != nil {
		panic(err.Error())
	}

	blk := &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}
	err = pem.Encode(&c.PemBytes, blk)
	if err != nil {
		panic(err.Error())
	}

	return c
}

func (re *TlsCert) String() string {
	return StringFromBuffer(&re.PemBytes)
}
