package signer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
)

type Signer struct {
	h hash.Hash
}

func NewSigner(key string) *Signer {
	return &Signer{h: hmac.New(sha256.New, []byte(key))}
}

func (s *Signer) WriteJSON(data any) (int, error) {
	ba, err := json.Marshal(data)
	if err != nil {
		return -1, fmt.Errorf("json error: %w", err)
	}
	return s.h.Write(ba) //nolint //all ok
}

func (s *Signer) Write(p []byte) (int, error) {
	return s.h.Write(p) //nolint //all ok
}

func (s *Signer) GetHash() string {
	h := s.h.Sum(nil)
	return base64.StdEncoding.EncodeToString(h) // hex.EncodeToString(h)
}

func (s *Signer) GetHashJSON(data any) (string, error) {
	s.Reset()

	_, err := s.WriteJSON(data)
	if err != nil {
		return "", err
	}

	return s.GetHash(), nil
}
func (s *Signer) GetHashBA(data []byte) (string, error) {
	s.Reset()

	_, err := s.Write(data)
	if err != nil {
		return "", err
	}

	return s.GetHash(), nil
}

func (s *Signer) ValidateJSON(data any, h string) bool {
	exp, err := s.GetHashJSON(data)
	if err != nil {
		return false
	}
	return exp == h
}

func (s *Signer) Reset() {
	s.h.Reset()
}
