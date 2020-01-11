package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	const baseURL string = "https://discordapp.com/api/v6/channels"

	var accCfgPath string
	var channelID string

	app := cli.NewApp()
	app.Name = "GOD (Go Off Discord)"
	app.Authors = []*cli.Author{
		&cli.Author{
			Name: "Mihai Stan",
		},
	}

	app.Version = "1.1.0"
	app.Usage = "make it an accord"
	app.Description = "Get off discord completely with a single command."

	app.Commands = []*cli.Command{
		{
			Name:    "delete",
			Aliases: []string{"del"},
			Usage:   "Perform a deletion action on a server or channel",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "account",
					Usage:       "Load account configuration from `FILE`",
					Aliases:     []string{"a"},
					TakesFile:   true,
					Required:    true,
					Destination: &accCfgPath,
				},
			},
			Subcommands: []*cli.Command{
				{
					Name: "channel",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "id",
							Usage:       "Specify the `ID` of the channel",
							Required:    true,
							Destination: &channelID,
						},
					},
					Action: func(ctx *cli.Context) error {
						cfg, err := loadConfiguration(accCfgPath)
						if err != nil {
							return errors.New("couldn't load account configuration, please verify your file")
						}

						client := Client{baseURL, cfg}
						client.startDeletion(channelID)

						return nil
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
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
