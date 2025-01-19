// Program generates key and certificate.
package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dir := flag.String("d", "", "output directory with /")
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("error generating key: %w", err)
	}

	var pubKeyPEM bytes.Buffer
	err = pem.Encode(&pubKeyPEM, &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	})
	if err != nil {
		return fmt.Errorf("error writing certificate: %w", err)
	}

	var privateKeyPEM bytes.Buffer
	err = pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		return fmt.Errorf("error writing private key: %w", err)
	}

	fc, err := os.Create(*dir + "pubkey.pem")
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer func() { _ = fc.Close() }()
	_, err = fc.Write(pubKeyPEM.Bytes())
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	fk, err := os.Create(*dir + "key.pem")
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer func() { _ = fc.Close() }()
	_, err = fk.Write(privateKeyPEM.Bytes())
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}
	return nil
}
