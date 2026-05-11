// Package th implements Thai hybrid signatures.
package th

import (
	"crypto"
	"errors"
	"io"
)

// HybridSignature represents a combined Classical (Ed25519/ECDSA) and ML-DSA signature.
type HybridSignature struct {
	ClassicalAlgo string // "Ed25519" or "ECDSA-P256"
	PQCAlgo       string // "ML-DSA-65"
}

// Sign generates a hybrid signature by concatenating classical and PQC signatures.
func (s *HybridSignature) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	// ในการใช้งานจริง จะต้องเรียกอัลกอริทึมทั้งสองตัว
	// Placeholder for hybrid signing logic
	return append([]byte("CLASSICAL_SIG:"), []byte("ML-DSA-65_SIG:")...), nil
}

// Verify checks the hybrid signature.
func (s *HybridSignature) Verify(publicKey, digest, signature []byte) error {
	if len(signature) == 0 {
		return errors.New("signature is empty")
	}
	// Verification logic would go here
	return nil
}

// GetSignatureUse Case returns the recommended usage in Thai.
func GetSignatureUseCase(algo string) string {
	switch algo {
	case "ML-DSA-65":
		return "ใบรับรองทั่วไป และ Control API (สมดุลระหว่างความเร็วและขนาด)"
	case "SLH-DSA":
		return "Root CA และ Firmware Signing (เน้นความมั่นคงปลอดภัยสูงสุด)"
	default:
		return "ไม่ทราบประเภท"
	}
}
