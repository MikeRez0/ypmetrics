// Program generates key and certificate.
package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"log"
	"os"
)

func main() {
	dir := flag.String("d", "", "output directory with /")
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatalf("error generating key: %v", err)
	}

	var pubKeyPEM bytes.Buffer
	err = pem.Encode(&pubKeyPEM, &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	})
	if err != nil {
		log.Fatalf("error writing certificate: %v", err)
	}

	var privateKeyPEM bytes.Buffer
	err = pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		log.Fatalf("error writing private key: %v", err)
	}

	fc, err := os.Create(*dir + "pubkey.pem")
	if err != nil {
		log.Fatalf("error creating file: %v", err)
	}
	defer fc.Close()
	fc.Write(pubKeyPEM.Bytes())

	fk, err := os.Create(*dir + "key.pem")
	if err != nil {
		log.Fatalf("error creating file: %v", err)
	}
	defer fk.Close()
	fk.Write(privateKeyPEM.Bytes())
}
