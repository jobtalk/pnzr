package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"io"
)

type Secret struct {
	plain   []byte
	Chipher []byte `json:"chipher"`
	Mac     []byte `json:"mac"`
}

func New(p []byte) *Secret {
	return &Secret{plain: p}
}

func NewString(s string) *Secret {
	return &Secret{plain: []byte(s)}
}

func (s *Secret) Encrypt(key string) error {
	keysha := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(keysha[0:len(keysha)])
	if err != nil {
		return err
	}
	cipherText := make([]byte, aes.BlockSize+len(s.plain))
	iv := cipherText[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return err
	}
	encryptStream := cipher.NewCTR(block, iv)
	encryptStream.XORKeyStream(cipherText[aes.BlockSize:], s.plain)
	s.Chipher = cipherText

	h := hmac.New(sha512.New, []byte(key))
	if _, err := h.Write(s.plain); err != nil {
		return err
	}
	s.Mac = h.Sum(nil)
	return nil
}

func (s *Secret) Decrypt(key string) ([]byte, error) {
	keysha := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(keysha[0:len(keysha)])
	if err != nil {
		return nil, err
	}
	decryptedText := make([]byte, len(s.Chipher[aes.BlockSize:]))
	decryptStream := cipher.NewCTR(block, s.Chipher[:aes.BlockSize])
	decryptStream.XORKeyStream(decryptedText, s.Chipher[aes.BlockSize:])

	h := hmac.New(sha512.New, []byte(key))
	if _, err := h.Write(decryptedText); err != nil {
		return nil, err
	}
	if !hmac.Equal(s.Mac, h.Sum(nil)) {
		return nil, errors.New("message authority code is not match")
	}
	return decryptedText, nil
}
