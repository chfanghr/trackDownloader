package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
)

type authBuf struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

func (a *authBuf) createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (a *authBuf) Encrypt(password string) ([]byte, error) {
	buf, err := json.Marshal(a)
	if err != nil {
		return nil, errors.New("Encrypt authBuf : " + err.Error())
	}
	block, _ := aes.NewCipher([]byte(a.createHash(password)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.New("Encrypt authBuf : " + err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.New("Encrypt authBuf : " + err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, buf, nil)
	return ciphertext, nil
}

func (a *authBuf) Decrypt(password string, encrypted []byte) error {
	key := []byte(a.createHash(password))
	block, err := aes.NewCipher(key)
	if err != nil {
		return errors.New("Decrypt authbuf : " + err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return errors.New("Decrypt authbuf : " + err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	err = json.Unmarshal(plaintext, a)
	if err != nil {
		return errors.New("Decrypt authbuf : " + err.Error())
	}
	return nil
}

func (a *authBuf) GetPassword() string {
	return a.Password
}

func (a *authBuf) GetUsername() string {
	return a.Username
}

func (a *authBuf) SetPassword(p string) {
	a.Password = p
}

func (a *authBuf) SetUsername(u string) {
	a.Username = u
}
