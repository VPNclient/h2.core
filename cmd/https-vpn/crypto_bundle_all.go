//go:build !bundle_custom

package main

// Default binary includes all providers/ciphersuites.
// Providers register themselves via init() on import.
import (
	_ "github.com/vpnclient/https-vpn/crypto/cn"
	_ "github.com/vpnclient/https-vpn/crypto/fr"
	_ "github.com/vpnclient/https-vpn/crypto/ko"
	_ "github.com/vpnclient/https-vpn/crypto/ru"
	_ "github.com/vpnclient/https-vpn/crypto/th"
	_ "github.com/vpnclient/https-vpn/crypto/ua"
	_ "github.com/vpnclient/https-vpn/crypto/uk"
	_ "github.com/vpnclient/https-vpn/crypto/us"
)

