package crypto

import "crypto/sha256"

func hashedSecret(secret string) []byte {
	hash := sha256.Sum256([]byte(secret))
	return hash[:]
}
