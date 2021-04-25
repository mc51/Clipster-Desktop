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
)

type Clips struct {
	Id         int
	User       string
	Text       string
	Device     string
	Created_at string
}

func APIShareClip(clip string) error {
	// APIShareClip sends encrypted Clip to API endpoint for sharing
	url := conf.Server + API_URI_COPY_PASTE
	payload, err := json.Marshal(map[string]string{
		"text":   clip,
		"device": "desktop",
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
	req.SetBasicAuth("m", conf.Hash_login)

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

func APIDownloadAllClips() ([]Clips, error) {
	// APIDownloadAllClips retrieves all encrypted Clips from server and returns them as Clips struct
	var clips []Clips
	url := conf.Server + API_URI_COPY_PASTE
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth("m", conf.Hash_login)

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

func APIRegister(host string, user string, hash_login string, ssl_disable bool) error {
	// APIRegister registers new account at API endpoint using hash created from creds
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
		log.Println("Error: registration failed", string(body))
		return errors.New("registration failed: " + string(body))
	}
}

func APILogin(host string, user string, hash_login string, ssl_disable bool) error {
	// APILogin authenticates against API endpoint using hash created from creds
	url := host + API_URI_LOGIN
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("Error:", err)
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth("m", hash_login)

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
		log.Println("Error: login failed", string(body))
		return errors.New("login failed: " + string(body))
	}
}
