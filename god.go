package main

import (
	"encoding/json"
	"errors"
	"github.com/Mehigh17/go-off-discord/discord"
	"io/ioutil"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	var accCfgPath, channelID, serverID string

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

						client := discord.Client {
							Configuration:	cfg,
						}
						client.DeleteChannel(channelID)

						return nil
					},
				},
				{
					Name: "server",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "id",
							Usage:       "Specify the `ID` of the server",
							Required:    true,
							Destination: &serverID,
						},
					},
					Action: func(ctx *cli.Context) error {
						cfg, err := loadConfiguration(accCfgPath)
						if err != nil {
							return errors.New("couldn't load account configuration, please verify your file")
						}

						client := discord.Client {
							Configuration:	cfg,
						}
						client.DeleteServer(serverID)

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

func loadConfiguration(cfgPath string) (discord.AccountConfiguration, error) {
	cfg := discord.AccountConfiguration{}

	cfgFile, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(cfgFile, &cfg)

	return cfg, err
}
