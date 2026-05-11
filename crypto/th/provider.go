// Package th implements the Thai cryptographic provider for HTTPS VPN.
// It supports Hybrid Post-Quantum Cryptography (FIPS 203, 204, 205).
package th

import (
	"crypto/tls"
	"fmt"

	"github.com/nativemind/https-vpn/crypto"
	thtls "github.com/nativemind/https-vpn/crypto/th/tls"
)

// Provider implements crypto.Provider for Thai cryptography.
type Provider struct {
	Profile string // "balanced" or "high-assurance"
}

func init() {
	// Register the default balanced provider
	crypto.Register(&Provider{Profile: "balanced"})
}

// Name returns the provider identifier.
func (p *Provider) Name() string {
	return "th"
}

// ConfigureTLS applies Thai PQC settings to tls.Config.
// This ensures compatibility with the HTTPS VPN architecture.
func (p *Provider) ConfigureTLS(cfg *tls.Config) error {
	// HTTPS VPN uses TLS 1.3 as the foundation
	cfg.MinVersion = tls.VersionTLS13
	cfg.MaxVersion = tls.VersionTLS13

	// Configure Hybrid KEM based on profile
	switch p.Profile {
	case "high-assurance":
		cfg.CurvePreferences = []tls.CurveID{
			thtls.P384_MLKEM1024,
			tls.X25519,
		}
	case "balanced":
		fallthrough
	default:
		cfg.CurvePreferences = []tls.CurveID{
			thtls.X25519_MLKEM768,
			tls.X25519,
		}
	}

	// Default cipher suites (PQC-ready)
	cfg.CipherSuites = p.SupportedCipherSuites()

	return nil
}

// SupportedCipherSuites returns the list of supported cipher suite IDs.
func (p *Provider) SupportedCipherSuites() []uint16 {
	return []uint16{
		thtls.TLS_TH_PQC_WITH_AES_256_GCM_SHA384,
		tls.TLS_AES_256_GCM_SHA384,
		tls.TLS_CHACHA20_POLY1305_SHA256,
	}
}

// Description returns the provider description in Thai.
func (p *Provider) Description() string {
	return fmt.Sprintf("คริปโตกราฟีไทยแบบไฮบริด (%s profile)", p.Profile)
}

// Algorithms returns the list of supported algorithms.
func (p *Provider) Algorithms() []string {
	return []string{
		"ML-KEM-768/1024 (FIPS 203)",
		"ML-DSA-65 (FIPS 204)",
		"SLH-DSA (FIPS 205)",
		"X25519/P-384 (Hybrid Classical)",
		"HQC (Backup KEM)",
	}
}
