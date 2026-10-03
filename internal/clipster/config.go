// Deal with reading, writing and loading configuration
package clipster

import (
	"log"
	"os"
	"path/filepath"
	"reflect"

	"fyne.io/fyne/v2"
	"github.com/spf13/viper"
	"golang.design/x/clipboard"
)

var CONFIG_PATHS []string
var CONFIG_FILEPATH string

const CONFIG_FILENAME = "config.yaml"
const CONFIG_TYPE = "yaml"

const HOST_DEFAULT string = "https://clipster.cc"
const RE_HOSTNAME string = `^(https):\/\/[^\s\/$.?#].[^\s]*|://localhost:|://127.0.0.1:|^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`

const API_URI_COPY_PASTE = "/copy-paste/"
const API_URI_REGISTER = "/register/"
const API_URI_LOGIN = "/verify-user/"
const API_REQ_TIMEOUT = 6

const MAX_NOTIFICATION_LENGTH = 200
const MSG_NOTIFY_GOT_IMAGE = "Got an image!"
const THUMBNAIL_HEIGHT = 200
const THUMBNAIL_WIDTH = 200
const DEFAULT_IMAGE_SAVE_NAME = "myclip.png"

// Must be same as the other clients
const HASH_ITERS_LOGIN = 20000
const HASH_ITERS_MSG = 10000
const HASH_LENGTH = 32

type Config struct {
	Server                 string
	Username               string
	Hash_login             string
	Hash_msg               string
	Disable_ssl_cert_check bool
}

var conf Config

// Init prepares the config paths and the clipboard. It remembers the Fyne app, which is
// needed for notifications and windows. Must be called once on startup
func Init(a fyne.App) {
	fyneApp = a
	initConfigPaths()
	if err := clipboard.Init(); err != nil {
		log.Println("Error: clipboard not available", err)
	}
}

// OpenConfigFile looks for config file in standard config folders and tries to open it
func OpenConfigFile() error {
	log.Println("Trying to open config file")
	viper.SetConfigName(CONFIG_FILENAME)
	viper.SetConfigType(CONFIG_TYPE)
	for _, v := range CONFIG_PATHS {
		viper.AddConfigPath(v)
	}
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	log.Println("Ok: Read Config", viper.ConfigFileUsed())
	return nil
}

// LoadConfigFromFile loads the credentials from the already opened config file
func LoadConfigFromFile() (Config, error) {
	log.Println("Loading config file to struct")
	if err := viper.Unmarshal(&conf); err != nil {
		log.Println("Error: Could not decode config into struct")
		return conf, err
	}
	log.Println("Ok: loaded config into struct for user", conf.Username, "on", conf.Server)
	return conf, nil
}

// WriteConfigFile writes config struct to file
func WriteConfigFile(c Config) error {
	log.Println("Writing config for user", c.Username, "on", c.Server)
	v := reflect.ValueOf(c)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		viper.Set(typeOfS.Field(i).Name, v.Field(i).Interface())
	}
	// contains the encryption key, so keep it private
	viper.SetConfigPermissions(0600)
	if err := viper.WriteConfigAs(CONFIG_FILEPATH); err != nil {
		log.Println("Error: writing config", err)
		return err
	}
	return os.Chmod(CONFIG_FILEPATH, 0600)
}

// initConfigPaths checks if at least one config folder exists, otherwise creates it
// it sets CONFIG_FILEPATH to this path
func initConfigPaths() {
	log.Printf("initConfigPaths")
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Panicln("Error:", err)
	}

	CONFIG_PATHS = []string{
		filepath.Join(homedir, ".config", "clipster"),
		filepath.Join(homedir, ".clipster"),
		filepath.FromSlash("/etc/clipster"),
	}

	for _, path := range CONFIG_PATHS {
		if fileExists(path) {
			log.Println("Config file folder exists", path)
			CONFIG_FILEPATH = filepath.Join(path, CONFIG_FILENAME)
			return
		}
	}

	log.Println("Error: No config file folder exists")
	log.Println("Creating config file folder", CONFIG_PATHS[0])
	if err := os.MkdirAll(CONFIG_PATHS[0], 0700); err != nil {
		log.Panicln(err)
	}
	CONFIG_FILEPATH = filepath.Join(CONFIG_PATHS[0], CONFIG_FILENAME)
	log.Println("Created config file folder", CONFIG_PATHS[0])
}

// fileExists checks if a file or folder exists
func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}
