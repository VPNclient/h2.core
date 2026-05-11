// Package tls provides TLS constants for Thai cryptography.
package tls

import "crypto/tls"

const (
	// Hybrid KEM Named Groups (Private Use range)
	// X25519 + ML-KEM-768 (Balanced Profile)
	X25519_MLKEM768 tls.CurveID = 0xEE42

	// P-384 + ML-KEM-1024 (High-Assurance Profile)
	P384_MLKEM1024 tls.CurveID = 0xEE43
)

const (
	// Thai PQC Cipher Suite (Placeholder)
	TLS_TH_PQC_WITH_AES_256_GCM_SHA384 uint16 = 0xEA01
)
