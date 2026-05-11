package th

import (
	"crypto/tls"
	"testing"

	"github.com/nativemind/https-vpn/crypto"
	thtls "github.com/nativemind/https-vpn/crypto/th/tls"
)

func TestProviderRegistration(t *testing.T) {
	p, ok := crypto.Get("th")
	if !ok {
		t.Fatal("Thai provider not registered")
	}
	if p.Name() != "th" {
		t.Errorf("Expected name 'th', got '%s'", p.Name())
	}
}

func TestConfigureTLS(t *testing.T) {
	p := &Provider{Profile: "balanced"}
	cfg := &tls.Config{}
	err := p.ConfigureTLS(cfg)
	if err != nil {
		t.Fatalf("ConfigureTLS failed: %v", err)
	}

	if cfg.MinVersion != tls.VersionTLS13 {
		t.Errorf("Expected TLS 1.3, got %v", cfg.MinVersion)
	}

	foundHybrid := false
	for _, curve := range cfg.CurvePreferences {
		if curve == thtls.X25519_MLKEM768 {
			foundHybrid = true
			break
		}
	}
	if !foundHybrid {
		t.Error("X25519_MLKEM768 not found in curve preferences")
	}
}

func TestHybridKeyCombination(t *testing.T) {
	hkem := &HybridKEM{}
	k1 := []byte("classical_secret")
	k2 := []byte("pqc_secret")
	combined, err := hkem.CombineKeys(k1, k2)
	if err != nil {
		t.Fatalf("CombineKeys failed: %v", err)
	}
	if len(combined) != 32 { // SHA-256 result
		t.Errorf("Expected 32 bytes, got %d", len(combined))
	}
}

func TestSLHDSASigning(t *testing.T) {
	slh := &SLHDSAProvider{}
	data := []byte("firmware_data")
	sig, err := slh.SignOffline(data)
	if err != nil {
		t.Fatalf("SignOffline failed: %v", err)
	}
	if len(sig) <= len(data) {
		t.Error("Expected signature to be larger than data")
	}
}
