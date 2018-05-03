package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
)

const (
	// keySize stores the key size to generate. This size is dictated by
	// https://golang.org/pkg/crypto/aes/#NewCipher
	keySize = 32
)

var (
	// ErrKeyFormat returns an error when the provided key is not valid
	ErrKeyFormat = errors.New("key format is invalid")
	// ErrDecryptionKeyMismatch returns an error when the provided key is not for the loaded key
	ErrDecryptionKeyMismatch = errors.New("decryption failure: current key not able to decrypt")
)

// Cipher stores the necessary encryption interface and latest key version
type Cipher struct {
	gcm cipher.AEAD
}

// NewCipher creates an AES-256 block cipher to encrypt/decrypt data
func NewCipher(password string) (*Cipher, error) {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	aesgcm, err := buildCipher(hasher.Sum(nil))
	if err != nil {
		return nil, ErrCryptoManager(err.Error())
	}
	return &Cipher{gcm: aesgcm}, nil
}

func buildCipher(keySlice []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(keySlice)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return aesgcm, nil
}

// Encrypt will take a plaintext byte array, encrypt it, and encode to a storageLayout
func (c *Cipher) Encrypt(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, c.gcm.NonceSize())
	// don't care about how many bytes are read back
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return []byte{}, ErrCryptoManager(err.Error())
	}
	ciphertext := c.gcm.Seal(nil, nonce, plaintext, nil)
	sl := storageLayout{
		nonce: nonce,
		value: ciphertext,
	}
	storageCiphertext, err := sl.Encode()
	if err != nil {
		return []byte{}, ErrCryptoManager(err.Error())
	}
	return storageCiphertext, nil
}

// Decrypt will take a []byte slice and storageLayout object and
func (c *Cipher) Decrypt(storageCiphertext []byte) ([]byte, error) {
	sl := new(storageLayout)
	err := sl.Decode(storageCiphertext)
	if err != nil {
		return []byte{}, ErrCryptoManager(err.Error())
	}
	plaintext, err := c.gcm.Open(nil, sl.Nonce(), sl.Value(), nil)
	if err != nil {
		return []byte{}, ErrCryptoManager(err.Error())
	}
	return plaintext, nil
}

// ErrCryptoManager is returned when a general context error for
type ErrCryptoManager string

// Error returns the formatted config error
func (ecm ErrCryptoManager) Error() string {
	return fmt.Sprintf("crypto manager failure: %s", string(ecm))
}
