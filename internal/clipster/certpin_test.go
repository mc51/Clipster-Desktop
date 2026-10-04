package clipster

import (
	"crypto/x509"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
)

// newTLSTestServer starts an untrusted HTTPS server and counts its requests
func newTLSTestServer(t *testing.T) (*httptest.Server, *atomic.Int32) {
	t.Helper()
	var requests atomic.Int32
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		w.Write([]byte("ok"))
	}))
	t.Cleanup(srv.Close)
	return srv, &requests
}

func TestCertFingerprint(t *testing.T) {
	// SHA-256 of the empty input
	want := "E3:B0:C4:42:98:FC:1C:14:9A:FB:F4:C8:99:6F:B9:24:27:AE:41:E4:64:9B:93:4C:A4:95:99:1B:78:52:B8:55"
	if got := CertFingerprint(nil); got != want {
		t.Errorf("CertFingerprint = %q, want %q", got, want)
	}
}

func TestSameFingerprint(t *testing.T) {
	fp := "E3:B0:C4:42"
	for _, other := range []string{fp, "e3:b0:c4:42", "E3B0C442", " e3 b0 c4 42\n"} {
		if !sameFingerprint(fp, other) {
			t.Errorf("%q should equal %q", other, fp)
		}
	}
	for _, other := range []string{"", "E3:B0:C4:43", "E3:B0:C4"} {
		if sameFingerprint(fp, other) {
			t.Errorf("%q should not equal %q", other, fp)
		}
	}
	if sameFingerprint("", "") {
		t.Error("an empty pin must never match")
	}
}

func TestAPIRequestCertificate(t *testing.T) {
	srv, requests := newTLSTestServer(t)
	fp := CertFingerprint(srv.Certificate().Raw)
	request := func(pin string) ([]byte, error) {
		return apiRequest(http.MethodGet, srv.URL, nil, "alice", "secret-hash", pin)
	}

	t.Run("unknown certificate is refused", func(t *testing.T) {
		before := requests.Load()
		_, err := request("")
		var certErr *UntrustedCertError
		if !errors.As(err, &certErr) {
			t.Fatalf("got error %v, want *UntrustedCertError", err)
		}
		if certErr.Fingerprint != fp || certErr.Pinned != "" || certErr.Reason == nil {
			t.Errorf("got %+v, want fingerprint %s without pin and a reason", certErr, fp)
		}
		if requests.Load() != before {
			t.Error("the request, including the credentials, reached a server that is not trusted")
		}
	})

	t.Run("pinned certificate is accepted", func(t *testing.T) {
		for _, pin := range []string{fp, strings.ToLower(fp), strings.ReplaceAll(fp, ":", "")} {
			if body, err := request(pin); err != nil || string(body) != "ok" {
				t.Errorf("pin %q: got %q, %v", pin, body, err)
			}
		}
	})

	t.Run("other certificate than the pinned one is refused", func(t *testing.T) {
		before := requests.Load()
		pin := "00:11:22"
		_, err := request(pin)
		var certErr *UntrustedCertError
		if !errors.As(err, &certErr) {
			t.Fatalf("got error %v, want *UntrustedCertError", err)
		}
		if certErr.Pinned != pin || certErr.Fingerprint != fp {
			t.Errorf("got %+v, want pin %s and fingerprint %s", certErr, pin, fp)
		}
		if !strings.Contains(certErr.Error(), "changed") {
			t.Errorf("message %q should tell that the certificate changed", certErr)
		}
		if requests.Load() != before {
			t.Error("the request, including the credentials, reached a server with a changed certificate")
		}
	})

	t.Run("wrapped error can be detected", func(t *testing.T) {
		err := APILogin(srv.URL, "alice", "secret-hash", "")
		var certErr *UntrustedCertError
		if !errors.As(err, &certErr) || !strings.HasPrefix(err.Error(), "login failed: ") {
			t.Errorf("got %v", err)
		}
	})
}

func TestVerifyServerCert(t *testing.T) {
	srv, _ := newTLSTestServer(t)
	cert := srv.Certificate()
	fp := CertFingerprint(cert.Raw)
	chain := []*x509.Certificate{cert}
	roots := x509.NewCertPool()
	roots.AddCert(cert)

	tests := []struct {
		name  string
		host  string
		pin   string
		roots *x509.CertPool
		ok    bool
	}{
		{"trusted CA and matching host", "127.0.0.1", "", roots, true},
		{"trusted CA, pin is not needed", "127.0.0.1", "00:11", roots, true},
		{"trusted CA but wrong host", "example.org", "", roots, false},
		{"pin matches, no CA", "127.0.0.1", fp, nil, true},
		{"pin matches, host name does not", "example.org", fp, nil, true},
		{"nothing matches", "127.0.0.1", "00:11", nil, false},
		{"no pin, no CA", "127.0.0.1", "", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifyServerCert(tt.host, tt.pin, chain, tt.roots)
			if (err == nil) != tt.ok {
				t.Errorf("got %v, want ok=%v", err, tt.ok)
			}
		})
	}
	if err := verifyServerCert("127.0.0.1", fp, nil, nil); err == nil {
		t.Error("a missing certificate must not be accepted")
	}
}

// The old disable_ssl_cert_check must be ignored without breaking the load
func TestLoadConfigFromFileLegacySSLOption(t *testing.T) {
	oldPath := CONFIG_FILEPATH
	CONFIG_FILEPATH = filepath.Join(t.TempDir(), CONFIG_FILENAME)
	t.Cleanup(func() {
		CONFIG_FILEPATH = oldPath
		setConf(Config{})
	})
	content := "server = \"https://example.com\"\nusername = \"a\"\nhash_login = \"l\"\nhash_msg = \"m\"\n" +
		"disable_ssl_cert_check = true\n"
	if err := os.WriteFile(CONFIG_FILEPATH, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	c, err := LoadConfigFromFile()
	if err != nil {
		t.Fatal(err)
	}
	if c.Pinned_cert != "" {
		t.Errorf("legacy option must not pin a certificate, got %q", c.Pinned_cert)
	}
}
