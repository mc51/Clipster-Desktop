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

// UntrustedCertError means the certificate is neither valid nor the pinned one
type UntrustedCertError struct {
	Host        string
	Fingerprint string // of the presented certificate
	Pinned      string // trusted before, may be empty
	Reason      error  // why the normal check failed
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

// CertFingerprint returns the SHA-256 of a DER certificate as colon separated hex
func CertFingerprint(der []byte) string {
	sum := sha256.Sum256(der)
	h := strings.ToUpper(hex.EncodeToString(sum[:]))
	parts := make([]string, 0, len(sum))
	for i := 0; i < len(h); i += 2 {
		parts = append(parts, h[i:i+2])
	}
	return strings.Join(parts, ":")
}

// sameFingerprint ignores case, colons and spaces
func sameFingerprint(a, b string) bool {
	clean := func(s string) string {
		return strings.ToLower(strings.NewReplacer(":", "", " ", "").Replace(strings.TrimSpace(s)))
	}
	return clean(a) != "" && clean(a) == clean(b)
}

// verifyServerCert accepts a chain that is valid for host (roots nil: system CAs)
// or whose leaf has the fingerprint pin. Otherwise it returns *UntrustedCertError
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

// newTLSConfig verifies via verifyServerCert, before anything is sent to the server
func newTLSConfig(host, pin string) *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true, // replaced by VerifyConnection
		VerifyConnection: func(cs tls.ConnectionState) error {
			return verifyServerCert(host, pin, cs.PeerCertificates, nil)
		},
	}
}
