// API calls to Server
package clipster

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Clips struct {
	Id            int
	User          string
	Text          string
	Format        string
	Device        string
	Created_at    string
	TextDecrypted string `json:"-"`
	ImageBytes    []byte `json:"-"`
	ThumbBytes    []byte `json:"-"`
}

// apiRequest sends a JSON request to the API endpoint and returns the response body.
// If user is not empty, basic auth is used. Non 2xx/3xx status codes are returned as error.
// Certificates are checked as described at verifyServerCert, pin is the fingerprint of the
// certificate the user trusts (can be empty). A certificate that is not trusted is returned as *UntrustedCertError
func apiRequest(method string, url string, payload any, user string, hash_login string,
	pin string) ([]byte, error) {
	var reqBody io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	if user != "" {
		req.SetBasicAuth(user, hash_login)
	}

	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = newTLSConfig(req.URL.Hostname(), pin)
	client := http.Client{Timeout: API_REQ_TIMEOUT * time.Second, Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		var certErr *UntrustedCertError
		if errors.As(err, &certErr) {
			return nil, certErr // without the url wrapped around it
		}
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		msg := strings.TrimSpace(string(body))
		if msg == "" {
			msg = resp.Status
		}
		return nil, errors.New(msg)
	}
	return body, nil
}

// APIShareClip sends encrypted Clip to API endpoint for sharing
func APIShareClip(clip string, format string) error {
	payload := map[string]string{
		"text":   clip,
		"device": "desktop",
		"format": format,
	}
	c := getConf()
	_, err := apiRequest(http.MethodPost, c.Server+API_URI_COPY_PASTE, payload,
		c.Username, c.Hash_login, c.Pinned_cert)
	if err != nil {
		log.Println("Error: sharing clip failed", err)
		return fmt.Errorf("sharing clip failed: %w", err)
	}
	log.Println("Ok: sharing clip successful")
	return nil
}

// APIDownloadAllClips retrieves all encrypted Clips from server and returns them as Clips struct
func APIDownloadAllClips() ([]Clips, error) {
	var clips []Clips
	c := getConf()
	body, err := apiRequest(http.MethodGet, c.Server+API_URI_COPY_PASTE, nil,
		c.Username, c.Hash_login, c.Pinned_cert)
	if err != nil {
		log.Println("Error: download Clips", err)
		return nil, fmt.Errorf("download Clips failed: %w", err)
	}
	if err := json.Unmarshal(body, &clips); err != nil {
		log.Println("Error:", err)
		return nil, err
	}
	log.Println("Ok: downloaded Clips:", len(clips))
	return clips, nil
}

// APIRegister registers new account at API endpoint using hash created from creds
func APIRegister(host string, user string, hash_login string, pin string) error {
	payload := map[string]string{
		"username": user,
		"password": hash_login,
	}
	if _, err := apiRequest(http.MethodPost, host+API_URI_REGISTER, payload, "", "",
		pin); err != nil {
		return fmt.Errorf("registration failed: %w", err)
	}
	log.Println("Ok: registration successful")
	return nil
}

// APILogin authenticates against API endpoint using hash created from creds
func APILogin(host string, user string, hash_login string, pin string) error {
	if _, err := apiRequest(http.MethodGet, host+API_URI_LOGIN, nil, user, hash_login,
		pin); err != nil {
		return fmt.Errorf("login failed: %w", err)
	}
	log.Println("Ok: logged in")
	return nil
}
