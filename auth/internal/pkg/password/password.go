package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type PassEncryptParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func NewDefaultPassEncrypt() *PassEncryptParams {
	return &PassEncryptParams{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  32,
		keyLength:   32,
	}
}

func (h *PassEncryptParams) EncryptPassword(passwd string) (string, error) {
	salt, err := generateRandomBytes(h.saltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(passwd), salt, h.iterations, h.memory, h.parallelism, h.keyLength)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation.
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, h.memory, h.iterations, h.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func (h *PassEncryptParams) ValidatePass(pass, encryptedPass string) (bool, error) {
	p, salt, hash, err := h.decodeHash(encryptedPass)
	if err != nil {
		return false, fmt.Errorf("cant validate password")
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(pass), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (h *PassEncryptParams) decodeHash(encodedHash string) (*PassEncryptParams, []byte, []byte, error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, fmt.Errorf("invalid hash")
	}

	var version int
	_, err := fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, fmt.Errorf("invalid argon version")
	}

	p := &PassEncryptParams{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	hash, err := base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}
