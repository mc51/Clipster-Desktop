package clipster

import (
	"github.com/spf13/viper"
)

const HOST_DEFAULT string = "https://clipster.cc"
const RE_HOSTNAME string = `^(https):\/\/[^\s\/$.?#].[^\s]*|://localhost:|://127.0.0.1:|^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
const ICON_FILENAME = "assets/clipster_icon_128.png"

const API_URI_COPY_PASTE = "/copy-paste/"
const API_URI_REGISTER = "/register/"
const API_URI_LOGIN = "/verify-user/"
const API_TIMEOUT = 6

func LoadConfig() (bool, error) {

	viper.SetDefault("server", HOST_DEFAULT)
	viper.SetDefault("verify_ssl_cert", "True")
	if err := ReadConfigFile(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found - ask for credentials
			// viper.WriteConfig()
			return false, nil
		} else {
			// Config file was found but another error was produced
			return false, err
		}
	}
	return true, nil
}

func ReadConfigFile() error {

	viper.SetConfigName("config")
	viper.SetConfigType("ini")
	viper.AddConfigPath("$HOME/.config/clipster")
	viper.AddConfigPath("$HOME/.clipster")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		return err
	}
	return nil
}
