package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type AccountConfiguration struct {
	AuthenticationToken string `json:"authToken"`
	ChannelID           uint32 `json:"channelId"`
	UserID              uint32 `json:"userId"`
}

func main() {
	if len(os.Args) <= 1 {
		panic("Path to account configuration was not provided.")
	}

	cfg, err := loadConfiguration(os.Args[1])
	if err != nil {
		panic(err)
	}

	fmt.Println("Your account configuration is: ", cfg.AuthenticationToken, cfg.ChannelID, cfg.UserID)
}

func loadConfiguration(cfgPath string) (AccountConfiguration, error) {
	cfg := AccountConfiguration{}

	cfgFile, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(cfgFile, &cfg)

	return cfg, err
}
