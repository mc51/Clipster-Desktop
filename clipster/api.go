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

func Register(host, user, pw) {
	// Register created new account on endpoint using creds
	// Same as login but POST using payload user and login_hash
}

func Login(host string, user string, pw string) error {
	// Login authenticates against login endpoint using basic auth
	url := host + API_URI_LOGIN
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth("m", "5wio7yYuSZuRPQr6Fk65r-zamVUDFFn_nd6Tg7KGWdc=")

	client := http.Client{Timeout: API_TIMEOUT * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Println("Ok: logged in")
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error:", err)
			return err
		}
		log.Println("Error: login failed", string(body))
		return errors.New("login failed: " + string(body))
	}
	return nil
}
