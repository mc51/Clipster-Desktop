// API calls to Server
package clipster

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gotk3/gotk3/gtk"
)

type Clips struct {
	Id            int
	User          string
	Text          string
	Format        string
	Device        string
	Created_at    string
	TextDecrypted string
	ImageBytes    []byte
	GtkThumb      *gtk.Image
}

// APIShareClip sends encrypted Clip to API endpoint for sharing
func APIShareClip(clip string, format string) error {
	url := conf.Server + API_URI_COPY_PASTE
	payload, err := json.Marshal(map[string]string{
		"text":   clip,
		"device": "desktop",
		"format": format,
	})
	if err != nil {
		log.Println("Error", err)
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(conf.Username, conf.Hash_login)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: conf.Disable_ssl_cert_check},
	}
	client := http.Client{Timeout: API_REQ_TIMEOUT * time.Second, Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		log.Println("Ok: sharing clip successfull")
		return nil
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error:", err)
			return err
		}
		log.Println("Error: sharing clip failed", string(body))
		return errors.New("Sharing clip failed: " + string(body))
	}
}

// APIDownloadAllClips retrieves all encrypted Clips from server and returns them as Clips struct
func APIDownloadAllClips() ([]Clips, error) {
	var clips []Clips
	url := conf.Server + API_URI_COPY_PASTE
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(conf.Username, conf.Hash_login)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: conf.Disable_ssl_cert_check},
	}
	client := http.Client{Timeout: API_REQ_TIMEOUT * time.Second, Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		log.Println("Ok: Download Clips")
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error:", err)
			return nil, err
		}
		log.Println("Ok: Download Clips", string(body))
		err = json.Unmarshal(body, &clips)
		if err != nil {
			log.Println("Error:", err)
			return nil, err
		}
		return clips, nil
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error:", err)
			return nil, err
		}
		log.Println("Error: download Clips", string(body))
		return nil, errors.New("download Clips failed: " + string(body))
	}
}

// APIRegister registers new account at API endpoint using hash created from creds
func APIRegister(host string, user string, hash_login string, ssl_disable bool) error {
	url := host + API_URI_REGISTER
	payload, err := json.Marshal(map[string]string{
		"username": user,
		"password": hash_login,
	})
	if err != nil {
		log.Println("Error", err)
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: ssl_disable},
	}
	client := http.Client{Timeout: API_REQ_TIMEOUT * time.Second, Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		log.Println("Ok: registration successfull")
		return nil
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error:", err)
			return err
		}
		return errors.New("registration failed: " + string(body))
	}
}

// APILogin authenticates against API endpoint using hash created from creds
func APILogin(host string, user string, hash_login string, ssl_disable bool) error {
	url := host + API_URI_LOGIN
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("Error:", err)
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(user, hash_login)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: ssl_disable},
	}
	client := http.Client{Timeout: API_REQ_TIMEOUT * time.Second, Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		log.Println("Ok: logged in")
		return nil
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error:", err)
			return err
		}
		return errors.New("login failed: " + string(body))
	}
}
