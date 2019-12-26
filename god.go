package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/urfave/cli"
)

type MessageAuthor struct {
	ID string `json:"id"`
}

type ChannelMessage struct {
	ID      string        `json:"id"`
	Content string        `json:"content"`
	Author  MessageAuthor `json:"author"`
	Hit     bool          `json:"hit"`
}

type ChannelMessagesResponse struct {
	TotalResults uint32             `json:"total_results"`
	Messages     [][]ChannelMessage `json:"messages"`
}

type AccountConfiguration struct {
	AuthenticationToken string `json:"authToken"`
	UserID              string `json:"userId"`
}

type Client struct {
	baseURL       string
	Configuration AccountConfiguration
}

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
	app.Version = "1.0.0"
	app.Usage = "make it an accord"
	app.Description = "Get off discord completely with a single command."
	app.Action = func(ctx *cli.Context) error {
		cfg, err := loadConfiguration(accCfgPath)
		if err != nil {
			return errors.New("couldn't load account configuration, please verify your file")
		}

		client := Client{baseURL, cfg}
		client.startDeletion(channelID)

		return nil
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "account",
			Usage:       "Load account configuration from `FILE`",
			Aliases:     []string{"a"},
			TakesFile:   true,
			Required:    true,
			Destination: &accCfgPath,
		},
		&cli.StringFlag{
			Name:        "channel",
			Aliases:     []string{"c"},
			Usage:       "Specify the `ID` of the channel",
			Required:    true,
			Destination: &channelID,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func (client Client) startDeletion(channelID string) {
	resp, err := client.loadMessages(channelID)
	if err != nil {
		panic(err)
	}

	initialMsgCount := resp.TotalResults
	removedMesages := 0

	for resp.TotalResults != 0 {
		const maxRetries int = 15

		for _, messages := range resp.Messages {
			for _, msg := range messages {
				if msg.Author.ID == client.Configuration.UserID && msg.Hit {
					waitTime := 250
					retries := 0
					for retries < maxRetries {
						time.Sleep(time.Duration(waitTime) * time.Millisecond)

						err = client.deleteMessage(channelID, msg)
						if err != nil {
							retries++
							waitTime += 50
							fmt.Printf("(%s) %s (retry no. %d/%d, waiting %dms) ['%s']\n", channelID, err, retries, maxRetries, waitTime, msg.Content)
						} else {
							removedMesages++
							fmt.Printf("(%s) [%d/%d] Removed message: \"%s\" \n", channelID, removedMesages, initialMsgCount, msg.Content)
							break
						}
					}
				}
			}
		}

		resp, err = client.loadMessages(channelID)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("No more messages have been found on channel %s.\n", channelID)
}

func (client Client) loadMessages(channelID string) (ChannelMessagesResponse, error) {
	var serverResponse ChannelMessagesResponse
	httpClient := http.Client{}

	url := fmt.Sprintf("%s/%s/messages/search?author_id=%s", client.baseURL, channelID, client.Configuration.UserID)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return serverResponse, err
	}
	request.Header.Set("Authorization", client.Configuration.AuthenticationToken)

	serverIndexed := false
	var resp *http.Response
	for !serverIndexed {
		resp, err = httpClient.Do(request)
		if err != nil {
			return serverResponse, err
		}

		serverIndexed = resp.StatusCode == http.StatusOK
		if serverIndexed {
			break
		}

		fmt.Printf("The server has not been yet indexed. Trying again in 1s... (Received code: %d) \n", resp.StatusCode)
		time.Sleep(1 * time.Second)
	}

	defer resp.Body.Close()
	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return serverResponse, err
	}

	err = json.Unmarshal(bodyData, &serverResponse)
	if err != nil {
		return serverResponse, err
	}

	return serverResponse, nil
}

func (client Client) deleteMessage(channelID string, message ChannelMessage) error {
	delurl := fmt.Sprintf("%s/%s/messages/%s", client.baseURL, channelID, message.ID)

	request, err := http.NewRequest("DELETE", delurl, nil)
	if err != nil {
		return nil
	}
	request.Header.Set("Authorization", client.Configuration.AuthenticationToken)

	httpClient := http.Client{}
	resp, err := httpClient.Do(request)
	if err != nil {
		return err
	}

	if (resp.StatusCode >= 200 && resp.StatusCode <= 299) || (resp.StatusCode == 404) {
		return nil
	}

	return fmt.Errorf("failed to remove the message with id %s (Response code: %d)", message.ID, resp.StatusCode)
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
