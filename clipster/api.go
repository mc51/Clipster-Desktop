package clipster

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func justPrint() {
	fmt.Println("Just Print!")
}

func GetLastClip() {
	// TODO
	fmt.Println("Get Last Clip")
}

func GetAllClips() {
	// TODO
	fmt.Println("Get all Clips")
}

func Register(host string, user string, pw string) {
	host, user, pw, err := AreCredsComplete(host, user, pw)
	if err != nil {

	} else {
		log.Println("Error:", err)
	}
	log.Println("Register:", host, user, pw)
}

func APILogin(host string, user string, hash_login string) error {
	// APILogin tries to auth against API endpoint using hash created from creds
	url := host + API_URI_LOGIN
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("Error:", err)
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth("m", hash_login)

	client := http.Client{Timeout: API_TIMEOUT * time.Second}
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
