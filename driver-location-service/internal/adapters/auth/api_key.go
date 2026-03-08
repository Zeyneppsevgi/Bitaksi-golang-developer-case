package auth

import "crypto/subtle"

type APIKeyVerifier struct {
	expected string
}

func NewAPIKeyVerifier(expected string) *APIKeyVerifier {
	return &APIKeyVerifier{expected: expected}
}

func (a *APIKeyVerifier) Verify(got string) bool {
	if got == "" || a.expected == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(got), []byte(a.expected)) == 1
}
