// Trust on first use for server certificates
package clipster

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// UntrustedCertError is returned when the server presents a certificate that neither
// passes the normal verification, nor is the one the user has trusted before
type UntrustedCertError struct {
	Host        string // host name the connection was made to
	Fingerprint string // SHA-256 fingerprint of the certificate presented by the server
	Pinned      string // fingerprint the user has trusted before, empty if there is none
	Reason      error  // why the normal verification failed
}

func (e *UntrustedCertError) Error() string {
	if e.Pinned != "" {
		return fmt.Sprintf("The certificate of %s changed and is not the one you trusted! "+
			"The connection was refused. If you expect this, open Edit Credentials to trust the new one", e.Host)
	}
	return fmt.Sprintf("The certificate of %s is not trusted (%v). "+
		"Open Edit Credentials to review and trust it", e.Host, e.Reason)
}

func (e *UntrustedCertError) Unwrap() error { return e.Reason }

// CertFingerprint returns the SHA-256 fingerprint of a DER encoded certificate
// as upper case hex bytes separated by colons
func CertFingerprint(der []byte) string {
	sum := sha256.Sum256(der)
	h := strings.ToUpper(hex.EncodeToString(sum[:]))
	parts := make([]string, 0, len(sum))
	for i := 0; i < len(h); i += 2 {
		parts = append(parts, h[i:i+2])
	}
	return strings.Join(parts, ":")
}

// sameFingerprint compares two fingerprints, ignoring case, colons and whitespace
func sameFingerprint(a, b string) bool {
	clean := func(s string) string {
		return strings.ToLower(strings.NewReplacer(":", "", " ", "").Replace(strings.TrimSpace(s)))
	}
	return clean(a) != "" && clean(a) == clean(b)
}

// verifyServerCert accepts the certificate chain presented by the server if
//   - it is valid for host and signed by a trusted CA (roots, nil for the system's), or
//   - its leaf is exactly the certificate with the fingerprint pin (trust on first use).
//     Because the certificate itself is trusted, its expiry, issuer and host name are not checked.
//
// Anything else is returned as *UntrustedCertError
func verifyServerCert(host, pin string, certs []*x509.Certificate, roots *x509.CertPool) error {
	if len(certs) == 0 {
		return errors.New("the server did not send a certificate")
	}
	leaf := certs[0]
	intermediates := x509.NewCertPool()
	for _, c := range certs[1:] {
		intermediates.AddCert(c)
	}
	_, err := leaf.Verify(x509.VerifyOptions{DNSName: host, Roots: roots, Intermediates: intermediates})
	if err == nil {
		return nil
	}
	fp := CertFingerprint(leaf.Raw)
	if sameFingerprint(pin, fp) {
		return nil
	}
	return &UntrustedCertError{Host: host, Fingerprint: fp, Pinned: strings.TrimSpace(pin), Reason: err}
}

// newTLSConfig returns the TLS config for connecting to host. The certificate is checked
// by verifyServerCert. This is done in VerifyConnection, so that it also applies to resumed
// sessions. Go's own verification is switched off, as it would reject self signed certificates
// before the pin could be looked at. Nothing (not even the credentials) is sent to the
// server before the handshake has succeeded
func newTLSConfig(host, pin string) *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true, // only because VerifyConnection does the verification
		VerifyConnection: func(cs tls.ConnectionState) error {
			return verifyServerCert(host, pin, cs.PeerCertificates, nil)
		},
	}
}
