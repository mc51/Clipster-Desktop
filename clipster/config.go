package clipster

import (
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/spf13/viper"
)

var CONFIG_PATHS = make([]string, 4)
var CONFIG_FILEPATH string

const CONFIG_FILENAME = "config.yaml"
const CONFIG_TYPE = "yaml"

const HOST_DEFAULT string = "https://clipster.cc"
const RE_HOSTNAME string = `^(https):\/\/[^\s\/$.?#].[^\s]*|://localhost:|://127.0.0.1:|^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
const ICON_FILENAME = "assets/clipster_icon_128.png"

const API_URI_COPY_PASTE = "/copy-paste/"
const API_URI_REGISTER = "/register/"
const API_URI_LOGIN = "/verify-user/"
const API_REQ_TIMEOUT = 6

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

func init() {
	initConfigPaths()
}

func OpenConfigFile() (bool, error) {
	// OpenConfigFile looks for config file in standard config folders and tries to open it
	log.Println("Trying to open config file")
	viper.SetConfigName(CONFIG_FILENAME)
	viper.SetConfigType(CONFIG_TYPE)
	for _, v := range CONFIG_PATHS {
		viper.AddConfigPath(v)
	}
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		return false, err
	}
	log.Println("Ok: Read Config\n", viper.AllSettings())
	return true, nil
}

func LoadConfigFromFile() (Config, error) {
	// LoadConfigFromFile loads the credentials from the already opened config file
	log.Println("Loading config file to struct")
	if err := viper.Unmarshal(&conf); err != nil {
		log.Println("Error: Could not decode config into struct")
		return conf, err
	}
	log.Println("Ok: loaded config into struct")
	return conf, nil
}

func WriteConfigFile(c Config) {
	// WriteConfigFile writes config struct to file
	log.Printf("Writing config: %+v", c)
	v := reflect.ValueOf(c)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		log.Printf("Field: %s\tValue: %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
		viper.Set(typeOfS.Field(i).Name, v.Field(i).Interface())
	}
	if err := viper.WriteConfigAs(CONFIG_FILEPATH); err != nil {
		log.Fatalln("Error:", err)
	}
}

func initConfigPaths() {
	// initConfigPaths checks if at least one config folder exists, otherwise creates it
	// it sets CONFIG_FILEPATH to this path
	log.Printf("initConfigPaths")
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("Error:", err)
	}

	CONFIG_PATHS[0] = filepath.Join(homedir, ".config", "clipster")
	CONFIG_PATHS[1] = filepath.Join(homedir, ".clipster")
	CONFIG_PATHS[2] = filepath.FromSlash("/etc/clipster")

	for _, path := range CONFIG_PATHS {
		if fileExists(path) {
			log.Println("Config file folder exists", path)
			CONFIG_FILEPATH = path + string(os.PathSeparator) + CONFIG_FILENAME
			return
		}
	}

	log.Println("Error: No config file folder exists")
	log.Println("Creating config file folder", CONFIG_PATHS[0])
	if err := os.MkdirAll(CONFIG_PATHS[0], 0775); err != nil {
		log.Fatal(err)
		return
	}
	CONFIG_FILEPATH = CONFIG_PATHS[0] + string(os.PathSeparator) + CONFIG_FILENAME
	log.Println("Created config file folder", CONFIG_PATHS[0])
}

func fileExists(p string) bool {
	// fileExists checks if a file or folder exists
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		return true
	}
	return false
}
