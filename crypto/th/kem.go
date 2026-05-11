// Package th implements Thai hybrid cryptography.
package th

import (
	"crypto/sha256"
	"errors"
	"fmt"
)

// HybridKEM represents a combined Classical and PQC Key Encapsulation Mechanism.
type HybridKEM struct {
	ClassicalName string
	PQCName       string
}

// CombineKeys implements the concatenation or XOR-based key combination.
// According to NIST recommendations, HKDF-based concatenation is preferred.
func (h *HybridKEM) CombineKeys(classicalKey, pqcKey []byte) ([]byte, error) {
	if len(classicalKey) == 0 || len(pqcKey) == 0 {
		return nil, errors.New("input keys cannot be empty")
	}

	// Simple concatenation and hash for the hybrid shared secret
	// In a real implementation, this would follow a specific RFC for Hybrid KEM
	hsh := sha256.New()
	hsh.Write(classicalKey)
	hsh.Write(pqcKey)
	
	return hsh.Sum(nil), nil
}

// GetProfileDetails returns technical details for the selected profile.
func GetProfileDetails(profile string) string {
	switch profile {
	case "high-assurance":
		return "P-384 + ML-KEM-1024 (Security Level 5)"
	case "balanced":
		return "X25519 + ML-KEM-768 (Security Level 3)"
	case "backup":
		return "HQC-128 (Hamming Quasi-Cyclic - Code-based)"
	default:
		return "Unknown profile"
	}
}

// FalconReadiness represents the future support for compact signatures.
func FalconReadiness() string {
	return "อยู่ระหว่างการติดตามมาตรฐาน NIST (เตรียมพร้อมสำหรับ Compact Signatures)"
}

// LogGapAnalysis outputs remaining gaps to the system log (Thai).
func LogGapAnalysis() {
	fmt.Println("การวิเคราะห์ช่องว่าง (Gap Analysis):")
	fmt.Println("1. IANA Codepoints: ยังใช้ช่วง Private Use (0xEE42, 0xEE43)")
	fmt.Println("2. Fragmentation: ต้องระวังเรื่อง MTU สำหรับ SLH-DSA")
	fmt.Println("3. HW Acceleration: ยังไม่ได้ปรับแต่งประสิทธิภาพระดับ Hardware")
}
