package signer_test

import (
	"crypto/rand"
	"testing"

	"github.com/MikeRez0/ypmetrics/internal/logger"
	"github.com/MikeRez0/ypmetrics/internal/utils/signer"
	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	l := logger.GetLogger("debug")

	enc, err := signer.NewEncrypter(".tmp/pubkey.pem", l.Named("encrypter"))
	assert.NoError(t, err)
	dec, err := signer.NewDecrypter(".tmp/key.pem", l.Named("decrypter"))
	assert.NoError(t, err)

	t.Run("base", func(t *testing.T) {
		testData := []byte("MY SECRET DATA")

		env, err := enc.Encrypt(testData)
		assert.NoError(t, err)

		data, err := dec.Decrypt(env)
		assert.NoError(t, err)
		assert.NotEqual(t, testData, env.Data)

		assert.Equal(t, testData, data)
	})

	t.Run("fail key", func(t *testing.T) {
		testData := "MY SECRET DATA"

		env, err := enc.Encrypt([]byte(testData))
		assert.NoError(t, err)

		env.Key = []byte("BADKEY")

		_, err = dec.Decrypt(env)
		assert.Error(t, err)
	})

	t.Run("big data", func(t *testing.T) {
		testData := make([]byte, 1024*1024*1024)
		_, err := rand.Read(testData)
		assert.NoError(t, err)

		env, err := enc.Encrypt(testData)
		assert.NoError(t, err)

		data, err := dec.Decrypt(env)
		assert.NoError(t, err)
		assert.NotEqual(t, testData, env.Data)

		assert.Equal(t, testData, data)
	})
}

func TestSigner(t *testing.T) {
	s := signer.NewSigner("MYKEY")

	t.Run("json hash", func(t *testing.T) {
		data := struct {
			name string
			num  int
		}{name: "test", num: 1}

		val, err := s.GetHashJSON(data)
		assert.NoError(t, err)
		assert.True(t, s.ValidateJSON(data, val))

		val += "BAD"
		assert.False(t, s.ValidateJSON(data, val))
	})

	t.Run("byte array hash", func(t *testing.T) {
		data := []byte("MY DATA")

		val, err := s.GetHashBA(data)
		assert.NoError(t, err)
		assert.True(t, s.Validate(data, val))

		val += "BAD"
		assert.False(t, s.Validate(data, val))
	})
}
