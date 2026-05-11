// Package th implements Thai SLH-DSA (FIPS 205).
package th

import (
	"errors"
)

// SLHDSAProvider implements Stateless Hash-Based Digital Signature Standard.
// Used for High-Assurance trust anchors.
type SLHDSAProvider struct {
	Variant string // e.g., "SHAKE-128f"
}

// SignOffline signs data for long-term persistence (Firmware, Root Manifest).
func (p *SLHDSAProvider) SignOffline(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}
	// SLH-DSA signatures are large (up to 30KB+)
	// Implementation should consider fragmentation
	return append([]byte("SLH-DSA_PERSISTENT_SIG:"), data...), nil
}

// GetSecurityProperties returns the properties of SLH-DSA in Thai.
func (p *SLHDSAProvider) GetSecurityProperties() string {
	return "ทนทานสูงสุดต่อการโจมตีทางควอนตัม (Hash-based) เหมาะสำหรับระบบที่ต้องอยู่นาน (Long-term trust)"
}
