package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

type PasswordService interface {
	// Encrypt criptografa uma senha em texto plano e retorna a versão criptografada
	Encrypt(plainPassword string) (string, error)

	// Decrypt descriptografa uma senha criptografada e retorna o texto plano
	Decrypt(encryptedPassword string) (string, error)

	// Compare compara uma senha em texto plano com uma senha criptografada
	// Retorna true se corresponderem, false caso contrário
	Compare(plainPassword, encryptedPassword string) (bool, error)
}

type passwordService struct {
	key []byte
}

func NewPasswordService(key []byte) PasswordService {
	return &passwordService{
		key: key,
	}
}

func (p *passwordService) Encrypt(plainPassword string) (string, error) {
	block, err := aes.NewCipher(p.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plainPassword), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (p *passwordService) Decrypt(encryptedPassword string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedPassword)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(p.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("dados criptografados inválidos")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func (p *passwordService) Compare(plainPassword, encryptedPassword string) (bool, error) {
	decrypted, err := p.Decrypt(encryptedPassword)
	if err != nil {
		return false, err
	}

	return plainPassword == decrypted, nil
}
