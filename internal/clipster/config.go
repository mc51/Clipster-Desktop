// Deal with reading, writing and loading configuration
package clipster

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync"

	"fyne.io/fyne/v2"
	"github.com/BurntSushi/toml"
	"golang.design/x/clipboard"
)

var CONFIG_PATHS []string
var CONFIG_FILEPATH string

const CONFIG_FILENAME = "config.toml"

const HOST_DEFAULT string = "https://clipster.cc"
const API_URI_COPY_PASTE = "/copy-paste/"
const API_URI_REGISTER = "/register/"
const API_URI_LOGIN = "/verify-user/"
const API_REQ_TIMEOUT = 6

const MAX_NOTIFICATION_LENGTH = 200
const MSG_NOTIFY_SHARED = "Your clip was shared"
const MSG_NOTIFY_GOT_CLIP = "The clip was copied to your clipboard"
const THUMBNAIL_HEIGHT = 200
const THUMBNAIL_WIDTH = 200
const DEFAULT_IMAGE_SAVE_NAME = "myclip.png"

// Must be same as the other clients
const HASH_ITERS_LOGIN = 20000
const HASH_ITERS_MSG = 10000
const HASH_LENGTH = 32

type Config struct {
	Server     string `toml:"server"`
	Username   string `toml:"username"`
	Hash_login string `toml:"hash_login"`
	Hash_msg   string `toml:"hash_msg"`
	// SHA-256 fingerprint of the server certificate the user has chosen to trust,
	// although it can not be verified (e.g. self signed). Empty if there is none.
	// Replaces the former option disable_ssl_cert_check
	Pinned_cert string `toml:"pinned_cert"`
}

var (
	// conf is read by flows running in the background, so only access it via getConf and setConf
	conf   Config
	confMu sync.RWMutex
)

// getConf returns a copy of the active config
func getConf() Config {
	confMu.RLock()
	defer confMu.RUnlock()
	return conf
}

// setConf makes c the active config
func setConf(c Config) {
	confMu.Lock()
	defer confMu.Unlock()
	conf = c
}

// Init prepares the config paths and the clipboard. It remembers the Fyne app, which is
// needed for notifications and windows. Must be called once on startup
func Init(a fyne.App) {
	fyneApp = a
	initConfigPaths()
	if err := clipboard.Init(); err != nil {
		log.Println("Error: clipboard not available", err)
	}
}

// LoadConfigFromFile reads the credentials from the config file and makes them the active config
func LoadConfigFromFile() (Config, error) {
	log.Println("Loading config file")
	debugf("Config file is %s", CONFIG_FILEPATH)
	var c Config
	md, err := toml.DecodeFile(CONFIG_FILEPATH, &c)
	if err != nil {
		return c, err
	}
	if md.IsDefined("disable_ssl_cert_check") {
		log.Println("Warning: disable_ssl_cert_check in the config is no longer supported and ignored. " +
			"A certificate that can not be verified must be trusted via Edit Credentials")
	}
	if c.Server == "" || c.Username == "" || c.Hash_login == "" || c.Hash_msg == "" {
		return c, errors.New("config file is incomplete: " + CONFIG_FILEPATH)
	}
	setConf(c)
	log.Println("Ok: loaded config")
	debugf("Loaded config for user %s on %s", c.Username, c.Server)
	return c, nil
}

// WriteConfigFile writes config struct to file. A temporary file is renamed over the
// config, so that it is never left half written and always only readable by the user
func WriteConfigFile(c Config) error {
	log.Println("Writing config")
	debugf("Writing config for user %s on %s", c.Username, c.Server)
	if err := os.MkdirAll(filepath.Dir(CONFIG_FILEPATH), 0700); err != nil {
		return err
	}
	// created with permissions 0600
	f, err := os.CreateTemp(filepath.Dir(CONFIG_FILEPATH), CONFIG_FILENAME+".*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name()) // fails harmlessly once renamed
	if err := toml.NewEncoder(f).Encode(c); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), CONFIG_FILEPATH)
}

// initConfigPaths sets CONFIG_FILEPATH to the config file in the first of the standard
// config folders that contains one. Otherwise it points into the user's config folder
func initConfigPaths() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Panicln("Error:", err)
	}
	CONFIG_PATHS = []string{
		filepath.Join(homedir, ".config", "clipster"),
		filepath.Join(homedir, ".clipster"),
		filepath.FromSlash("/etc/clipster"),
	}
	CONFIG_FILEPATH = findConfigFile(CONFIG_PATHS)
	debugf("Config file is %s", CONFIG_FILEPATH)
}

// findConfigFile returns the path of the config file in the first folder of dirs that
// contains one, or the path in the first folder if none does
func findConfigFile(dirs []string) string {
	for _, dir := range dirs {
		if path := filepath.Join(dir, CONFIG_FILENAME); fileExists(path) {
			return path
		}
	}
	return filepath.Join(dirs[0], CONFIG_FILENAME)
}

// fileExists checks if a file or folder exists
func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}
